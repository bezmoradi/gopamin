# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Gopamin is a CLI tool for scaffolding Go projects based on Clean Architecture principles and the Standard Go Project Layout. It generates boilerplate Go projects with various configurations for web apps, APIs, broker-driven workers, and hello-world applications.

## Build and Run Commands

```bash
# Build the CLI
go build -o gopamin github.com/bezmoradi/gopamin

# Install globally
go install github.com/bezmoradi/gopamin@latest

# Run commands
gopamin version
gopamin new -h
gopamin new -n my-app -t api -f echo -l log -d mysql
# an API that also produces/consumes events and owns a database:
gopamin new -n my-app -t api -f echo -b kafka -d postgres -l slog
# a broker-driven worker that persists what it consumes:
gopamin new -n my-worker -t worker -b kafka -d postgres -l slog
```

## Architecture

### Package Structure

- `main.go` - Entry point, invokes root command
- `internal/commands/` - Cobra CLI commands and validation
- `internal/project/` - The `Project` type every template renders with, plus the derived predicates (`PublishesEvents`, `HasCompose`, `IsRelational`, `HasHTTP`, ...) that both the scaffolder and the templates branch on
- `internal/scaffolder/` - Project scaffolding logic (data-driven recipes + one build function per project type)
- `internal/templates/` - Embedded template files (using Go's `embed`)

### Key Design Patterns

**Data-driven generation**: The scaffolder keeps the *variation* between project types / platforms / databases / loggers as plain recipe data and shares the build logic. The flow is:
1. CLI parses flags → `Project` struct created
2. `generate` creates a `generator` (sticky first error, collected modules) and `buildersMap[ProjectType]` dispatches to `buildAPI` / `buildWebApp` / `buildWorker` / `buildHelloWorld` (`builder.go`)
3. That build function reads its recipe (e.g. `apiRecipes[Platform]`), emits its own files and its `.env`/README parts; `buildDatabase(g, db)` emits the repository stack; `buildCommon(g)` then emits everything the types share (logger, broker, otel/auth/openapi, `.env`, README, Makefile, compose, configs, tools, lint configs, and a Dockerfile for every type but hello-world)
4. `fileGenerator` renders each template, runs `go/format.Source` on `.go` output (so generated code is always `gofmt`-clean), then writes the file; every failure, including an unknown template key, is returned as an error
5. External commands run (git init; `go mod init` pinned to `generatedGoVersion`; `go get module@version` from `modulePins` in `versions.go`, never `-u`, plus the `securityFloor` transitive pins; `go mod tidy`). On failure the partial project directory is rolled back.

(There is no `director`/builder-interface ceremony — those were removed in favor of the recipe tables.)

**Template Registry**: `internal/templates/main.go` embeds every template in one `embed.FS` and maps each key to `{src, out}` (`registry`); `Lookup(key)` returns the content and output path, or an error for an unknown key. The `api/`, `common/`, `database/`, `hello-world/`, `microservice/` and `web-app/` directories are template data, not Go packages. Templates that render to the same path (`.env`, `README.md`, `Makefile`) are concatenated in the order the builder lists them.

### Supported Configurations

- **Project Types** (`-t`): hello-world, web-app, api, worker
- **HTTP Frameworks** (`-f`): echo, chi, gin, http, gorilla, httprouter, graphql (api only; web-app has the first six) — required for api/web-app
- **Message Brokers** (`-b`): kafka, rabbitmq, redis — required for worker, optional for api/web-app
- **Databases** (`-d`): mysql, mariadb, postgres, sqlite, mongodb, cassandra, dynamodb, badgerdb — optional for all types
- **Loggers** (`-l`): log, slog, logrus, zap — required
- **CI** (`-c`): github, gitlab — optional (omit for none); emits a build/test/`make lint`/`make arch-check` pipeline (`.github/workflows/ci.yml` or `.gitlab-ci.yml`)
- **Auth** (`-a`): jwt — optional, api only; bearer-token middleware on every app route plus a demo `POST /login` (`internal/adapters/handlers/auth`, `cmd/hashpassword`, a `JWT_*`/`AUTH_*` block in `.env`)
- **Observability** (`-o`): otel — optional for api/web-app/worker; OpenTelemetry traces and runtime metrics over OTLP (`internal/adapters/observability`, an `OTEL_*` block, a collector service in compose)
- **OpenAPI** (`-s`): swagger — optional, REST api only; a templated `openapi.yaml` plus Swagger UI at `/swagger`

A `worker` is a broker-driven service with no HTTP server. An `api`/`web-app` given a `-b` broker emits a producer (publishes `user.created`/`user.updated`/`user.deleted` domain events on writes, via an `EventPublisher` port + the service's `WithPublisher`) **and** a consumer loop running alongside the HTTP server (one binary, shared core, graceful shutdown; a consumer failure stops the server and exits 1). `redis` is **not** a database — it's a cache/broker, offered only via `-b`.

Events: the core owns `domain.Event`; `internal/adapters/events` (arch component `events`) owns the wire format, a versioned JSON envelope with a UUID id, and `events.NewPublisher(broker)` implements the port. Broker adapters only move bytes (`Produce`, `Consume(ctx, func(ctx, []byte) error)`), commit/ack after the handler succeeds and retry a failing handler with backoff (Kafka: consumer group `KAFKA_GROUP_ID`; RabbitMQ: durable queue, confirms, manual ack; Redis: at-most-once pub/sub). Topic/queue/channel and group default to `Project.Slug()`. A worker with a database persists `user.created` via `CreateOnce(eventID, name)` (worker-only port method `InsertOnce`; SQL/Mongo use a unique `event_id` column/index, Cassandra/DynamoDB/Badger key by the event id), so redeliveries don't duplicate rows.

**Combination matrix** — the validators (`internal/commands/validators.go`, which read the allowed values from the recipe tables via `scaffolder.ProjectTypes()`/`Frameworks()`/`Brokers()`/`Databases()`/`Loggers()`/`CIs()`) accept **2,016** distinct flag combinations of `gopamin new -t <type> -f <framework> -b <broker> -d <database> -l <logger> -n <name>`. Rules: `-l` is required for every type; `-f` is required for api/web-app and rejected for hello-world/worker; `-b` is required for worker and optional for api/web-app; `-d` is optional for all. Selection is orthogonal — any framework pairs with any broker (or none), any database (or none), and any logger.

| Type (`-t`) | Framework (`-f`) | Broker (`-b`) | Database (`-d`) | × Loggers | = Combos |
|------|------|------|------|------|------|
| hello-world | none | none | none + any of 8 = 9 | 4 | 36 |
| api | echo, chi, gin, http, gorilla, httprouter, graphql (7) | none + any of 3 = 4 | none + any of 8 = 9 | 4 | 1,008 |
| web-app | the same minus graphql (6) | none + any of 3 = 4 | none + any of 8 = 9 | 4 | 864 |
| worker | none | any of 3 (required) = 3 | none + any of 8 = 9 | 4 | 108 |
| **Total** | | | | | **2,016** |

The 8 databases (`-d`): mysql, mariadb, postgres, sqlite, mongodb, cassandra, dynamodb, badgerdb. The 3 brokers (`-b`): kafka, rabbitmq, redis. The 4 loggers (`-l`): log, slog, logrus, zap. The 2,016 counts what the validators *accept* — a representative cross-section is build-tested. That's 504 structural (type × framework × broker × database) layouts × 4 loggers. Recompute when adding a framework/broker/database/logger.

The optional `-c` (CI: none/github/gitlab) is an **orthogonal, combination-invariant** dimension: it multiplies the accepted total by 3 (→ **6,048**) but the generated pipeline is identical regardless of type/framework/broker/database/logger, so only one CI-on project per platform needs build-testing. The table above is left at 2,016 (the structural matrix) since `-c` adds nothing structural. `-a`, `-o` and `-s` are left out of the count too; test them on a few shapes each.

### Recipe & dispatch tables (in `internal/scaffolder/`)

- `buildersMap` - project type → build function (`buildAPI`, `buildWebApp`, `buildWorker`, `buildHelloWorld`)
- `apiRecipes` / `webRecipes` - framework → template-file keys + Go modules to `go get`
- `databaseRecipes` - database → repository file, modules (consumed by `buildDatabase`)
- `loggerRecipes` - logger → logger file + modules (consumed by `buildLogger`)
- `workerRecipes` - broker → modules (also used to fetch the broker module when api/web-app add `-b`)
- `Project.HasCompose()` (in `internal/project`) - whether to emit the merged `docker-compose.yml` + `docker-*` make targets (a containerized database, a broker, or `-o otel`; sqlite/badgerdb are embedded). AGENTS.md uses the same method, so the docs can't drift from the files

To add a framework/database/broker/logger, add one recipe entry (plus its templates) rather than a new builder struct. Docker artifacts are assembled from one conditional `common/files/docker-compose.tmpl` (services chosen by `.Database`/`.Broker`), so any database+broker combination merges into a single compose file. Every compose service has a healthcheck that passes only once clients can connect, and `make docker-up` runs `docker compose up -d --wait`: generated apps fail fast at startup (a 5 s ping, one broker dial), so "healthy" must mean "ready". The RabbitMQ check is a TCP connect, never the Erlang CLI (run as root before the broker writes its cookie, the CLI writes one the broker can't read, and the broker exits). (The CLI flag is `-f/--framework`; internally the framework value is still carried on `Project.Platform`, and the broker templates still live under `microservice/` with keys `*-microservice-*` — both pending an internal cosmetic rename.) One `http-main` and one `http-server` template serve every api/web-app (REST, GraphQL, web; with or without a database), one `worker-main` every worker, and the mysql templates serve mariadb too. Framework adapters return a plain `http.Handler` (`NewRouter` for REST/web, `NewHandler` for GraphQL); `server.Serve` mounts it behind the shared middleware chain, so no HTTP type lives in `internal/core`.

### Generated Project Structure (Clean Architecture)

Generated projects follow this layout:
- `cmd/` - Entry points (`main.go` with a `run() error`, so deferred `Close` calls run; plus `cmd/server/` for api/web, whose `Serve(ctx, handler, logger, checks...)` returns errors; plus `cmd/hashpassword` with `-a jwt`)
- `internal/core/user/{domain,ports,services}` - the user aggregate: entities and events (`domain.User{ID,Name}`, `domain.Event`), ports (`Update`/`Delete` return `domain.ErrUserNotFound`; constructors return errors), and `services.UserService`
- `internal/adapters/` - repositories (each `New()`/`Ping`/`Close`; `memory` when there is no `-d`), handlers (incl. `health` with readiness checks and `middleware`), loggers, brokers (bytes only), `events` (the envelope), `observability`, `api`
- `internal/web/` - web-app only: HTML templates (incl. `error.html`) and CSS embedded via `//go:embed` (binary is self-contained)
- `migrations/` - relational only; SQLite also embeds it (`migrations/migrations.go`) and applies it at startup
- `configs/`, `tools/`, and an `AGENTS.md` (agent/contributor guidance) at the project root

### Architecture Boundary Enforcement

**Every** generated project ships a `.go-arch-lint.yml` (template: `common/files/arch-lint.tmpl`) plus `make arch` / `make arch-check` targets — including a plain hello-world, whose config is a thin variant (no `domain`/`services`; just `cmd`/`configs`/`tools`/`loggers`/`ports`). This makes the Clean Architecture dependency rule (which AGENTS.md/README describe in prose) a *build-enforced* invariant via [go-arch-lint](https://github.com/fe3dback/go-arch-lint), pinned and invoked as `go run …@<version>` (never added to the generated `go.mod`; version lives in `makefile.tmpl` as `GO_ARCH_LINT_VERSION`).

Load-bearing decisions (all empirically verified against go-arch-lint **v1.15.0**, schema version 3 — do not change without re-verifying):
- **The component set is templated per project, not a static superset.** A component whose glob matches zero directories makes go-arch-lint exit 1, so `arch-lint.tmpl` declares only the components a given project actually has (guards on `.ProjectType`/`.Platform`/`.Broker`/`.Database`; e.g. a plain hello-world omits `domain`/`services` entirely, and `ports` becomes a leaf there). Adapter *categories* still use combination-invariant globs (`handlers/*`, `repositories/*`, `loggers/*`, `brokers/*`), so the specific framework/db/logger/broker needs no per-combo work — only whole-component presence varies.
- **`excludeFiles: ["^.*_test\\.go$"]`** — tests may wire real adapters to the code under test; boundaries are enforced on the production import graph. (The generated service test no longer needs it: it uses a `fakeRepository` defined in the test, and api/web-app without `-d` use `repositories/memory` at runtime — the old `mock` package had both roles.)
- **`deepScan: false`** — deepScan traces DI data-flow and misattributes `cmd`'s wiring as a phantom adapter→service edge.
- **`depOnAnyVendor: false`; vendors granted per-adapter via `anyVendorDeps`.** The core (domain/ports/services) is genuinely pure (zero third-party imports), so it gets **no** vendor grant — a driver/framework leaking into the core is a hard failure. `arch-check` greps the linter banner and refuses to pass if the base or vendor-import linter has been silently turned off (exits 0 while `Off`).

To add a new generated adapter *category* (a new top-level dir), add a matching component to `arch-lint.tmpl` — otherwise deny-by-default (`workdir: .`) makes `make arch` fail on the unattached package.

## Dependencies

The CLI itself depends on:
- `github.com/spf13/cobra` - CLI framework
- `github.com/fatih/color` - Terminal colors
- `gopkg.in/yaml.v3` - test-only (validates the generated `.go-arch-lint.yml` in `templates_test.go`)

`github.com/joho/godotenv` is **not** a CLI dependency — it is `go get`-installed into *generated* projects (for `.env` loading), passed as a string in `scaffolder.go`. (`promptui`/`readline` were dropped in the `ec0c90c` refactor and removed from `go.mod` by `go mod tidy` here.)

## Testing

The CLI has one test file on purpose: verify changes by generating real projects, not by adding CLI unit tests. Tests inside generated projects (the templates' `*_test.tmpl`) can change freely.

`go test ./...` runs `internal/templates/templates_test.go` — a hermetic smoke test that (a) renders every template with **one** representative api shape and `go/parser`-parses the Go output (branches for other shapes — worker-only blocks, other databases/brokers — are checked only by generating projects) and (b) renders `.go-arch-lint.yml` across every project shape, YAML-parsing each and asserting no dependency rule references an undeclared component (no project generation, no network). This catches template-directive/Go-syntax errors and arch-config regressions fast. After changing templates or the scaffolder, also generate a few real projects and run `go build` / `go vet` / `golangci-lint` / `go test` **and `make arch`** on them (the YAML parse alone can't catch an empty-glob component or an unattached package — only real `go-arch-lint check` on the generated tree can).

For a change that reaches runtime (compose, adapters, `cmd/`), pick shapes that between them use every affected flag value and run them: `make docker-up`, `make migrate-up` for postgres/mysql/mariadb, then the app or its image (AGENTS.md "Container image" lists the host overrides), and drive it: `/healthz`, `/readyz`, a create/read/update/delete round trip, events consumed, a clean stop. Generate into a throwaway directory, give Docker resources unique names, and remove them afterwards (`make docker-down-rm-vol`).

## Conventions

- **Echo is on v5 on purpose.** Don't move it back to v4; raise the trade-offs before any other major-version bump.
- **Transport types live in the adapters.** API request/response and error types belong in `internal/adapters/api` (or a handler package), never in `internal/core`; domain entities carry no `json` tags.
- **Compose images are pinned to series tags** (`postgres:17`, `rabbitmq:4.4-management`, ...), never `:latest`. Kafka runs `apache/kafka` in KRaft mode, with no ZooKeeper.
- **Correctness beats compatibility.** Fix a mismodeled concept properly even when that changes the CLI or the generated code; generated projects are one-shot, so there is no upgrade path to protect.

## Version Management

Version is defined in `internal/commands/constants.go` (keep the `VERSION = "v…"` form: releases up to v1.2.0 read it from `main` for their update check). `gopamin new` runs an **advisory, timed** update check against the Go module proxy's `@latest` (semver comparison): it prints a hint only when a newer release exists, **never blocks**, silently skips on network/timeout failure, and `GOPAMIN_NO_UPDATE_CHECK=1` turns it off. `version`/`help`/completion never run it.

Generated projects' `go` directive is pinned by the `generatedGoVersion` constant in `internal/scaffolder/scaffolder.go` (keep it in lockstep with this repo's `go.mod`). Their runtime modules are pinned in `internal/scaffolder/versions.go` (`modulePins`); bump them deliberately and re-run a generated cross-section (build, vet, test, lint, arch-check, govulncheck) before releasing.
