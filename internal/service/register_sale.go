package service

import (
	"context"
	"time"

	"github.com/ajlima/go-ms-arch-example/internal/app"
	"github.com/ajlima/go-ms-arch-example/internal/config"
)

type RegisterSaleService struct {
	applicationContext *app.ApplicationContext
}

func NewRegisterSaleService(appContext *app.ApplicationContext) RegisterSaleService {
	return RegisterSaleService{
		applicationContext: appContext,
	}
}

func (r RegisterSaleService) SendMessage(ctx context.Context, msg []byte) (err error) {
	log := r.applicationContext.Log

	kafkaConnection, err := r.applicationContext.KafkaConnectionPool.BorrowObject(ctx)
	if err != nil {
		log.Fatal("It was impossible to get one connection from kafka connection pool")
	}
	defer r.applicationContext.KafkaConnectionPool.ReturnObject(ctx, kafkaConnection)

	conn := kafkaConnection.(*config.KafkaConnection).Conn
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if _, err = conn.Write(msg); err != nil {
		log.Fatal("Failed to write messages: ", err)
		return err
	}

	return nil
}
