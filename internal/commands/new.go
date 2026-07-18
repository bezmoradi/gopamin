package commands

import (
	"os"

	"github.com/bezmoradi/gopamin/internal/scaffolder"
	"github.com/spf13/cobra"
)

var (
	database      string
	framework     string
	broker        string
	name          string
	projectType   string
	logger        string
	ci            string
	observability string
	auth          string
	openapi       string
)

var newCmd = &cobra.Command{
	Example: `To create a simple hello world app, run the following command:
 gopamin new -n HelloWorld -t hello-world -l slog
To create a hello world app with MySQL integration, run the following command:
 gopamin new -n HelloWorld -t hello-world -d mysql -l slog
To create a web application using the built-in http package and MySQL support, run the following command:
 gopamin new -n HelloWorld -t web-app -f http -d mysql -l slog
To create a RESTful API using the echo framework, run the following command:
 gopamin new -n HelloWorld -t api -f echo -l slog
To create a RESTful API using the built-in http package and MongoDB integration, run the following command:
 gopamin new -n HelloWorld -t api -f http -d mongodb -l slog
To create a worker (a broker-driven service with no HTTP server) using the Kafka message broker, run the following command:
 gopamin new -n HelloWorld -t worker -b kafka -l slog`,
	Use:   "new",
	Short: "Create a new project",
	Run: func(cmd *cobra.Command, args []string) {
		// argsValidator prints a specific message for the first rule it fails.
		// Exit non-zero on rejection so scripts and CI can detect the failure
		// (the flag values are invalid input, not a successful no-op).
		if !argsValidator() {
			os.Exit(1)
		}
		scaffolder.New(projectType, framework, broker, name, database, logger, ci, observability, auth, openapi)
	},
}

func init() {
	newCmd.Flags().StringVarP(&name, "name", "n", "", `Name of project (Keep in mind the all white spaces will be replaced by dash).
If you want to use space-separated words, place them inside double quotes like "my demo app" then it will be converted to "my-demo-app".`)

	newCmd.Flags().StringVarP(&projectType, "type", "t", "", `Type of the project. Available types are:
 - hello-world (A simple "Hello World" app. A database can also be added with the "-d" flag).
 - web-app (A server-rendered HTML app. Requires the "-f" flag; "-d" is optional).
 - api (An HTTP API. Requires the "-f" flag; "-d" is optional).
 - worker (A broker-driven service with no HTTP server. Requires the "-b" flag).`)

	newCmd.Flags().StringVarP(&framework, "framework", "f", "", `The HTTP framework, required for "api" and "web-app" projects.
Available values for the "api" type are:
 - echo
 - chi
 - gin
 - httprouter
 - gorilla
 - http (the built-in net/http package will be used)
 - graphql
Available values for the "web-app" type are the same minus graphql.`)

	newCmd.Flags().StringVarP(&broker, "broker", "b", "", `The message broker. Required for the "worker" type; optional for "api"/"web-app" (adds a producer on writes + a consumer loop alongside the HTTP server). Available values are:
 - kafka
 - rabbitmq
 - redis`)

	newCmd.Flags().StringVarP(&database, "database", "d", "", `Type of the database. Optional. Available values are:
 - mysql
 - mariadb
 - cassandra
 - postgres
 - mongodb
 - sqlite
 - dynamodb
 - badgerdb`)

	newCmd.Flags().StringVarP(&logger, "logger", "l", "", `Type of the logger (required for every project type). Available values are:
 - log
 - slog
 - logrus
 - zap`)

	newCmd.Flags().StringVarP(&ci, "ci", "c", "", `Optional CI pipeline to generate. Available values are:
 - github (.github/workflows/ci.yml)
 - gitlab (.gitlab-ci.yml)
The pipeline runs build, test, "make lint", and "make arch-check".`)

	newCmd.Flags().StringVarP(&observability, "observability", "o", "", `Optional observability stack to generate (for "api", "web-app", and "worker"). Available values are:
 - otel (OpenTelemetry: traces + Go runtime metrics, OTLP export to a collector; HTTP
   request spans for api/web-app, a per-message span for a worker)
It adds an "internal/adapters/observability" package, an OTEL_* block in .env, and a
collector service in docker-compose.yml. Configure it via the standard OTEL_* env vars.`)

	newCmd.Flags().StringVarP(&auth, "auth", "a", "", `Optional authentication to generate (for "api" projects only). Available values are:
 - jwt (bearer-token JWT: a validation middleware protecting every app route, plus a
   demo POST /login token issuer, using golang-jwt/v5 with HS256 + algorithm pinning)
It adds an "internal/adapters/handlers/auth" package and a JWT_*/AUTH_* block in .env.`)

	newCmd.Flags().StringVarP(&openapi, "openapi", "s", "", `Optional API documentation to generate (for REST "api" projects only). Available values are:
 - swagger (a templated openapi.yaml describing the generated API, plus embedded
   Swagger UI served at /swagger with the spec at /openapi.yaml, via swaggest/swgui)
It adds "openapi.yaml" + an OpenAPI serving handler to "internal/adapters/api".`)

	rootCmd.AddCommand(newCmd)
}
