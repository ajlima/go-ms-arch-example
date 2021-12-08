package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ajlima/go-ms-arch-example/internal/app"
	"github.com/ajlima/go-ms-arch-example/internal/http/datastruct"
	"github.com/ajlima/go-ms-arch-example/internal/service"
	"github.com/gin-gonic/gin"
)

type RegisterSaleHandler struct {
	applicationContext  *app.ApplicationContext
	registerSaleService service.RegisterSaleService
	router              *gin.Engine
}

func NewRegisterSaleHandler(appContext *app.ApplicationContext, rss service.RegisterSaleService, router *gin.Engine) RegisterSaleHandler {
	h := RegisterSaleHandler{
		applicationContext:  appContext,
		registerSaleService: rss,
		router:              router,
	}
	h.configureRoutes()
	return h
}

func (h RegisterSaleHandler) configureRoutes() {
	h.router.POST("/register/sale/", h.registerSale)
}

func (h RegisterSaleHandler) registerSale(c *gin.Context) {
	var transaction datastruct.Transaction
	if err := c.Bind(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rowdata, err := json.Marshal(transaction)
	if err != nil {
		log.Panic("Error transforming body in []byte: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = h.registerSaleService.SendMessage(rowdata)
	if err != nil {
		log.Panic("Error sending message to kakfa: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "OK"})
}
