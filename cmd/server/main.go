package main

import (
	"context"
	"fmt"
	"log"

	"github.com/streadway/amqp"

	"github.com/ArtemGretsov/golang-rabbitmq-template/internal/config"
	"github.com/ArtemGretsov/golang-rabbitmq-template/internal/rabbitmq"
	"github.com/ArtemGretsov/golang-rabbitmq-template/internal/shutdown"
)

func main() {
	ctx := context.Background()
	configurator := config.NewConfigurator()

	shutdownModule := &shutdown.Module{Ctx: ctx}
	shutdownModule.Safe()

	rabbitMQModule := &rabbitmq.Module{Config: configurator, Ctx: ctx, Shutdown: shutdownModule}

	// examples
	appConfig := configurator.Get()
	connection := rabbitMQModule.Connection("name_connection", appConfig.RabbitMQConnectionURL)

	connection.RegisterConsumer(
		"first_queue",
		true,
		func(ctx context.Context, delivery amqp.Delivery) error {
			log.Printf("message: %s\n", delivery.Body)

			return nil
		})

	connection.RegisterConsumer(
		"second_queue",
		true,
		func(ctx context.Context, delivery amqp.Delivery) error {
			fmt.Printf("message: %s\n", delivery.Body)

			return nil
		})

	err := connection.Publish("", "first_queue", false, false, amqp.Publishing{Body: []byte("Hello")})

	if err != nil {
		log.Printf("publish message error (queue: first_queue): %v", err)
	}

	<-make(chan struct{})
}
