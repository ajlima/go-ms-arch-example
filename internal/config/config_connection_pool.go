package config

import (
	"context"
	"runtime"
	"strconv"

	"github.com/ajlima/go-ms-arch-example/internal/app"
	pool "github.com/jolestar/go-commons-pool/v2"
	"github.com/segmentio/kafka-go"
)

type KafkaConnection struct {
	Conn *kafka.Conn
}

type KafkaConnectionFactory struct {
	applicationContext *app.ApplicationContext
}

func NewKafkaConnectionPool(ctx context.Context, appContext *app.ApplicationContext) *pool.ObjectPool {
	factory := &KafkaConnectionFactory{
		applicationContext: appContext,
	}

	pconfig := pool.NewDefaultPoolConfig()
	pconfig.MaxTotal = runtime.NumCPU() * 2
	pconfig.MinIdle = runtime.NumCPU() / 2
	pconfig.MaxIdle = runtime.NumCPU() * 2

	p := pool.NewObjectPool(ctx, factory, pconfig)

	return p
}

func (f *KafkaConnectionFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	log := f.applicationContext.Log

	log.Info("Starting kafka connection")

	partition, err := strconv.Atoi(f.applicationContext.Viper.GetString(KAFKA_PARTITION))
	if err != nil {
		partition = 0
	}

	kafkaConn, err := kafka.DialLeader(
		ctx,
		"tcp",
		f.applicationContext.Viper.GetString(KAFKA_BROKERS),
		f.applicationContext.Viper.GetString(KAFKA_OUT_TOPIC),
		partition,
	)
	if err != nil {
		log.Panic("Failed to dial leader: ", err)
	}

	return pool.NewPooledObject(&KafkaConnection{
		Conn: kafkaConn,
	}), nil
}

func (f *KafkaConnectionFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	log := f.applicationContext.Log

	log.Info("Closing kafka connection")

	conn := f.GetNativeConnection(object)
	if err := conn.Close(); err != nil {
		log.Error("Problem in Kafka close connection ", err)
	}
	return nil
}

func (f *KafkaConnectionFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	// do validate
	return true
}

func (f *KafkaConnectionFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	// do activate
	return nil
}

func (f *KafkaConnectionFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	// do passivate
	return nil
}

func (f *KafkaConnectionFactory) GetNativeConnection(object *pool.PooledObject) *kafka.Conn {
	return object.Object.(*KafkaConnection).Conn
}
