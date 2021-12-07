package main

import (
	"github.com/ajlima/go-ms-arch-example/internal/app"
	"github.com/ajlima/go-ms-arch-example/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	_vipper     = viper.GetViper()
	_logger     = logrus.New()
	_appContext *app.ApplicationContext
)

func init() {
	_vipper = config.ConfigureEnvironment(_vipper)
	_logger = config.ConfigureLogger(_logger, viper.GetString(config.LOG_FILE), viper.GetString(config.LOG_LEVEL))
	_appContext = app.NewApplicationContext(
		_vipper,
		_logger,
	)
}

func main() {
	_logger.Info("Teste de log com configuração no init")
}
