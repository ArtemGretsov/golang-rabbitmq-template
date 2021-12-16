WIP

# golang-rabbitmq-template

### features:
1. Simple application configuration via json (default.json / local.json) and env variables.
2. Stop the application safely. Your handlers will complete their execution before stopping the application.
3. Reconnecting to RabbitMQ in case of error.
4. Reconnecting to the queue RabbitMQ in case of an error.
5. The number of RabbitMQ connections is not limited.
6. You don't have to wait for a connection to RabbitMQ and register your handlers instantly. When a RabbitMQ connection is opened, all handlers will be initialized.


### Run
```bash
docker-compose up -d
make dev
 ```


