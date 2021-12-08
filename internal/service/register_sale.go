package service

import "github.com/ajlima/go-ms-arch-example/internal/app"

type RegisterSaleService struct {
	applicationContext *app.ApplicationContext
}

func NewRegisterSaleService(appContext *app.ApplicationContext) RegisterSaleService {
	return RegisterSaleService{
		applicationContext: appContext,
	}
}
