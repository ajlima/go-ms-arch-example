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
	router              *gin.RouterGroup
	registerSaleService service.RegisterSaleService
}

func NewRegisterSaleHandler(appContext *app.ApplicationContext, rss service.RegisterSaleService, router *gin.RouterGroup) RegisterSaleHandler {
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

// RegisterSale godoc
// @Summary      Register one sale
// @Description  Register one sale record
// @Tags         sale
// @Accept       json
// @Produce      json
// @Param        transaction	body      datastruct.Transaction  true  "Register sale"
// @Success      200      		{object}  datastruct.Transaction
// @Failure      400      		{object}  datastruct.Err
// @Failure      404      		{object}  datastruct.Err
// @Failure      500      		{object}  datastruct.Err
// @Router       /register/sale [post]
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
	c.JSON(http.StatusOK, transaction)
}
