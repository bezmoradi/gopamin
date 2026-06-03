package scaffolder

const BUILD_SUCCESS_MESSAGE = "project has been successfully created. Please consult the README.md file for instructions and guidance."

// coreRepositoryFiles are the domain/ports/services/mock files generated for any
// project that ships the user repository stack (api, web-app, and hello-world
// with a database).
var coreRepositoryFiles = []string{
	"mock-repository",
	"user-test",
	"user",
	"user-repository-interface",
	"user-service-interface",
	"user-service",
	"user-service-test",
}

// databaseRecipe describes what varies between database adapters: the repository
// file, whether the database runs in a container (and so contributes a service to
// the merged docker-compose.yml), and the modules to fetch.
type databaseRecipe struct {
	repository    string
	containerized bool
	packages      []string
}

// usesBrokerPublisher reports whether the project's user service should publish
// events on writes: true for an api/web-app that selected a broker.
func usesBrokerPublisher(p *Project) bool {
	return (p.ProjectType == "api" || p.ProjectType == "web-app") && p.Broker != ""
}

// hasContainerizedBackend reports whether the project needs a docker-compose.yml
// and the docker make targets: true when it uses a containerized database and/or
// a broker. (sqlite and badgerdb are embedded, so they need no container.)
func hasContainerizedBackend(p *Project) bool {
	return databaseRecipes[p.Database].containerized || p.Broker != ""
}

// buildDatabase generates the user repository stack plus the chosen database's
// adapter. `database` is passed explicitly (rather than read from p.Database) so
// callers can request the in-memory "mock" without mutating the Project.
func buildDatabase(p *Project, database string) {
	for _, key := range coreRepositoryFiles {
		// An api/web-app with a broker replaces the plain user service with the
		// publisher-aware variant (and its EventPublisher port) in place, so the
		// base service file is never written and then overwritten.
		if key == "user-service" && usesBrokerPublisher(p) {
			fileGenerator([]string{"event-publisher-interface"}, p)
			fileGenerator([]string{"user-service-broker"}, p)
			continue
		}
		fileGenerator([]string{key}, p)
	}

	recipe := databaseRecipes[database]
	if recipe.repository != "" {
		fileGenerator([]string{recipe.repository}, p)
	}

	goGetPackages(p.Path, recipe.packages)
}

// loggerRecipe describes what varies between logger adapters.
type loggerRecipe struct {
	logger   string
	packages []string
}

func buildLogger(p *Project) {
	recipe := loggerRecipes[p.Logger]
	fileGenerator([]string{recipe.logger}, p)
	fileGenerator([]string{"logger-interface"}, p)
	goGetPackages(p.Path, recipe.packages)
}

// buildersMap dispatches a project type to its build function.
var buildersMap = map[string]func(p *Project){
	"api":         buildAPI,
	"web-app":     buildWebApp,
	"worker":      buildWorker,
	"hello-world": buildHelloWorld,
}
