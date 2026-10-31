package templates

import (
	"embed"
	"fmt"
)

// files holds every template. The directories group templates by the project
// shape that uses them; they are data, not Go packages.
//
//go:embed api/files common/files database/files hello-world/files microservice/files web-app/files
var files embed.FS

// entry is one registry entry: where the template lives in files, and the
// path it renders to inside the generated project. Several entries may share an
// output path (README.md, .env, Makefile): the scaffolder concatenates them.
type entry struct {
	src string
	out string
}

// registry maps a template key, as used by the scaffolder's recipes, to its
// template.
var registry = map[string]entry{
	"env":                       {"common/files/env.tmpl", ".env"},
	"readme":                    {"common/files/readme/readme.tmpl", "README.md"},
	"log-readme":                {"common/files/readme/readme-log.tmpl", "README.md"},
	"slog-readme":               {"common/files/readme/readme-slog.tmpl", "README.md"},
	"logrus-readme":             {"common/files/readme/readme-logrus.tmpl", "README.md"},
	"zap-readme":                {"common/files/readme/readme-zap.tmpl", "README.md"},
	"makefile":                  {"common/files/makefile.tmpl", "Makefile"},
	"docker-compose":            {"common/files/docker-compose.tmpl", "docker-compose.yml"},
	"docker-makefile":           {"common/files/docker-makefile.tmpl", "Makefile"},
	"migrate-makefile":          {"common/files/migrate-makefile.tmpl", "Makefile"},
	"migration-up":              {"common/files/migration-up.tmpl", "migrations/000001_create_users.up.sql"},
	"migration-down":            {"common/files/migration-down.tmpl", "migrations/000001_create_users.down.sql"},
	"migrations-embed":          {"database/files/migrations/embed.tmpl", "migrations/migrations.go"},
	"observability":             {"common/files/internal/adapters/observability/otel.tmpl", "internal/adapters/observability/otel.go"},
	"otel-collector":            {"common/files/otel-collector.tmpl", "configs/otel-collector.yaml"},
	"otel-env":                  {"common/files/otel-env.tmpl", ".env"},
	"auth":                      {"common/files/internal/adapters/handlers/auth/auth.tmpl", "internal/adapters/handlers/auth/auth.go"},
	"auth-env":                  {"common/files/auth-env.tmpl", ".env"},
	"dockerfile":                {"common/files/dockerfile.tmpl", "Dockerfile"},
	"dockerignore":              {"common/files/dockerignore.tmpl", ".dockerignore"},
	"gitignore":                 {"common/files/gitignore.tmpl", ".gitignore"},
	"license":                   {"common/files/license.tmpl", "LICENSE"},
	"agents":                    {"common/files/AGENTS.tmpl", "AGENTS.md"},
	"arch-lint":                 {"common/files/arch-lint.tmpl", ".go-arch-lint.yml"},
	"golangci":                  {"common/files/golangci.tmpl", ".golangci.yml"},
	"github-ci":                 {"common/files/github-ci.tmpl", ".github/workflows/ci.yml"},
	"gitlab-ci":                 {"common/files/gitlab-ci.tmpl", ".gitlab-ci.yml"},
	"health-handler":            {"common/files/internal/adapters/handlers/health/health.tmpl", "internal/adapters/handlers/health/health.go"},
	"middleware":                {"common/files/internal/adapters/handlers/middleware/middleware.tmpl", "internal/adapters/handlers/middleware/middleware.go"},
	"mock-repository":           {"common/files/internal/adapters/repositories/mock/repository.tmpl", "internal/adapters/repositories/mock/repository.go"},
	"user-test":                 {"common/files/internal/core/user/domain/user_test.tmpl", "internal/core/user/domain/user_test.go"},
	"user":                      {"common/files/internal/core/user/domain/user.tmpl", "internal/core/user/domain/user.go"},
	"user-repository-interface": {"common/files/internal/core/user/ports/user-repository.interface.tmpl", "internal/core/user/ports/user-repository.interface.go"},
	"user-service-interface":    {"common/files/internal/core/user/ports/user-service.interface.tmpl", "internal/core/user/ports/user-service.interface.go"},
	"user-service":              {"common/files/internal/core/user/services/user.service.tmpl", "internal/core/user/services/user.service.go"},
	"user-service-test":         {"common/files/internal/core/user/services/user.service_test.tmpl", "internal/core/user/services/user.service_test.go"},
	"logger-interface":          {"common/files/internal/core/shared/logger.interface.tmpl", "internal/core/shared/logger.interface.go"},
	"log-logger":                {"common/files/internal/adapters/loggers/log/logger.tmpl", "internal/adapters/loggers/log/logger.go"},
	"slog-logger":               {"common/files/internal/adapters/loggers/slog/logger.tmpl", "internal/adapters/loggers/slog/logger.go"},
	"logrus-logger":             {"common/files/internal/adapters/loggers/logrus/logger.tmpl", "internal/adapters/loggers/logrus/logger.go"},
	"zap-logger":                {"common/files/internal/adapters/loggers/zap/logger.tmpl", "internal/adapters/loggers/zap/logger.go"},
	"api-errors":                {"common/files/internal/adapters/api/errors.tmpl", "internal/adapters/api/errors.go"},
	"api-response":              {"common/files/internal/adapters/api/response.tmpl", "internal/adapters/api/response.go"},
	"api-user":                  {"common/files/internal/adapters/api/user.tmpl", "internal/adapters/api/user.go"},
	"openapi-spec":              {"common/files/internal/adapters/api/openapi.tmpl", "internal/adapters/api/openapi.yaml"},
	"openapi-serve":             {"common/files/internal/adapters/api/openapi-serve.tmpl", "internal/adapters/api/openapi.go"},
	"configs":                   {"common/files/configs/configs.tmpl", "configs/configs.go"},
	"configs-test":              {"common/files/configs/configs_test.tmpl", "configs/configs_test.go"},
	"user-events":               {"common/files/internal/core/user/domain/events.tmpl", "internal/core/user/domain/events.go"},
	"events":                    {"common/files/internal/adapters/events/events.tmpl", "internal/adapters/events/events.go"},
	"event-publisher-interface": {"common/files/internal/core/user/ports/event-publisher.interface.tmpl", "internal/core/user/ports/event-publisher.interface.go"},
	"tools":                     {"common/files/tools/tools.tmpl", "tools/tools.go"},
	"tools-test":                {"common/files/tools/tools_test.tmpl", "tools/tools_test.go"},
	"http-main":                 {"common/files/cmd/http-main.tmpl", "cmd/main.go"},
	"http-server":               {"common/files/cmd/server/server.tmpl", "cmd/server/server.go"},
	"hash-password-cmd":         {"common/files/cmd/hashpassword/main.tmpl", "cmd/hashpassword/main.go"},

	"hello-world-main":           {"hello-world/files/cmd/main.tmpl", "cmd/main.go"},
	"hello-world-readme-with-db": {"hello-world/files/readme/readme-with-db.tmpl", "README.md"},
	"hello-world-main-test":      {"hello-world/files/cmd/main_test.tmpl", "cmd/main_test.go"},
	"hello-world-main-with-db":   {"hello-world/files/cmd/main-with-db.tmpl", "cmd/main.go"},

	"api-http-routes":       {"api/files/internal/adapters/handlers/http/routes.tmpl", "internal/adapters/handlers/http/routes.go"},
	"api-http-users":        {"api/files/internal/adapters/handlers/std/users.tmpl", "internal/adapters/handlers/http/users.go"},
	"api-readme":            {"api/files/readme/rest-readme.tmpl", "README.md"},
	"graphql-readme":        {"api/files/readme/graphql-readme.tmpl", "README.md"},
	"api-env":               {"api/files/env.tmpl", ".env"},
	"api-chi-routes":        {"api/files/internal/adapters/handlers/chi/routes.tmpl", "internal/adapters/handlers/chi/routes.go"},
	"api-chi-users":         {"api/files/internal/adapters/handlers/std/users.tmpl", "internal/adapters/handlers/chi/users.go"},
	"api-echo-routes":       {"api/files/internal/adapters/handlers/echo/routes.tmpl", "internal/adapters/handlers/echo/routes.go"},
	"api-echo-users":        {"api/files/internal/adapters/handlers/echo/users.tmpl", "internal/adapters/handlers/echo/users.go"},
	"api-gin-routes":        {"api/files/internal/adapters/handlers/gin/routes.tmpl", "internal/adapters/handlers/gin/routes.go"},
	"api-gin-users":         {"api/files/internal/adapters/handlers/gin/users.tmpl", "internal/adapters/handlers/gin/users.go"},
	"api-gorilla-routes":    {"api/files/internal/adapters/handlers/gorilla/routes.tmpl", "internal/adapters/handlers/gorilla/routes.go"},
	"api-gorilla-users":     {"api/files/internal/adapters/handlers/std/users.tmpl", "internal/adapters/handlers/gorilla/users.go"},
	"api-httprouter-routes": {"api/files/internal/adapters/handlers/httprouter/routes.tmpl", "internal/adapters/handlers/httprouter/routes.go"},
	"api-httprouter-users":  {"api/files/internal/adapters/handlers/std/users.tmpl", "internal/adapters/handlers/httprouter/users.go"},
	"graphql-schema":        {"api/files/internal/adapters/handlers/graphql/graphql.tmpl", "internal/adapters/handlers/graphql/graphql.go"},

	"dynamodb-env":         {"database/files/env/env-dynamodb.tmpl", ".env"},
	"dynamodb-repository":  {"database/files/internal/adapters/repositories/dynamodb/repository.tmpl", "internal/adapters/repositories/dynamodb/repository.go"},
	"dynamodb-readme":      {"database/files/readme/readme-dynamodb.tmpl", "README.md"},
	"mongodb-env":          {"database/files/env/env-mongodb.tmpl", ".env"},
	"mongodb-repository":   {"database/files/internal/adapters/repositories/mongodb/repository.tmpl", "internal/adapters/repositories/mongodb/repository.go"},
	"mongodb-readme":       {"database/files/readme/readme-mongodb.tmpl", "README.md"},
	"mysql-env":            {"database/files/env/env-mysql.tmpl", ".env"},
	"mysql-repository":     {"database/files/internal/adapters/repositories/mysql/repository.tmpl", "internal/adapters/repositories/mysql/repository.go"},
	"mysql-readme":         {"database/files/readme/readme-mysql.tmpl", "README.md"},
	"postgres-env":         {"database/files/env/env-postgres.tmpl", ".env"},
	"postgres-repository":  {"database/files/internal/adapters/repositories/postgres/repository.tmpl", "internal/adapters/repositories/postgres/repository.go"},
	"postgres-readme":      {"database/files/readme/readme-postgres.tmpl", "README.md"},
	"sqlite-env":           {"database/files/env/env-sqlite.tmpl", ".env"},
	"sqlite-repository":    {"database/files/internal/adapters/repositories/sqlite/repository.tmpl", "internal/adapters/repositories/sqlite/repository.go"},
	"sqlite-readme":        {"database/files/readme/readme-sqlite.tmpl", "README.md"},
	"badgerdb-env":         {"database/files/env/env-badgerdb.tmpl", ".env"},
	"badgerdb-repository":  {"database/files/internal/adapters/repositories/badgerdb/repository.tmpl", "internal/adapters/repositories/badgerdb/repository.go"},
	"badgerdb-readme":      {"database/files/readme/readme-badgerdb.tmpl", "README.md"},
	"mariadb-env":          {"database/files/env/env-mysql.tmpl", ".env"},
	"mariadb-repository":   {"database/files/internal/adapters/repositories/mysql/repository.tmpl", "internal/adapters/repositories/mariadb/repository.go"},
	"mariadb-readme":       {"database/files/readme/readme-mysql.tmpl", "README.md"},
	"cassandra-env":        {"database/files/env/env-cassandra.tmpl", ".env"},
	"cassandra-repository": {"database/files/internal/adapters/repositories/cassandra/repository.tmpl", "internal/adapters/repositories/cassandra/repository.go"},
	"cassandra-readme":     {"database/files/readme/readme-cassandra.tmpl", "README.md"},

	"web-app-readme":              {"web-app/files/readme/readme.tmpl", "README.md"},
	"web-app-env":                 {"web-app/files/env.tmpl", ".env"},
	"web-app-http-routes":         {"web-app/files/internal/adapters/handlers/http/routes.tmpl", "internal/adapters/handlers/http/routes.go"},
	"web-app-http-users":          {"web-app/files/internal/adapters/handlers/std/users.tmpl", "internal/adapters/handlers/http/users.go"},
	"web-app-chi-routes":          {"web-app/files/internal/adapters/handlers/chi/routes.tmpl", "internal/adapters/handlers/chi/routes.go"},
	"web-app-chi-users":           {"web-app/files/internal/adapters/handlers/std/users.tmpl", "internal/adapters/handlers/chi/users.go"},
	"web-app-echo-routes":         {"web-app/files/internal/adapters/handlers/echo/routes.tmpl", "internal/adapters/handlers/echo/routes.go"},
	"web-app-echo-users":          {"web-app/files/internal/adapters/handlers/echo/users.tmpl", "internal/adapters/handlers/echo/users.go"},
	"web-app-gin-routes":          {"web-app/files/internal/adapters/handlers/gin/routes.tmpl", "internal/adapters/handlers/gin/routes.go"},
	"web-app-gin-users":           {"web-app/files/internal/adapters/handlers/gin/users.tmpl", "internal/adapters/handlers/gin/users.go"},
	"web-app-gorilla-routes":      {"web-app/files/internal/adapters/handlers/gorilla/routes.tmpl", "internal/adapters/handlers/gorilla/routes.go"},
	"web-app-gorilla-users":       {"web-app/files/internal/adapters/handlers/std/users.tmpl", "internal/adapters/handlers/gorilla/users.go"},
	"web-app-httprouter-routes":   {"web-app/files/internal/adapters/handlers/httprouter/routes.tmpl", "internal/adapters/handlers/httprouter/routes.go"},
	"web-app-httprouter-users":    {"web-app/files/internal/adapters/handlers/std/users.tmpl", "internal/adapters/handlers/httprouter/users.go"},
	"web-app-users-html-template": {"web-app/files/assets/users.tmpl", "internal/web/templates/users.html"},
	"web-app-user-html-template":  {"web-app/files/assets/user.tmpl", "internal/web/templates/user.html"},
	"web-app-styles":              {"web-app/files/assets/styles.tmpl", "internal/web/static/styles.css"},
	"web-app-embed":               {"web-app/files/internal/web/web.tmpl", "internal/web/web.go"},

	"worker-main":                  {"microservice/files/cmd/main.tmpl", "cmd/main.go"},
	"redis-microservice-readme":    {"microservice/files/readme/readme-redis.tmpl", "README.md"},
	"redis-microservice-env":       {"microservice/files/env/env-redis.tmpl", ".env"},
	"redis-microservice-broker":    {"microservice/files/internal/adapters/brokers/redis/redis.tmpl", "internal/adapters/brokers/redis/redis.go"},
	"kafka-microservice-readme":    {"microservice/files/readme/readme-kafka.tmpl", "README.md"},
	"kafka-microservice-env":       {"microservice/files/env/env-kafka.tmpl", ".env"},
	"kafka-microservice-broker":    {"microservice/files/internal/adapters/brokers/kafka/kafka.tmpl", "internal/adapters/brokers/kafka/kafka.go"},
	"rabbitmq-microservice-readme": {"microservice/files/readme/readme-rabbitmq.tmpl", "README.md"},
	"rabbitmq-microservice-env":    {"microservice/files/env/env-rabbitmq.tmpl", ".env"},
	"rabbitmq-microservice-broker": {"microservice/files/internal/adapters/brokers/rabbitmq/rabbitmq.tmpl", "internal/adapters/brokers/rabbitmq/rabbitmq.go"},
}

// Lookup returns a template's content and output path. An unknown key is a
// scaffolder bug; it is reported as an error so the caller can roll back the
// half-generated project instead of panicking.
func Lookup(key string) ([]byte, string, error) {
	t, ok := registry[key]
	if !ok {
		return nil, "", fmt.Errorf("unknown template %q", key)
	}
	content, err := files.ReadFile(t.src)
	if err != nil {
		return nil, "", fmt.Errorf("template %q: %w", key, err)
	}
	return content, t.out, nil
}

// Mapper returns every template as a key -> loader map, for callers that walk
// the whole registry (the template smoke test).
func Mapper() map[string]func() ([]byte, string) {
	m := make(map[string]func() ([]byte, string), len(registry))
	for key := range registry {
		m[key] = func() ([]byte, string) {
			content, out, err := Lookup(key)
			if err != nil {
				panic(err)
			}
			return content, out
		}
	}
	return m
}
