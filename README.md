# go-ms-arch-example

The aim of this project is to learn Go.

As a case study, I chose to develop something similar to what we do every day at work.

Therefore, this example must implement a microservice that aims to receive requests and send them to a kafka topic which will be consumed by another microservice and saved in a database.

I preferred not to use Go-kit, so I chose several frameworks to cover various parts of what I think is fundamental in creating a microservice.

Configuration: Viper
HTTP: Gin
Swagger: Gin-swagger
Logger: Logrus
Kafka: Kafka-go
ORM: Gorm

As a development architecture I chose something similar to Spring, with dependency injection by the constructor and an initial context concept.

