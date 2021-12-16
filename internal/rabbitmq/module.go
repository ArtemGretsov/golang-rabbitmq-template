package rabbitmq

import (
	"context"
	"sync"

	"github.com/ArtemGretsov/golang-rabbitmq-template/internal/config"
	"github.com/ArtemGretsov/golang-rabbitmq-template/internal/shutdown"
)

// Module - int module.
type Module struct {
	Ctx      context.Context
	Config   config.Configurator
	Shutdown *shutdown.Module

	storage sync.Map
}

// Connection creates a new or returns an already created connection by amqp url
// Runs a background task to check the connection and restore it in case of disconnection.
// Runs a background task to check listeners and restore them in case of error.
func (m *Module) Connection(name, url string) *Connection {
	connection, ok := m.storage.LoadOrStore(url, &Connection{
		shutdown: m.Shutdown,
		config:   m.Config,
		ctx:      m.Ctx,
		URL:      url,
		Name:     name,
	})

	connectionInstance := connection.(*Connection)

	if !ok {
		connectionInstance.connectMutex.Lock()

		go func() {
			connectionInstance.connect()
			connectionInstance.connectMutex.Unlock()

			connectionInstance.check()
			connectionInstance.checkConsumers()
		}()
	}

	return connection.(*Connection)
}
