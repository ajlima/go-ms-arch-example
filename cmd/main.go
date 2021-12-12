package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/ajlima/go-ms-arch-example/internal/app"
	"github.com/ajlima/go-ms-arch-example/internal/config"
	"github.com/ajlima/go-ms-arch-example/internal/http/handler"
	"github.com/ajlima/go-ms-arch-example/internal/service"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
	ginlogrus "github.com/toorop/gin-logrus"
)

var (
	vip        = viper.GetViper()
	log        = logrus.New()
	appContext *app.ApplicationContext
)

func init() {
	vip = config.ConfigureEnvironment(vip)
	log = config.ConfigureLogger(log, viper.GetString(config.LOG_FILE), viper.GetString(config.LOG_LEVEL))

	log.Println("*")
	log.Println("* Starting with following configuration: ")
	log.Println("*")
	for _, key := range vip.AllKeys() {
		log.Printf("* %s = %s", key, vip.GetString(key))
	}

	appContext = app.NewApplicationContext(
		vip,
		log,
	)

	// ctx := context.Background()
	// cp = config.NewKafkaConnectionPool(ctx, appContext)
	// appContext.SetKafkaConnectionPool(cp)
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
	router := gin.New()
	router.Use(gin.Recovery())

	// Using this request logging, the number of requests handled by this application will be decreased
	// if the application need to handle millions of requests per minute this log should be configured to OFF
	requestLog := strings.ToUpper(viper.GetString(config.HTTP_REQUEST_LOG))
	switch requestLog {
	case "YES", "TRUE", "ON":
		router.Use(ginlogrus.Logger(log))
	}

	deliveryChan := make(chan []byte, runtime.NumCPU()*2)
	closeDelivery := make(chan byte, 1)

	apiV1 := router.Group("/api/v1")

	router.GET("/swagger/*any", ginswagger.WrapHandler(swaggerFiles.Handler))

	registerSaleService := service.NewRegisterSaleService(appContext, deliveryChan)
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

	kafkaServers := viper.GetString(config.KAFKA_BROKERS)
	topic := viper.GetString(config.KAFKA_OUT_TOPIC)
	conf := kafka.ConfigMap{
		"bootstrap.servers": kafkaServers,
	}
	p, err := kafka.NewProducer(&conf)
	if err != nil {
		log.Panicf("It wasn't possible to create a producer to kafka %s", kafkaServers)
	}
	defer p.Close()
	defer p.Flush(100)

	// Kafka delivery channel goroutine
	run := true
	go func() {
		for run {
			select {
			case msg := <-deliveryChan:
				err = p.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
					Value:          msg},
					nil,
				)
				if err != nil {
					log.Error("Error sending a message to kafka topic: ", err)
					time.Sleep(500 * time.Millisecond)
				}
			case <-closeDelivery:
				log.Info("Closing delivery channel")
				run = false
			}
		}
		close(deliveryChan)
		close(closeDelivery)
	}()

	// Kafka delivery channel listener goroutine
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Errorf("Failed to deliver message: %v", ev.TopicPartition)
				} else {
					log.Infof("Successfully produced record to topic %s partition [%d] @ offset %v",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			}
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	closeDelivery <- byte(1)

	log.Println("")
	log.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	// catching ctx.Done(). timeout of 1 seconds.
	<-ctx.Done()
	log.Info("Server exiting")
}
