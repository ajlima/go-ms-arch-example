package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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

	partition, err := strconv.Atoi(viper.GetString(config.KAFKA_PARTITION))
	if err != nil {
		partition = 0
	}

	kafkaConn, err = kafka.DialLeader(context.Background(), "tcp", viper.GetString(config.KAFKA_BROKERS), viper.GetString(config.KAFKA_OUT_TOPIC), partition)
	if err != nil {
		log.Panic("failed to dial leader:", err)
	}

	appContext = app.NewApplicationContext(
		vip,
		log,
		kafkaConn,
	)
}

func main() {
	defer kafkaConn.Close()

	router := gin.New()
	router.Use(ginlogrus.Logger(log), gin.Recovery())

	registerSaleService := service.NewRegisterSaleService(appContext)
	handler.NewRegisterSaleHandler(appContext, registerSaleService, router)

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
