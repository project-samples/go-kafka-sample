package app

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/core-go/health"
	hm "github.com/core-go/health/mongo"
	hk "github.com/core-go/health/sarama"
	"github.com/core-go/kafka/sarama"
	w "github.com/core-go/mongo/writer"
	"github.com/core-go/mq"
	v "github.com/core-go/mq/validator"
	"github.com/core-go/mq/zap"
)

type ApplicationContext struct {
	HealthHandler *health.Handler
	Receive       func(ctx context.Context, handle func(context.Context, []byte))
	Handle        func(context.Context, []byte)
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

	consumer, er2 := kafka.NewConsumerByConfig(cfg.Consumer, logError, true)
	if er2 != nil {
		log.Error(ctx, "Cannot create a new consumer. Error: "+er2.Error())
		return nil, er2
	}
	validator, err := v.NewValidator[*User]()
	if err != nil {
		return nil, err
	}
	errorHandler := mq.NewErrorHandler[*User](logError)

	writer := w.NewWriter[*User](db, "user")
	handler := mq.NewHandlerByConfig[User](cfg.Handler, writer.Write, validator.Validate, errorHandler.Reject, errorHandler.HandleError, logError, logInfo)
	mongoChecker := hm.NewHealthChecker(client)
	consumerChecker := hk.NewKafkaHealthChecker(cfg.Consumer.Brokers, "kafka_consumer")
	healthHandler := health.NewHandler(mongoChecker, consumerChecker)

	return &ApplicationContext{
		HealthHandler: healthHandler,
		Receive:       consumer.ConsumeValue,
		Handle:        handler.Handle,
	}, nil
}
