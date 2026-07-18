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

// isRelational reports whether a database uses the SQL-migration workflow
// (versioned migrations/ + `make migrate-*`, DDL out of the constructor). The
// other backends are schemaless (mongodb/badgerdb), API-provisioned (dynamodb),
// or a different dialect (cassandra/CQL) and keep create-on-connect behavior.
func isRelational(database string) bool {
	switch database {
	case "postgres", "mysql", "mariadb", "sqlite":
		return true
	}
	return false
}

// hasObservability reports whether the project ships the OpenTelemetry stack
// (selected with -o otel). Only api/web-app/worker are eligible; the validator
// rejects it for hello-world.
func hasObservability(p *Project) bool {
	return p.Observability == "otel"
}

// hasAuth reports whether the project ships the JWT auth stack (selected with
// -a jwt). Only api is eligible; the validator rejects it for other types.
func hasAuth(p *Project) bool {
	return p.Auth == "jwt"
}

// hasOpenAPI reports whether the project ships the OpenAPI/Swagger stack (selected
// with -s swagger). REST api only — graphql has its own introspection, so it is
// excluded here as well as by the validator.
func hasOpenAPI(p *Project) bool {
	return p.OpenAPI == "swagger" && p.Platform != "graphql"
}

// hasCompose reports whether the project needs a docker-compose.yml + the docker-*
// make targets: true for a containerized backend (database/broker) OR observability
// (which contributes an otel-collector service).
func hasCompose(p *Project) bool {
	return hasContainerizedBackend(p) || hasObservability(p)
}

// observabilityPackages are the OpenTelemetry runtime modules fetched into the
// generated go.mod when -o otel is set (real deps, unlike the pinned go-run tools).
var observabilityPackages = []string{
	"go.opentelemetry.io/otel",
	"go.opentelemetry.io/otel/sdk",
	"go.opentelemetry.io/otel/sdk/metric",
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc",
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc",
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp",
	"go.opentelemetry.io/contrib/instrumentation/runtime",
}

// buildObservability emits the OpenTelemetry adapter and the local collector config
// and fetches the OTel modules. The OTEL_* env block (otel-env) and the compose
// collector service are handled by the caller's env/compose assembly.
func buildObservability(p *Project) {
	if !hasObservability(p) {
		return
	}
	fileGenerator([]string{"observability"}, p)
	fileGenerator([]string{"otel-collector"}, p)
	goGetPackages(p.Path, observabilityPackages)
}

// buildAuth emits the JWT auth adapter (bearer-validation middleware + demo /login
// issuer) and fetches golang-jwt. The JWT_*/AUTH_* env block (auth-env) is handled
// by the caller's env assembly.
func buildAuth(p *Project) {
	if !hasAuth(p) {
		return
	}
	fileGenerator([]string{"auth"}, p)
	goGetPackages(p.Path, []string{"github.com/golang-jwt/jwt/v5"})
}

// buildOpenAPI emits the templated OpenAPI document plus the Swagger-UI serving
// adapter, and fetches the embedded Swagger-UI package. The /openapi.yaml and
// /swagger routes are mounted by the cmd/server template (guarded on .OpenAPI).
func buildOpenAPI(p *Project) {
	if !hasOpenAPI(p) {
		return
	}
	fileGenerator([]string{"openapi-spec"}, p)
	fileGenerator([]string{"openapi-serve"}, p)
	goGetPackages(p.Path, []string{"github.com/swaggest/swgui/v5emb"})
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

	// Relational databases own their schema via versioned migrations rather than
	// constructor DDL; ship the initial create-users migration pair.
	if isRelational(database) {
		fileGenerator([]string{"migration-up"}, p)
		fileGenerator([]string{"migration-down"}, p)
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
