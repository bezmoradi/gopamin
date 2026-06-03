package templates

import (
	_ "embed"
)

//go:embed files/cmd/main.tmpl
var microserviceMain []byte

func MicroserviceMainTemplate() ([]byte, string) {
	return microserviceMain, "cmd/main.go"
}

//go:embed files/cmd/main-with-db.tmpl
var workerMainWithDb []byte

func WorkerMainWithDbTemplate() ([]byte, string) {
	return workerMainWithDb, "cmd/main.go"
}

//go:embed files/env/env-kafka.tmpl
var kafkaMicroserviceEnv []byte

func KafkaMicroserviceEnvTemplate() ([]byte, string) {
	return kafkaMicroserviceEnv, ".env"
}

//go:embed files/readme/readme-kafka.tmpl
var kafkaMicroserviceReadme []byte

func KafkaMicroserviceReadmeTemplate() ([]byte, string) {
	return kafkaMicroserviceReadme, "README.md"
}

//go:embed files/internal/adapters/brokers/kafka/kafka.tmpl
var kafkaMicroserviceBroker []byte

func KafkaMicroserviceBrokerTemplate() ([]byte, string) {
	return kafkaMicroserviceBroker, "internal/adapters/brokers/kafka/kafka.go"
}

//go:embed files/env/env-rabbitmq.tmpl
var rabbitmqMicroserviceEnv []byte

func RabbitmqMicroserviceEnvTemplate() ([]byte, string) {
	return rabbitmqMicroserviceEnv, ".env"
}

//go:embed files/readme/readme-rabbitmq.tmpl
var rabbitmqMicroserviceReadme []byte

func RabbitmqMicroserviceReadmeTemplate() ([]byte, string) {
	return rabbitmqMicroserviceReadme, "README.md"
}

//go:embed files/internal/adapters/brokers/rabbitmq/rabbitmq.tmpl
var rabbitmqMicroserviceBroker []byte

func RabbitmqMicroserviceBrokerTemplate() ([]byte, string) {
	return rabbitmqMicroserviceBroker, "internal/adapters/brokers/rabbitmq/rabbitmq.go"
}

//go:embed files/env/env-redis.tmpl
var redisMicroserviceEnv []byte

func RedisMicroserviceEnvTemplate() ([]byte, string) {
	return redisMicroserviceEnv, ".env"
}

//go:embed files/readme/readme-redis.tmpl
var redisMicroserviceReadme []byte

func RedisMicroserviceReadmeTemplate() ([]byte, string) {
	return redisMicroserviceReadme, "README.md"
}

//go:embed files/internal/adapters/brokers/redis/redis.tmpl
var redisMicroserviceBroker []byte

func RedisMicroserviceBrokerTemplate() ([]byte, string) {
	return redisMicroserviceBroker, "internal/adapters/brokers/redis/redis.go"
}
