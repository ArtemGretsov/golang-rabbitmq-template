package rabbitmq

import (
	"context"
	"log"
	"runtime/debug"
	"sync/atomic"

	"github.com/streadway/amqp"

	"github.com/ArtemGretsov/golang-rabbitmq-template/internal/shutdown"
)

const (
	lockState uint32 = iota
	unlockState
)

// Consumer - message listener with RabbitMQ.
type Consumer struct {
	ctx        context.Context
	shutdown   *shutdown.Module
	connection *Connection
	channel    *amqp.Channel
	handler    HandlerConsumer
	locker     uint32

	// Amqp's Consume arguments.
	queue     string
	consumer  string
	autoAck   bool
	exclusive bool
	noLocal   bool
	noWait    bool
	args      amqp.Table
}

// Run starts a message listener from RabbitMQ.
// The listener will not be started again if it is in the active phase or in the connect phase.
func (c *Consumer) Run() {
	if !atomic.CompareAndSwapUint32(&c.locker, unlockState, lockState) {
		return
	}

	go func() {
		defer atomic.StoreUint32(&c.locker, unlockState)
		var err error

		c.connection.connectMutex.RLock()
		c.channel, err = c.connection.connection.Channel()
		c.connection.connectMutex.RUnlock()

		if err != nil {
			log.Printf("channel connection error: %v\n", err)

			return
		}

		msgs, err := c.channel.Consume(
			c.queue,
			c.consumer,
			c.autoAck,
			c.exclusive,
			c.noLocal,
			c.noWait,
			c.args,
		)

		if err != nil {
			log.Printf("rabbitMQ queue connection error: %v\n", err)

			return
		}

		log.Printf("rabbitMQ queue connection successful (queue: %s)\n", c.queue)

		for delivery := range msgs {
			message := delivery

			go func() {
				c.shutdown.SafeRun()
				defer c.shutdown.SafeComplete()

				defer func() {
					if panicMessage := recover(); panicMessage != nil {
						log.Printf("unhandled exception: %s Stack: %s", panicMessage, debug.Stack())
					}
				}()

				log.Printf("receiving messages from RabbitMQ (queue: %s)\n", c.queue)

				err = c.handler(c.ctx, message)

				if err != nil {
					log.Printf("error processing message from RabbitMQ (queue: %s)\n", c.queue)
					return
				}

				log.Printf("successful message processing from RabbitMQ (queue: %s)\n", c.queue)
			}()
		}
	}()
}
