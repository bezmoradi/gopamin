package scaffolder

import "fmt"

// workerRecipes maps a broker to the modules its adapter needs. The file keys are
// derived uniformly from the broker name (<broker>-microservice-*).
var workerRecipes = map[string][]string{
	"kafka":    {"github.com/segmentio/kafka-go"},
	"rabbitmq": {"github.com/rabbitmq/amqp091-go"},
	"redis":    {"github.com/redis/go-redis/v9"},
}

func buildWorker(p *Project) error {
	broker := p.Broker
	packages := workerRecipes[broker]

	if err := buildLogger(p); err != nil {
		return err
	}
	if err := buildObservability(p); err != nil {
		return err
	}
	readme := []string{"readme", p.Logger + "-readme", broker + "-microservice-readme"}
	env := []string{"env", broker + "-microservice-env"}
	if p.Database != "" {
		if err := buildDatabase(p, p.Database); err != nil {
			return err
		}
		readme = append(readme, p.Database+"-readme")
		env = append(env, p.Database+"-env")
	}
	if hasObservability(p) {
		env = append(env, "otel-env")
	}

	fileGenerator(readme, p)
	fileGenerator(env, p)
	makefile := []string{"makefile", "docker-makefile"}
	if isRelational(p.Database) {
		makefile = append(makefile, "migrate-makefile")
	}
	fileGenerator(makefile, p)
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

	if err := goGetPackages(p.Path, packages); err != nil {
		return err
	}

	fmt.Printf("%v "+BUILD_SUCCESS_MESSAGE+"\n", p.Name)
	return nil
}
