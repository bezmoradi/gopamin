package scaffolder

// workerRecipes maps a broker to the modules its adapter needs. The file keys are
// derived uniformly from the broker name (<broker>-microservice-*).
var workerRecipes = map[string][]string{
	"kafka":    {"github.com/segmentio/kafka-go"},
	"rabbitmq": {"github.com/rabbitmq/amqp091-go"},
	"redis":    {"github.com/redis/go-redis/v9"},
}

func buildWorker(g *generator) {
	g.readme = append(g.readme, g.p.Broker+"-microservice-readme")
	if g.p.Database != "" {
		buildDatabase(g, g.p.Database)
	} else {
		// Events are about users, so a worker without a database still needs
		// the user domain.
		g.file("user")
		g.file("user-test")
	}
	g.file("worker-main")
}
