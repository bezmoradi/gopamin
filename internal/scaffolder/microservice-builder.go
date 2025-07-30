package scaffolder

import "fmt"

type kafkaMicroserviceBuilder struct {
	project *Project
}

func (b *kafkaMicroserviceBuilder) build() {
	loggerSelector(b.project)
	fileGenerator([]string{"readme", b.project.Logger + "-readme", "kafka-microservice-readme"}, b.project)
	fileGenerator([]string{"env", "kafka-microservice-env"}, b.project)
	fileGenerator([]string{"makefile", "kafka-microservice-makefile"}, b.project)
	fileGenerator([]string{"kafka-microservice-docker-compose"}, b.project)
	fileGenerator([]string{"kafka-microservice-broker"}, b.project)
	fileGenerator([]string{"configs"}, b.project)
	fileGenerator([]string{"configs-test"}, b.project)
	fileGenerator([]string{"tools"}, b.project)
	fileGenerator([]string{"tools-test"}, b.project)
	fileGenerator([]string{"message"}, b.project)
	fileGenerator([]string{"broker-service-interface"}, b.project)
	fileGenerator([]string{"broker-service"}, b.project)
	fileGenerator([]string{"message-broker-interface"}, b.project)
	fileGenerator([]string{"microservice-main"}, b.project)

	goGetPackages(b.project.Path, []string{"github.com/segmentio/kafka-go"})

	fmt.Printf("%v "+BUILD_SUCCESS_MESSAGE+"\n", b.project.Name)
}

type rabbitmqMicroserviceBuilder struct {
	project *Project
}

func (b *rabbitmqMicroserviceBuilder) build() {
	loggerSelector(b.project)
	fileGenerator([]string{"readme", b.project.Logger + "-readme", "rabbitmq-microservice-readme"}, b.project)
	fileGenerator([]string{"env", "rabbitmq-microservice-env"}, b.project)
	fileGenerator([]string{"makefile", "rabbitmq-microservice-makefile"}, b.project)
	fileGenerator([]string{"rabbitmq-microservice-docker-compose"}, b.project)
	fileGenerator([]string{"rabbitmq-microservice-broker"}, b.project)
	fileGenerator([]string{"configs"}, b.project)
	fileGenerator([]string{"configs-test"}, b.project)
	fileGenerator([]string{"tools"}, b.project)
	fileGenerator([]string{"tools-test"}, b.project)
	fileGenerator([]string{"message"}, b.project)
	fileGenerator([]string{"broker-service-interface"}, b.project)
	fileGenerator([]string{"broker-service"}, b.project)
	fileGenerator([]string{"message-broker-interface"}, b.project)
	fileGenerator([]string{"microservice-main"}, b.project)

	goGetPackages(b.project.Path, []string{"github.com/streadway/amqp"})

	fmt.Printf("%v "+BUILD_SUCCESS_MESSAGE+"\n", b.project.Name)
}

type redisMicroserviceBuilder struct {
	project *Project
}

func (b *redisMicroserviceBuilder) build() {
	loggerSelector(b.project)
	fileGenerator([]string{"readme", b.project.Logger + "-readme", "redis-microservice-readme"}, b.project)
	fileGenerator([]string{"env", "redis-microservice-env"}, b.project)
	fileGenerator([]string{"makefile", "redis-microservice-makefile"}, b.project)
	fileGenerator([]string{"redis-microservice-docker-compose"}, b.project)
	fileGenerator([]string{"redis-microservice-broker"}, b.project)
	fileGenerator([]string{"configs"}, b.project)
	fileGenerator([]string{"configs-test"}, b.project)
	fileGenerator([]string{"tools"}, b.project)
	fileGenerator([]string{"tools-test"}, b.project)
	fileGenerator([]string{"message"}, b.project)
	fileGenerator([]string{"broker-service-interface"}, b.project)
	fileGenerator([]string{"broker-service"}, b.project)
	fileGenerator([]string{"message-broker-interface"}, b.project)
	fileGenerator([]string{"microservice-main"}, b.project)

	goGetPackages(b.project.Path, []string{"github.com/redis/go-redis/v9"})

	fmt.Printf("%v "+BUILD_SUCCESS_MESSAGE+"\n", b.project.Name)
}

type microserviceBuilderFactory func(p *Project) boilerplateBuilder

var microserviceBuilderMap = map[string]microserviceBuilderFactory{
	"kafka": func(p *Project) boilerplateBuilder {
		return &kafkaMicroserviceBuilder{p}
	},
	"rabbitmq": func(p *Project) boilerplateBuilder {
		return &rabbitmqMicroserviceBuilder{p}
	},
	"redis": func(p *Project) boilerplateBuilder {
		return &redisMicroserviceBuilder{p}
	},
}
