package main

import (
	"github.com/ajlima/go-ms-arch-example/internal/config"
	"github.com/spf13/viper"
)

func main() {
	config.NewConfig(viper.GetViper())

}
