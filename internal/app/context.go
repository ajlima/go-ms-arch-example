package app

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ApplicationContext struct {
	Viper *viper.Viper
	Log   *logrus.Logger
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
