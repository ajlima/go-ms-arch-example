package app

import (
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ApplicationContext struct {
	Viper     *viper.Viper
	Log       *logrus.Logger
	KafkaConn *kafka.Conn
}

func NewApplicationContext(
	viper *viper.Viper,
	log *logrus.Logger,
	kafkaConn *kafka.Conn,
) *ApplicationContext {
	return &ApplicationContext{
		Viper:     viper,
		Log:       log,
		KafkaConn: kafkaConn,
	}
}
