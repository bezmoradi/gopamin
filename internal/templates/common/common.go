package templates

import (
	_ "embed"
)

//go:embed files/env.tmpl
var env []byte

func EnvTemplate() ([]byte, string) {
	return env, ".env"
}

//go:embed files/makefile.tmpl
var makefile []byte

func MakefileTemplate() ([]byte, string) {
	return makefile, "Makefile"
}

//go:embed files/docker-makefile.tmpl
var dockerMakefile []byte

// DockerMakefileTemplate renders the generic `docker-*` make targets, appended to
// the Makefile for any project that ships a containerized backend (a database
// and/or a broker). The targets are backend-agnostic because the merged
// docker-compose.yml already contains every service.
func DockerMakefileTemplate() ([]byte, string) {
	return dockerMakefile, "Makefile"
}

//go:embed files/docker-compose.tmpl
var dockerCompose []byte

// DockerComposeTemplate renders a single docker-compose.yml that contains every
// containerized backend the project uses (database and/or broker), assembled by
// conditionals on .Database and .Broker.
func DockerComposeTemplate() ([]byte, string) {
	return dockerCompose, "docker-compose.yml"
}

//go:embed files/dockerfile.tmpl
var dockerfile []byte

// DockerfileTemplate renders a multi-stage Dockerfile that compiles the service and
// ships only the binary on a minimal, non-root distroless base. sqlite selects the
// cgo driver, so it builds with cgo on and runs on a libc base (distroless/base);
// every other database builds a static binary on the smaller distroless/static. The
// build image tag comes from .GoVersion so it never drifts from the go.mod directive.
func DockerfileTemplate() ([]byte, string) {
	return dockerfile, "Dockerfile"
}

//go:embed files/dockerignore.tmpl
var dockerignore []byte

// DockerignoreTemplate keeps the build context lean and secret-free (excludes .env,
// .git, docs, and runtime artifacts).
func DockerignoreTemplate() ([]byte, string) {
	return dockerignore, ".dockerignore"
}

//go:embed files/internal/core/ports/event-publisher.interface.tmpl
var eventPublisherInterface []byte

func EventPublisherInterfaceTemplate() ([]byte, string) {
	return eventPublisherInterface, "internal/core/ports/event-publisher.interface.go"
}

//go:embed files/internal/core/services/user.service-broker.tmpl
var userServiceBroker []byte

// UserServiceBrokerTemplate renders to the same path as UserServiceTemplate, so
// when a broker is selected the builder generates it after buildDatabase to
// replace the plain user service with the publisher-aware variant.
func UserServiceBrokerTemplate() ([]byte, string) {
	return userServiceBroker, "internal/core/services/user.service.go"
}

//go:embed files/readme/readme.tmpl
var readme []byte

func ReadmeTemplate() ([]byte, string) {
	return readme, "README.md"
}

//go:embed files/readme/readme-log.tmpl
var readmeLog []byte

func ReadmeLogTemplate() ([]byte, string) {
	return readmeLog, "README.md"
}

//go:embed files/readme/readme-slog.tmpl
var readmeSlog []byte

func ReadmeSlogTemplate() ([]byte, string) {
	return readmeSlog, "README.md"
}

//go:embed files/readme/readme-logrus.tmpl
var readmeLogrus []byte

func ReadmeLogrusTemplate() ([]byte, string) {
	return readmeLogrus, "README.md"
}

//go:embed files/readme/readme-zap.tmpl
var readmeZap []byte

func ReadmeZapTemplate() ([]byte, string) {
	return readmeZap, "README.md"
}

//go:embed files/license.tmpl
var license []byte

func LicenseTemplate() ([]byte, string) {
	return license, "LICENSE"
}

//go:embed files/gitignore.tmpl
var gitIgnore []byte

func GitIgnoreTemplate() ([]byte, string) {
	return gitIgnore, ".gitignore"
}

//go:embed files/AGENTS.tmpl
var agents []byte

func AgentsTemplate() ([]byte, string) {
	return agents, "AGENTS.md"
}

//go:embed files/arch-lint.tmpl
var archLint []byte

// ArchLintTemplate renders the go-arch-lint boundary-enforcement config. It is
// emitted for every project type that has adapters to enforce (api, web-app,
// worker, and hello-world with a database); the component/dep set is shaped by the
// project's type/platform/broker/database so no empty component is ever declared.
func ArchLintTemplate() ([]byte, string) {
	return archLint, ".go-arch-lint.yml"
}

//go:embed files/internal/adapters/repositories/mock/repository.tmpl
var mockRepository []byte

func MockRepositoryTemplate() ([]byte, string) {
	return mockRepository, "internal/adapters/repositories/mock/repository.go"
}

//go:embed files/internal/core/domain/user_test.tmpl
var userTest []byte

func UserTestTemplate() ([]byte, string) {
	return userTest, "internal/core/domain/user_test.go"
}

//go:embed files/internal/core/domain/user.tmpl
var user []byte

func UserTemplate() ([]byte, string) {
	return user, "internal/core/domain/user.go"
}

//go:embed files/internal/core/ports/user-repository.interface.tmpl
var userRepositoryInterface []byte

func UserRepositoryInterfaceTemplate() ([]byte, string) {
	return userRepositoryInterface, "internal/core/ports/user-repository.interface.go"
}

//go:embed files/internal/core/ports/user-service.interface.tmpl
var userServiceInterface []byte

func UserServiceInterfaceTemplate() ([]byte, string) {
	return userServiceInterface, "internal/core/ports/user-service.interface.go"
}

//go:embed files/internal/core/services/user.service.tmpl
var userService []byte

func UserServiceTemplate() ([]byte, string) {
	return userService, "internal/core/services/user.service.go"
}

//go:embed files/internal/core/services/user.service_test.tmpl
var userServiceTest []byte

func UserServiceTestTemplate() ([]byte, string) {
	return userServiceTest, "internal/core/services/user.service_test.go"
}

//go:embed files/internal/core/ports/router.interface.tmpl
var routerInferface []byte

func RouterInterfaceTemplate() ([]byte, string) {
	return routerInferface, "internal/core/ports/router.interface.go"
}

//go:embed files/internal/core/ports/logger.interface.tmpl
var loggerInferface []byte

func LoggerInterfaceTemplate() ([]byte, string) {
	return loggerInferface, "internal/core/ports/logger.interface.go"
}

//go:embed files/internal/adapters/api/errors.tmpl
var apiErrors []byte

func ApiErrorsTemplate() ([]byte, string) {
	return apiErrors, "internal/adapters/api/errors.go"
}

//go:embed files/internal/adapters/api/response.tmpl
var apiResponse []byte

func ApiResponseTemplate() ([]byte, string) {
	return apiResponse, "internal/adapters/api/response.go"
}

//go:embed files/configs/configs.tmpl
var configs []byte

func ConfigsTemplate() ([]byte, string) {
	return configs, "configs/configs.go"
}

//go:embed files/configs/configs_test.tmpl
var configsTest []byte

func ConfigsTestTemplate() ([]byte, string) {
	return configsTest, "configs/configs_test.go"
}

//go:embed files/internal/core/domain/message.tmpl
var message []byte

func MessageTemplate() ([]byte, string) {
	return message, "internal/core/domain/message.go"
}

//go:embed files/internal/core/ports/broker-service.interface.tmpl
var brokerServiceInterface []byte

func BrokerServiceInterfaceTemplate() ([]byte, string) {
	return brokerServiceInterface, "internal/core/ports/broker-service.interface.go"
}

//go:embed files/internal/core/ports/message-broker.interface.tmpl
var messageBrokerInterface []byte

func MessageBrokerInterfaceTemplate() ([]byte, string) {
	return messageBrokerInterface, "internal/core/ports/message-broker.interface.go"
}

//go:embed files/internal/core/services/broker.service.tmpl
var brokerService []byte

func BrokerServiceTemplate() ([]byte, string) {
	return brokerService, "internal/core/services/broker.service.go"
}

//go:embed files/tools/tools.tmpl
var tools []byte

func ToolsTemplate() ([]byte, string) {
	return tools, "tools/tools.go"
}

//go:embed files/tools/tools_test.tmpl
var toolsTest []byte

func ToolsTestTemplate() ([]byte, string) {
	return toolsTest, "tools/tools_test.go"
}

//go:embed files/internal/adapters/loggers/log/logger.tmpl
var logLogger []byte

func LogLoggerTemplate() ([]byte, string) {
	return logLogger, "internal/adapters/loggers/log/logger.go"
}

//go:embed files/internal/adapters/loggers/logrus/logger.tmpl
var logrusLogger []byte

func LogrusLoggerTemplate() ([]byte, string) {
	return logrusLogger, "internal/adapters/loggers/logrus/logger.go"
}

//go:embed files/internal/adapters/loggers/slog/logger.tmpl
var slogLogger []byte

func SlogLoggerTemplate() ([]byte, string) {
	return slogLogger, "internal/adapters/loggers/slog/logger.go"
}

//go:embed files/internal/adapters/loggers/zap/logger.tmpl
var zapLogger []byte

func ZapLoggerTemplate() ([]byte, string) {
	return zapLogger, "internal/adapters/loggers/zap/logger.go"
}
