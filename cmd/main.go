package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ajlima/go-ms-arch-example/internal/app"
	"github.com/ajlima/go-ms-arch-example/internal/config"
	"github.com/ajlima/go-ms-arch-example/internal/http/handler"
	"github.com/ajlima/go-ms-arch-example/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
	ginlogrus "github.com/toorop/gin-logrus"
)

var (
	vip        = viper.GetViper()
	log        = logrus.New()
	kafkaConn  *kafka.Conn
	appContext *app.ApplicationContext
)

func init() {
	vip = config.ConfigureEnvironment(vip)
	log = config.ConfigureLogger(log, viper.GetString(config.LOG_FILE), viper.GetString(config.LOG_LEVEL))

	log.Println("*")
	log.Println("* Starting with following configuration: ")
	log.Println("*")
	for _, key := range vip.AllKeys() {
		log.Printf("* %s = %s\n", key, vip.GetString(key))
	}

	appContext = app.NewApplicationContext(
		vip,
		log,
		kafkaConn,
	)
}

// @title           Microservice GO example
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
func main() {
	defer kafkaConn.Close()

	router := gin.New()
	router.Use(ginlogrus.Logger(log), gin.Recovery())
	apiV1 := router.Group("/api/v1")

	router.GET("/swagger/*any", ginswagger.WrapHandler(swaggerFiles.Handler))

	registerSaleService := service.NewRegisterSaleService(appContext)
	handler.NewRegisterSaleHandler(appContext, registerSaleService, apiV1)

	srv := &http.Server{
		Addr:    ":" + vip.GetString(config.HTTP_PORT),
		Handler: router,
	}

	// Execute gin engine in separated goroutine to use main goroutine to handle graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("")
	log.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 1 seconds.
	<-ctx.Done()
	log.Info("timeout of 1 seconds.")
	log.Info("Server exiting")
}
