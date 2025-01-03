package app

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/core-go/health"
	hk "github.com/core-go/health/kafka"
	hm "github.com/core-go/health/mongo"
	"github.com/core-go/kafka"
	w "github.com/core-go/mongo/writer"
	"github.com/core-go/mq"
	v "github.com/core-go/mq/validator"
	"github.com/core-go/mq/zap"
)

type ApplicationContext struct {
	HealthHandler *health.Handler
	Read          func(ctx context.Context, handle func(context.Context, []byte, map[string]string))
	Handle        func(context.Context, []byte, map[string]string)
}

func NewApp(ctx context.Context, cfg Config) (*ApplicationContext, error) {
	log.Initialize(cfg.Log)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.Uri))
	if err != nil {
		log.Error(ctx, "Cannot connect to MongoDB: Error: "+err.Error())
		return nil, err
	}
	db := client.Database(cfg.Mongo.Database)

	logError := log.ErrorMsg
	var logInfo func(context.Context, string)
	if log.IsInfoEnable() {
		logInfo = log.InfoMsg
	}

	reader, er2 := kafka.NewReaderByConfig(cfg.Reader, logError, true)
	if er2 != nil {
		log.Error(ctx, "Cannot create a new reader. Error: "+er2.Error())
		return nil, er2
	}
	validator, err := v.NewValidator[*User]()
	if err != nil {
		return nil, err
	}
	errorHandler := mq.NewErrorHandler[*User](logError)
	producer, err := kafka.NewWriterByConfig(*cfg.KafkaWriter, nil)
	if err != nil {
		return nil, err
	}
	writer := w.NewWriter[*User](db, "user")
	handler := mq.NewRetryHandlerByConfig[User](cfg.Retry, writer.Write, validator.Validate, errorHandler.RejectWithMap, nil, producer.Write, logError, logInfo)
	mongoChecker := hm.NewHealthChecker(client)
	readerChecker := hk.NewKafkaHealthChecker(cfg.Reader.Brokers, "kafka_reader")
	producerChecker := hk.NewKafkaHealthChecker(cfg.KafkaWriter.Brokers, "kafka_producer")
	healthHandler := health.NewHandler(mongoChecker, readerChecker, producerChecker)

	return &ApplicationContext{
		HealthHandler: healthHandler,
		Read:          reader.Read,
		Handle:        handler.Handle,
	}, nil
}
