package service

import (
	"context"
	"strconv"
	"time"

	"github.com/ajlima/go-ms-arch-example/internal/app"
	"github.com/ajlima/go-ms-arch-example/internal/config"
	"github.com/segmentio/kafka-go"
)

type RegisterSaleService struct {
	applicationContext *app.ApplicationContext
}

func NewRegisterSaleService(appContext *app.ApplicationContext) RegisterSaleService {
	return RegisterSaleService{
		applicationContext: appContext,
	}
}

func (r RegisterSaleService) SendMessage(msg []byte) (err error) {
	// conn := r.applicationContext.KafkaConn
	log := r.applicationContext.Log

	partition, err := strconv.Atoi(r.applicationContext.Viper.GetString(config.KAFKA_PARTITION))
	if err != nil {
		partition = 0
	}

	kafkaConn, err := kafka.DialLeader(
		context.Background(),
		"tcp",
		r.applicationContext.Viper.GetString(config.KAFKA_BROKERS),
		r.applicationContext.Viper.GetString(config.KAFKA_OUT_TOPIC),
		partition,
	)

	if err != nil {
		log.Panic("failed to dial leader:", err)
	}

	defer kafkaConn.Close()

	kafkaConn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = kafkaConn.Write(msg)
	if err != nil {
		log.Fatal("Failed to write messages: ", err)
		return err
	}

	return nil
}
