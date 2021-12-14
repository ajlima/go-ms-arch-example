package service

import (
	"context"

	"github.com/ajlima/go-ms-arch-example/internal/app"
)

type RegisterSaleService struct {
	applicationContext *app.ApplicationContext
	deliveryChan       chan []byte
}

type Producer interface {
	SendMessage(context.Context, []byte) error
}

func NewRegisterSaleService(appContext *app.ApplicationContext, dc chan []byte) *RegisterSaleService {
	return &RegisterSaleService{
		applicationContext: appContext,
		deliveryChan:       dc,
	}
}

func (r *RegisterSaleService) SendMessage(ctx context.Context, msg []byte) (err error) {
	r.deliveryChan <- msg
	return nil
}
