package app

import (
	pool "github.com/jolestar/go-commons-pool/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ApplicationContext struct {
	Viper               *viper.Viper
	Log                 *logrus.Logger
	KafkaConnectionPool *pool.ObjectPool
}

func NewApplicationContext(
	viper *viper.Viper,
	log *logrus.Logger,
) *ApplicationContext {
	return &ApplicationContext{
		Viper: viper,
		Log:   log,
	}
}

func (a *ApplicationContext) SetKafkaConnectionPool(kafkaConnectionPool *pool.ObjectPool) {
	a.KafkaConnectionPool = kafkaConnectionPool
}
