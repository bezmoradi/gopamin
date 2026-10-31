// Package project describes the project a user asked gopamin to generate: the
// selected flags, and the properties derived from them that both the scaffolder
// and the templates branch on. Keeping each predicate here, as a method the
// templates can call, means the builder and the files it renders can't disagree.
package project

import "strings"

// Project is the data every template renders with.
type Project struct {
	Database    string
	Platform    string
	Broker      string
	Name        string
	Path        string
	ProjectType string
	Logger      string
	// GoVersion is the pinned Go version exposed to templates so the Dockerfile's
	// build image never drifts from the go.mod `go` directive.
	GoVersion string
	// CI is the optional CI platform ("github" or "gitlab", empty for none)
	// selected with -c; it decides which pipeline file, if any, is emitted.
	CI string
	// Observability is the optional observability stack ("otel", empty for none)
	// selected with -o; when set, templates emit the OTel adapter, the OTLP
	// instrumentation wrap, the OTEL_* env block, and a collector in compose.
	Observability string
	// Auth is the optional authentication scheme ("jwt", empty for none) selected
	// with -a (api only); when set, templates emit the auth adapter (bearer
	// validation middleware + demo /login issuer) and a JWT_*/AUTH_* env block.
	Auth string
	// OpenAPI is the optional API-docs stack ("swagger", empty for none) selected
	// with -s (REST api only); when set, templates emit a templated openapi.yaml
	// plus embedded Swagger UI served at /swagger.
	OpenAPI string
	// Secrets are random per-project values the scaffolder generates and the
	// templates write into .env, so no project ships a known default credential.
	Secrets Secrets
}

// Secrets holds the generated values. Only what belongs in .env lives here: the
// demo login password itself is printed once by the scaffolder and never stored.
type Secrets struct {
	JWTSecret        string
	AuthPasswordHash string
}

// containerizedDatabases run as a service in docker-compose.yml. sqlite and
// badgerdb are embedded in the binary, so they need no container.
var containerizedDatabases = map[string]bool{
	"mysql":     true,
	"mariadb":   true,
	"postgres":  true,
	"mongodb":   true,
	"dynamodb":  true,
	"cassandra": true,
}

// PublishesEvents reports whether the user service publishes events on writes:
// true for an api/web-app that selected a broker. (A worker sets Broker too, but
// consumes rather than publishes.)
func (p Project) PublishesEvents() bool {
	return (p.ProjectType == "api" || p.ProjectType == "web-app") && p.Broker != ""
}

// HasHTTP reports whether the project runs an HTTP server.
func (p Project) HasHTTP() bool {
	return p.ProjectType == "api" || p.ProjectType == "web-app"
}

// Slug is the last element of the module path, such as "svc" for
// github.com/acme/svc. Topic, queue and consumer-group names derive from it, so
// two projects on one broker don't share them by default.
func (p Project) Slug() string {
	return p.Name[strings.LastIndex(p.Name, "/")+1:]
}

// HasContainerizedBackend reports whether the project uses a containerized
// database and/or a broker.
func (p Project) HasContainerizedBackend() bool {
	return containerizedDatabases[p.Database] || p.Broker != ""
}

// HasCompose reports whether the project ships a docker-compose.yml and the
// docker-* make targets: a containerized backend, or observability (which adds
// an otel-collector service).
func (p Project) HasCompose() bool {
	return p.HasContainerizedBackend() || p.HasObservability()
}

// ComposeServices names the docker-compose.yml services, joined by " + ", for
// the generated docs.
func (p Project) ComposeServices() string {
	var services []string
	if containerizedDatabases[p.Database] {
		services = append(services, p.Database)
	}
	if p.Broker != "" {
		services = append(services, p.Broker)
	}
	if p.HasObservability() {
		services = append(services, "otel-collector")
	}
	return strings.Join(services, " + ")
}

// IsRelational reports whether the database uses the SQL-migration workflow
// (versioned migrations/ + `make migrate-*`, no DDL in the constructor). The
// other backends are schemaless (mongodb/badgerdb), API-provisioned (dynamodb),
// or a different dialect (cassandra/CQL) and create what they need on connect.
func (p Project) IsRelational() bool {
	switch p.Database {
	case "postgres", "mysql", "mariadb", "sqlite":
		return true
	}
	return false
}

// HasObservability reports whether the project ships the OpenTelemetry stack
// (-o otel).
func (p Project) HasObservability() bool {
	return p.Observability == "otel"
}

// HasAuth reports whether the project ships the JWT auth stack (-a jwt).
func (p Project) HasAuth() bool {
	return p.Auth == "jwt"
}

// HasOpenAPI reports whether the project ships the OpenAPI/Swagger stack
// (-s swagger). REST api only: graphql describes itself.
func (p Project) HasOpenAPI() bool {
	return p.OpenAPI == "swagger" && p.Platform != "graphql"
}
