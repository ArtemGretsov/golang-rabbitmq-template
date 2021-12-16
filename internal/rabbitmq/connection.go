package rabbitmq

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/ArtemGretsov/golang-rabbitmq-template/internal/config"
	"github.com/ArtemGretsov/golang-rabbitmq-template/internal/shutdown"
)

const consumerCheckingDelay = 10 * time.Second

// HandlerConsumer - message handler.
type HandlerConsumer func(context.Context, amqp.Delivery) error

// Connection - connection RabbitMQ.
type Connection struct {
	ctx      context.Context
	config   config.Configurator
	shutdown *shutdown.Module

	connectMutex   sync.RWMutex
	consumersMutex sync.RWMutex

	connection *amqp.Connection
	consumers  []*Consumer

	Name       string
	URL        string
	PrintedURL string
}

// Publish publishes a message to the specified queue or exchange.
func (c *Connection) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	var (
		channel *amqp.Channel
		err     error
	)

	func() {
		c.connectMutex.RLock()
		defer c.connectMutex.RUnlock()

		channel, err = c.connection.Channel()
	}()

	defer channel.Close()

	if err != nil {
		return errors.Wrapf(err, "error opening RabbitMQ channel to post message (queue: %s)", key)
	}

	err = channel.Publish(exchange, key, mandatory, immediate, msg)

	if err != nil {
		return errors.Wrapf(err, "error posting message to RabbitMQ (queue: %s)", key)
	}

	return nil
}

// RegisterConsumer registers a new listener with default settings.
// If the connection is not open at the time of calling RegisterConsumer,
// the listener will be registered after successfully connecting to the queue.
func (c *Connection) RegisterConsumer(queue string, autoAck bool, handler HandlerConsumer) {
	c.consumersMutex.Lock()
	defer c.consumersMutex.Unlock()

	consumer := &Consumer{
		shutdown:   c.shutdown,
		ctx:        c.ctx,
		handler:    handler,
		connection: c,
		queue:      queue,
		locker:     unlockState,
		autoAck:    autoAck,
		exclusive:  false,
		noLocal:    false,
		noWait:     false,
		args:       nil,
	}

	consumer.Run()

	c.consumers = append(c.consumers, consumer)
}

func (c *Connection) connect() {
	var (
		err    error
		isInit bool
	)

	appConfig := c.config.Get()
	delayReconnect := time.Duration(appConfig.RabbitMQReconnectDelay) * time.Second

	for !isInit {
		c.connection, err = amqp.Dial(c.URL)

		if err != nil {
			log.Printf("rabbitMQ connection error: %v", err)
			time.Sleep(delayReconnect)
			continue
		}

		log.Println("connect to RabbitMQ")

		isInit = true
	}
}

func (c *Connection) check() {
	go func() {
		closeConnectionNotify := c.connection.NotifyClose(make(chan *amqp.Error))
		<-closeConnectionNotify

		func() {
			c.connectMutex.Lock()
			defer c.connectMutex.Unlock()

			c.connect()
		}()

		c.runAllConsumers()
	}()

}

func (c *Connection) checkConsumers() {
	go func() {
		for {
			ctx, stop := shutdown.Subscribe(c.ctx)
			select {
			case <-ctx.Done():
				stop()
				c.closeAllConsumers()
				return

			case <-time.After(consumerCheckingDelay):
				c.runAllConsumers()
			}
		}
	}()
}

func (c *Connection) runAllConsumers() {
	c.consumersMutex.RLock()
	defer c.consumersMutex.RUnlock()

	for _, consumer := range c.consumers {
		consumer.Run()
	}
}

func (c *Connection) closeAllConsumers() {
	c.consumersMutex.RLock()
	defer c.consumersMutex.RUnlock()

	for _, consumer := range c.consumers {
		_ = consumer.channel.Close()
	}
}
