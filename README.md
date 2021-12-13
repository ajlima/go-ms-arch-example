# go-ms-arch-example

The aim of this project is to learn Go.

As a case study, I chose to develop something similar to what I do every day at work.

Therefore, this example must implement a microservice that aims to receive requests and send them to a kafka topic which will be consumed by another microservice and saved in a database.

I preferred not to use Go-kit, so I chose several frameworks to cover various parts of what I think are fundamental to create a microservice.

- Configuration: Viper
- HTTP: Gin
- Swagger: Gin-swagger
- Logger: Logrus
- Kafka: Confluent-kafka-go
- ORM: Gorm

As a development architecture I chose something similar to Spring, with dependency injection by the constructor and an initial context concept.

## To build this app:

```
make build
```

## To run from docker:

Use this command to startup everything that is necessary to run this application, including kafka, zookeaper, postgres and the application

```
docker-compose up
```

## To execute from localhost:

```
GIN_MODE=release out/bin/ms-arch-example \
--HTTP_PORT=8081 \
--KAFKA_OUT_TOPIC=sales_transactions \
--KAFKA_BROKERS=kafka:9092 \
--LOG_FILE=./ms-arch-example.log
```

INFO[0000] *                                            
INFO[0000] * Starting with following configuration:     
INFO[0000] *                                            
INFO[0000] * database_password =                        
INFO[0000] * kafka_out_topic = sales_transactions       
INFO[0000] * log_level = debug                          
INFO[0000] * database_type = PostgreSQL                 
INFO[0000] * kafka_in_topic = in-ms-arch-example        
INFO[0000] * http_port = 8081                           
INFO[0000] * kafka_partition = 0                        
INFO[0000] * database_url = localhost                   
INFO[0000] * database_username =                        
INFO[0000] * kafka_in_offset =                          
INFO[0000] * database_port = 5432                       
INFO[0000] * database_name =                            
INFO[0000] * service_name = ms-arch-example             
INFO[0000] * profile = dev                              
INFO[0000] * kafka_brokers = kafka:9092                 
INFO[0000] * kafka_in_consumer_group =                  
INFO[0000] * log_file = stdout 

Some of these variables are configured with default values.

To be possible to run this app from localhost we need add one entry on /etc/hosts like this one:

```
127.0.0.1   kafka
```

## Stress test

In the directory ./stress you can found one stress script.  This script uses Apache Benchmark (ab) to send 100K request against the service using the sale.json as body of each request.