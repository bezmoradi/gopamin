package scaffolder

import "github.com/bezmoradi/gopamin/internal/project"

const BUILD_SUCCESS_MESSAGE = "project has been successfully created. Please consult the README.md file for instructions and guidance."

// Project is the data every template renders with.
type Project = project.Project

// generator renders a project's files and collects the modules it needs. The
// first error sticks: later calls do nothing, and the caller rolls the project
// back. Type builders add their own files plus their part of .env and README.md;
// buildCommon then emits everything the project types share.
type generator struct {
	p       *Project
	err     error
	modules []string
	env     []string // type-specific .env templates, concatenated after the common ones
	readme  []string // type-specific README.md templates, concatenated after the common ones
}

// file renders keys, concatenated in order, into their shared output path.
func (g *generator) file(keys ...string) {
	if g.err == nil {
		g.err = fileGenerator(keys, g.p)
	}
}

// need records modules to add to the project's go.mod at their pinned versions.
func (g *generator) need(modules ...string) {
	g.modules = append(g.modules, modules...)
}

// buildersMap dispatches a project type to its build function.
var buildersMap = map[string]func(g *generator){
	"api":         buildAPI,
	"web-app":     buildWebApp,
	"worker":      buildWorker,
	"hello-world": buildHelloWorld,
}

// ciRecipes maps a -c value to the template that renders its pipeline file.
var ciRecipes = map[string]string{
	"github": "github-ci",
	"gitlab": "gitlab-ci",
}

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

// databaseRecipe describes what varies between database adapters: the
// repository file and the modules it needs.
type databaseRecipe struct {
	repository string
	packages   []string
}

// buildCommon emits the files every project type shares. It runs after the type
// builder, so .env and README.md get the type's parts appended to the common ones.
func buildCommon(g *generator) {
	p := g.p
	g.file("gitignore")
	g.file("license")
	g.file("agents")
	if p.CI != "" {
		g.file(ciRecipes[p.CI])
	}

	g.file("logger-interface")
	g.file(loggerRecipes[p.Logger].logger)
	g.need(loggerRecipes[p.Logger].packages...)

	if p.Broker != "" {
		g.file(p.Broker + "-microservice-broker")
		g.need(workerRecipes[p.Broker]...)
	}
	if p.HasObservability() {
		g.file("observability")
		g.file("otel-collector")
		g.need(observabilityPackages...)
	}
	if p.HasAuth() {
		g.file("auth")
		g.file("hash-password-cmd")
		g.need("github.com/golang-jwt/jwt/v5")
	}
	if p.HasOpenAPI() {
		g.file("openapi-spec")
		g.file("openapi-serve")
		g.need("github.com/swaggest/swgui")
	}

	env := []string{"env"}
	env = append(env, g.env...)
	if p.Database != "" {
		env = append(env, p.Database+"-env")
	}
	if p.Broker != "" {
		env = append(env, p.Broker+"-microservice-env")
	}
	if p.HasObservability() {
		env = append(env, "otel-env")
	}
	if p.HasAuth() {
		env = append(env, "auth-env")
	}
	g.file(env...)

	readme := []string{"readme", p.Logger + "-readme"}
	if p.Database != "" {
		readme = append(readme, p.Database+"-readme")
	}
	g.file(append(readme, g.readme...)...)

	makefile := []string{"makefile"}
	if p.HasCompose() {
		makefile = append(makefile, "docker-makefile")
		g.file("docker-compose")
	}
	if p.IsRelational() {
		makefile = append(makefile, "migrate-makefile")
	}
	g.file(makefile...)

	for _, key := range []string{"configs", "configs-test", "tools", "tools-test", "arch-lint", "golangci"} {
		g.file(key)
	}
	if p.ProjectType != "hello-world" {
		g.file("dockerfile")
		g.file("dockerignore")
	}
}

// observabilityPackages are the OpenTelemetry runtime modules fetched into the
// generated go.mod when -o otel is set (real deps, unlike the go-run tools).
var observabilityPackages = []string{
	"go.opentelemetry.io/otel",
	"go.opentelemetry.io/otel/sdk",
	"go.opentelemetry.io/otel/sdk/metric",
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc",
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc",
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp",
	"go.opentelemetry.io/contrib/instrumentation/runtime",
}

// buildDatabase generates the user repository stack plus the chosen database's
// adapter. `database` is passed explicitly (rather than read from p.Database) so
// api/web-app can request the in-memory "mock" without mutating the Project.
func buildDatabase(g *generator, database string) {
	for _, key := range coreRepositoryFiles {
		// The user service is one template that grows publisher plumbing when
		// PublishesEvents() is true; all that's needed here is the outbound port
		// it depends on.
		if key == "user-service" && g.p.PublishesEvents() {
			g.file("event-publisher-interface")
		}
		g.file(key)
	}

	recipe := databaseRecipes[database]
	if recipe.repository != "" {
		g.file(recipe.repository)
	}

	// Relational databases own their schema via versioned migrations rather than
	// constructor DDL; ship the initial create-users migration pair.
	if g.p.IsRelational() {
		g.file("migration-up")
		g.file("migration-down")
	}

	g.need(recipe.packages...)
}

// loggerRecipe describes what varies between logger adapters.
type loggerRecipe struct {
	logger   string
	packages []string
}
