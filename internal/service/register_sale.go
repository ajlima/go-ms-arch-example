package service

import (
	"time"

	"github.com/ajlima/go-ms-arch-example/internal/app"
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
	conn := r.applicationContext.KafkaConn
	log := r.applicationContext.Log

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.Write(msg)
	if err != nil {
		log.Fatal("Failed to write messages: ", err)
		return err
	}
	return nil
}
