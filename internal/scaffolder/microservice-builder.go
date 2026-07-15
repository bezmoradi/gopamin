package scaffolder

import "fmt"

// workerRecipes maps a broker to the modules its adapter needs. The file keys are
// derived uniformly from the broker name (<broker>-microservice-*).
var workerRecipes = map[string][]string{
	"kafka":    {"github.com/segmentio/kafka-go"},
	"rabbitmq": {"github.com/rabbitmq/amqp091-go"},
	"redis":    {"github.com/redis/go-redis/v9"},
}

func buildWorker(p *Project) {
	broker := p.Broker
	packages := workerRecipes[broker]

	buildLogger(p)
	readme := []string{"readme", p.Logger + "-readme", broker + "-microservice-readme"}
	env := []string{"env", broker + "-microservice-env"}
	if p.Database != "" {
		buildDatabase(p, p.Database)
		readme = append(readme, p.Database+"-readme")
		env = append(env, p.Database+"-env")
	}

	fileGenerator(readme, p)
	fileGenerator(env, p)
	fileGenerator([]string{"makefile", "docker-makefile"}, p)
	fileGenerator([]string{"docker-compose"}, p)
	fileGenerator([]string{broker + "-microservice-broker"}, p)
	fileGenerator([]string{"configs"}, p)
	fileGenerator([]string{"configs-test"}, p)
	fileGenerator([]string{"tools"}, p)
	fileGenerator([]string{"tools-test"}, p)
	fileGenerator([]string{"arch-lint"}, p)
	fileGenerator([]string{"golangci"}, p)
	fileGenerator([]string{"dockerfile"}, p)
	fileGenerator([]string{"dockerignore"}, p)
	fileGenerator([]string{"message"}, p)

	if p.Database == "" {
		fileGenerator([]string{"microservice-main"}, p)
	} else {
		fileGenerator([]string{"worker-main-with-db"}, p)
	}

	goGetPackages(p.Path, packages)

	fmt.Printf("%v "+BUILD_SUCCESS_MESSAGE+"\n", p.Name)
}
