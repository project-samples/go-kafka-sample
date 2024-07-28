package app

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/core-go/health"
	hm "github.com/core-go/health/mongo"
	"github.com/core-go/kafka/confluent"
	w "github.com/core-go/mongo/writer"
	"github.com/core-go/mq"
	v "github.com/core-go/mq/validator"
	"github.com/core-go/mq/zap"
)

type ApplicationContext struct {
	HealthHandler *health.Handler
	Receive       func(ctx context.Context, handle func(context.Context, []byte, map[string]string))
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

	receiver, er2 := kafka.NewConsumerByConfig(cfg.Reader, logError)
	if er2 != nil {
		log.Error(ctx, "Cannot create a new receiver. Error: "+er2.Error())
		return nil, er2
	}
	validator, err := v.NewValidator[*User]()
	if err != nil {
		return nil, err
	}
	errorHandler := mq.NewErrorHandler[*User](logError)
	sender, err := kafka.NewProducerByConfig(*cfg.KafkaWriter, nil)
	if err != nil {
		return nil, err
	}
	writer := w.NewWriter[*User](db, "user")
	han := mq.NewRetryHandlerByConfig[User](cfg.Retry, writer.Write, validator.Validate, errorHandler.RejectWithMap, nil, sender.Produce, logError, logInfo)
	mongoChecker := hm.NewHealthChecker(client)
	receiverChecker := kafka.NewKafkaHealthChecker(cfg.Reader.Brokers, "kafka_consumer")
	senderChecker := kafka.NewKafkaHealthChecker(cfg.KafkaWriter.Brokers, "kafka_producer")
	healthHandler := health.NewHandler(mongoChecker, receiverChecker, senderChecker)

	return &ApplicationContext{
		HealthHandler: healthHandler,
		Receive:       receiver.Consume,
		Handle:        han.Handle,
	}, nil
}
