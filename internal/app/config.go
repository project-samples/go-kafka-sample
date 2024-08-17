package app

import (
	"github.com/core-go/kafka"
	"github.com/core-go/mq"
	"github.com/core-go/mq/server"
	"github.com/core-go/mq/zap"
)

type Config struct {
	Server      server.ServerConfig   `mapstructure:"server"`
	Log         log.Config            `mapstructure:"log"`
	Mongo       MongoConfig           `mapstructure:"mongo"`
	Reader      kafka.ReaderConfig    `mapstructure:"reader"`
	KafkaWriter *kafka.WriterConfig   `mapstructure:"writer"`
	Retry       mq.RetryHandlerConfig `mapstructure:"retry"`
}

type MongoConfig struct {
	Uri      string `yaml:"uri" mapstructure:"uri" json:"uri,omitempty" gorm:"column:uri" bson:"uri,omitempty" dynamodbav:"uri,omitempty" firestore:"uri,omitempty"`
	Database string `yaml:"database" mapstructure:"database" json:"database,omitempty" gorm:"column:database" bson:"database,omitempty" dynamodbav:"database,omitempty" firestore:"database,omitempty"`
}
