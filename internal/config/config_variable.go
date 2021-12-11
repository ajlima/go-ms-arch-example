package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Service struct {
	Id          string
	Description string
	Default     interface{}
}

const (
	SERVICE_NAME            = "SERVICE_NAME"
	HTTP_PORT               = "HTTP_PORT"
	DATABASE_URL            = "DATABASE_URL"
	DATABASE_PORT           = "DATABASE_PORT"
	DATABASE_TYPE           = "DATABASE_TYPE"
	DATABASE_NAME           = "DATABASE_NAME"
	DATABASE_USERNAME       = "DATABASE_USERNAME"
	DATABASE_PASSWORD       = "DATABASE_PASSWORD"
	KAFKA_BROKERS           = "KAFKA_BROKERS"
	KAFKA_PARTITION         = "KAFKA_PARTITION"
	KAFKA_IN_TOPIC          = "KAFKA_IN_TOPIC"
	KAFKA_IN_CONSUMER_GROUP = "KAFKA_IN_CONSUMER_GROUP"
	KAFKA_IN_OFFSET         = "KAFKA_IN_OFFSET"
	KAFKA_OUT_TOPIC         = "KAFKA_OUT_TOPIC"
	PROFILE                 = "PROFILE"
	LOG_FILE                = "LOG_FILE"
	LOG_LEVEL               = "LOG_LEVEL"
	HTTP_REQUEST_LOG        = "HTTP_REQUEST_LOG"
)

var (
	executableName = strings.ToLower(filepath.Base(os.Args[0]))
	services       = map[string]Service{
		SERVICE_NAME: {
			Id:          SERVICE_NAME,
			Description: "Name of service",
			Default:     executableName,
		},
		HTTP_PORT: {
			Id:          HTTP_PORT,
			Description: "HTTP port for service",
			Default:     8080,
		},
		DATABASE_URL: {
			Id:          DATABASE_URL,
			Description: "Database URL",
			Default:     "localhost",
		},
		DATABASE_PORT: {
			Id:          DATABASE_PORT,
			Description: "Database port",
			Default:     5432,
		},
		DATABASE_TYPE: {
			Id:          DATABASE_TYPE,
			Description: "Type of database (PostgreSQL, MySQL, MSSQL, Oracle)",
			Default:     "PostgreSQL",
		},
		DATABASE_NAME: {
			Id:          DATABASE_NAME,
			Description: "Name of database",
			Default:     "",
		},
		DATABASE_USERNAME: {
			Id:          DATABASE_USERNAME,
			Description: "Username to connect to database",
			Default:     "",
		},
		DATABASE_PASSWORD: {
			Id:          DATABASE_PASSWORD,
			Description: "Password to connect to database",
			Default:     "",
		},
		KAFKA_BROKERS: {
			Id:          KAFKA_BROKERS,
			Description: "Kafka brokers separated by ,",
			Default:     "localhost:9092",
		},
		KAFKA_PARTITION: {
			Id:          KAFKA_PARTITION,
			Description: "Kafka partition",
			Default:     "0",
		},
		KAFKA_IN_TOPIC: {
			Id:          KAFKA_IN_TOPIC,
			Description: "Name of inbound kafka topic",
			Default:     "in-" + filepath.Base(os.Args[0]),
		},
		KAFKA_IN_CONSUMER_GROUP: {
			Id:          KAFKA_IN_CONSUMER_GROUP,
			Description: "Name of inbound consumer group",
			Default:     "",
		},
		KAFKA_IN_OFFSET: {
			Id:          KAFKA_IN_OFFSET,
			Description: "Name of inbound offset position",
			Default:     "",
		},
		KAFKA_OUT_TOPIC: {
			Id:          KAFKA_OUT_TOPIC,
			Description: "Name of outbound kafka topic",
			Default:     "out-" + filepath.Base(os.Args[0]),
		},
		PROFILE: {
			Id:          PROFILE,
			Description: "Name of environment where the app is running (dev, int, uat, prd)",
			Default:     "dev",
		},
		LOG_FILE: {
			Id:          LOG_FILE,
			Description: "Path to logfile",
			Default:     "stdout",
		},
		LOG_LEVEL: {
			Id:          LOG_LEVEL,
			Description: "Log level (panic,fatal,error,warn,info,debug,trace)",
			Default:     "info",
		},
		HTTP_REQUEST_LOG: {
			Id:          HTTP_REQUEST_LOG,
			Description: "Turn http request log on (ON|OFF)",
			Default:     "ON",
		},
	}
)

func ConfigureEnvironment(viper *viper.Viper) *viper.Viper {
	viper.AutomaticEnv()
	configureDefaults(viper)
	bindCommandFlags(viper)
	return viper
}

func configureDefaults(viper *viper.Viper) {
	for _, value := range services {
		viper.SetDefault(value.Id, value.Default)
	}
}

func bindCommandFlags(viper *viper.Viper) {
	for _, value := range services {
		flag.String(value.Id, fmt.Sprintf("%v", value.Default), value.Description)
	}
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}
