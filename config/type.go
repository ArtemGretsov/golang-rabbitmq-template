package configtype

type Config struct {
	RabbitMQConnectionURL  string `mapstructure:"rabbitmq_connection_url"`
	RabbitMQReconnectDelay int    `mapstructure:"rabbitmq_reconnect_delay"`
	RabbitMQReconsumeDelay int    `mapstructure:"rabbitmq_reconsume_delay"`
}
