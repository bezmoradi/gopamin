<h1 align="center">⚙️ Engineered, Not Vibe Coded</h1>

<p align="center">
  <img src="https://img.shields.io/badge/Engineered-not%20vibe%20coded-blue?style=for-the-badge" alt="Engineered, not vibe coded">
</p>

## #1 STOP VIBE CODING <img src="https://img.shields.io/badge/Merged-8957e5?style=flat-square&logo=git&logoColor=white" alt="Merged">

```diff
- vibe coded: shipped because it compiled
+ engineered: designed, reviewed, documented
```

<sub>✅ Approved by human reviewers · All checks passed</sub>

## Table of Contents

- [Introduction](#introduction)
- [An Intro to The Clean Architecture](#an-intro-to-the-clean-architecture)
- [Prerequisites](#prerequisites)
- [Supported Flags](#supported-flags)
- [Recipes](#recipes)
- [Installation](#installation)
  - [Install somewhere else](#install-somewhere-else)
  - [Update](#update)
  - [Check the script first](#check-the-script-first)
  - [Other ways to install](#other-ways-to-install)
- [Guides](#guides)
- [Development](#development)
  - [Releasing](#releasing)
- [Author](#author)
- [License](#license)

## Introduction

Gopamin is a CLI which creates new projects based on ideas promoted by [Standard Go Project Layout](https://github.com/golang-standards/project-layout) which is a **well-accepted** architecture by the Go community (It's not an official standard defined by the core Go team; however, it is a set of common historical and emerging project layout patterns in the Go ecosystem).

## An Intro to The Clean Architecture

All boilerplates created by Gopamin are based on [The Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html). Simply put, it is a software architectural pattern introduced by Robert C. Martin in order to create maintainable, scalable, and testable software systems by decoupling the application logic from the infrastructure details (Although there are some minor differences, this architectural pattern is referred to by other names such as Ports & Adapters Architecture, Hexagonal Architecture, and Onion Architecture).

## Prerequisites

The minimum required tool for using the Gopamin CLI tool is Go 1.27.3 or higher, which can be downloaded from [Go All Releases](https://go.dev/dl). Go 1.21 or later also works: it downloads Go 1.27.3 by itself the first time Gopamin creates a project, unless `GOTOOLCHAIN` is set to `local` (or another setting that stops Go fetching toolchains). To have full development setup though, other tools are also recommended to be installed on your local machine:

- **[Git](https://git-scm.com/)**: By default each new project created by this tool initializes a Git repo; that's why you need to make sure Git is installed on your machine.
- **[Docker](https://www.docker.com)**: If you choose to create a new project with a containerized database, a message broker, and/or observability, a merged `docker-compose.yml` will be included in the root of the project for running those backends locally.
- **[GNU Make](https://www.gnu.org/software/make)**: This a tool which controls the generation of executables and other types of files. By default, each new project includes a `Makefile` for running some most-used commands like running an application (This tool is installed by default on Mac and some distributions of Linux). To check whether this tool is installed on your machine, open terminal and run `make --version` (If you do not have this tool installed on your machine though, still you can use this tools without any limitations).

## Supported Flags

A project is defined by **orthogonal capabilities**: the **type** (`-t`) picks the app's primary interface (HTML, JSON/GraphQL, or message-driven), while a **framework** (`-f`), a **broker** (`-b`), and a **database** (`-d`) are optional capabilities you attach to it. On top of those, a set of **cross-cutting features** — a CI pipeline (`-c`), observability (`-o`), authentication (`-a`), and API documentation (`-s`) — can be layered on. You can mix them freely.

| Flag | Description | Required | Options |
|------|-------------|----------|---------|
| `-n`, `--name` | Project name, or a full module path | Always | Any string starting with a letter or digit; characters outside `A-Z a-z 0-9 . - _` are replaced with `-`. A module path such as `github.com/You/repo` sets the `go.mod` module as written; the directory name is lowercased. |
| `-t`, `--type` | The app's primary interface | Always | `hello-world`, `web-app`, `api`, `worker` |
| `-f`, `--framework` | HTTP framework | `api` & `web-app` only | `echo`, `chi`, `gin`, `httprouter`, `gorilla`, `http` — plus `graphql` for `api` |
| `-b`, `--broker` | Message broker — publishes a versioned event on every write and runs a consumer loop | Required for `worker`; optional for `api` / `web-app` | `kafka`, `rabbitmq`, `redis` |
| `-d`, `--database` | Database | Optional (for any type) | `mysql`, `mariadb`, `postgres`, `mongodb`, `cassandra`, `sqlite`, `dynamodb`, `badgerdb` |
| `-l`, `--logger` | Logging library | Always | `log`, `slog`, `logrus`, `zap` |
| `-c`, `--ci` | CI pipeline (build, test, `make lint`, `make arch-check`) | Optional (for any type) | `github` (`.github/workflows/ci.yml`), `gitlab` (`.gitlab-ci.yml`) |
| `-o`, `--observability` | Observability stack | Optional (any type **except** `hello-world`) | `otel` — OpenTelemetry traces + Go runtime metrics, OTLP export to a collector |
| `-a`, `--auth` | Authentication | Optional (`api` only) | `jwt` — bearer-token middleware on every route + a demo `POST /login` issuer (golang-jwt/v5, HS256 + algorithm pinning) |
| `-s`, `--openapi` | API documentation | Optional (REST `api` only — not `graphql`) | `swagger` — a templated `openapi.yaml` + embedded Swagger UI at `/swagger` (via swaggest/swgui) |

A few things worth knowing:

- A `worker` is a **headless, broker-driven service** — it has no HTTP server (so it takes no `-f`) and is triggered by messages off its broker.
- Adding `-b` to an `api` or `web-app` gives you one binary that runs the HTTP server **and** a broker consumer side by side, publishing an event on every write. In other words, "an API that is also a worker" is `-t api -f <framework> -b <broker>`.
- When a containerized database and/or a broker is selected, a single merged `docker-compose.yml` (plus `docker-up` / `docker-down` make targets) is generated so every backend runs together for local development. Adding `-o otel` also merges an OpenTelemetry collector into that same compose file. Its ports are published on `127.0.0.1` only, and the passwords in `.env` are generated per project.
- Kafka and RabbitMQ deliver events **at least once**: a message is committed or acknowledged only after it is processed, and a worker with a database stores each event once. Redis pub/sub is **at most once**.
- Every `api`, `web-app` and `worker` ships a `Dockerfile` (non-root distroless image); the generated `AGENTS.md` explains running it next to the compose services.
- Redis is a broker/cache, offered only via `-b` — it is intentionally **not** a database option.
- The cross-cutting flags stack: a single `api` can carry a framework, a broker, a database, a logger, CI, observability, auth, **and** OpenAPI docs all at once.
- Every generated project ships a `.go-arch-lint.yml` plus `make arch` / `make arch-check` targets, so the Clean Architecture dependency rule is a **build-enforced** invariant, not just documentation.

## Recipes

Because the flags are orthogonal, you can combine them freely. Some representative projects:

| Command | What it scaffolds |
|---------|-------------------|
| `gopamin new -n app -t hello-world -l log` | Minimal starter app |
| `gopamin new -n app -t api -f echo -l slog` | REST API with an in-memory store |
| `gopamin new -n app -t api -f echo -d postgres -l slog` | REST API + PostgreSQL |
| `gopamin new -n app -t api -f echo -b kafka -l slog` | REST API **+ Kafka** producer/consumer in one binary |
| `gopamin new -n app -t api -f echo -b kafka -d postgres -l slog` | REST API that persists writes **and** publishes/consumes events |
| `gopamin new -n app -t api -f graphql -d mysql -l zap` | GraphQL API + MySQL |
| `gopamin new -n app -t web-app -f gin -d mongodb -l zap` | Server-rendered web app + MongoDB |
| `gopamin new -n app -t worker -b rabbitmq -d postgres -l slog` | Headless worker: consume `user.created` events → store each once |
| `gopamin new -n app -t worker -b kafka -l logrus` | Headless worker with no database |
| `gopamin new -n app -t api -f echo -d postgres -l slog -a jwt -s swagger` | REST API + PostgreSQL, JWT-protected, with Swagger UI at `/swagger` |
| `gopamin new -n app -t api -f chi -d mysql -l slog -o otel -c github` | REST API + MySQL, OpenTelemetry traces/metrics, and a GitHub Actions pipeline |
| `gopamin new -n app -t api -f echo -b kafka -d postgres -l slog -o otel -a jwt -s swagger -c gitlab` | The works: persistence, events, auth, docs, tracing, and CI in one binary |

## Installation

On Linux or macOS, amd64 or arm64, run:

```text
$ curl -fsSL https://github.com/bezmoradi/gopamin/releases/latest/download/install.sh | sh
```

The script:

- runs on linux and darwin, amd64 and arm64, and stops with one line on anything else, or when `curl`, `tar`, `git`, `sha256sum` or `shasum`, or Go is missing;
- needs Go because `gopamin new` runs it to create projects: Go 1.21 or later, which downloads Go 1.27.3 by itself, or Go 1.27.3 or later when `GOTOOLCHAIN` keeps Go from downloading toolchains (such as `GOTOOLCHAIN=local`);
- downloads the release's archive for your platform and checks it against the SHA-256 written into the script and against the archive's own `SHA256SUMS`, and checks that its `gopamin version` prints the release's version;
- installs `gopamin` into `~/.local/bin` without sudo, and prints the `export PATH=...` line to add to `~/.bashrc` or `~/.zshrc` if that folder isn't on your `PATH`. A failed download or check changes nothing.

Then make sure it works:

```text
$ gopamin version
```

### Install somewhere else

Set `GOPAMIN_INSTALL_DIR` to the folder to install into:

```text
$ curl -fsSL https://github.com/bezmoradi/gopamin/releases/latest/download/install.sh | GOPAMIN_INSTALL_DIR="$HOME/bin" sh
```

### Update

Run the same command again. `gopamin new` asks GitHub for the latest release and, when it is newer than yours, prints the command that updates the `gopamin` in the folder it runs from before generating the project:

```text
Gopamin v1.4.0 is available (you have v1.3.0). To update:
  curl -fsSL https://github.com/bezmoradi/gopamin/releases/latest/download/install.sh | GOPAMIN_INSTALL_DIR='/home/you/.local/bin' sh
Or, with Go: go install github.com/bezmoradi/gopamin@v1.4.0
Set GOPAMIN_NO_UPDATE_CHECK=1 to stop this check.
```

The check asks `https://api.github.com/repos/bezmoradi/gopamin/releases/latest`, sending only `gopamin/<version>` as the User-Agent. It never blocks: it gives up silently after 3 seconds or on any error. Other commands (`version`, `help`, completion) don't run it. Set `GOPAMIN_NO_UPDATE_CHECK=1` to turn it off; it is also off when `CI` or `GITHUB_ACTIONS` is set.

### Check the script first

To check the script before it runs, download it and verify that the release workflow built it, with a signed-in [GitHub CLI](https://cli.github.com):

```text
$ curl -fsSLO https://github.com/bezmoradi/gopamin/releases/latest/download/install.sh
$ gh attestation verify install.sh -R bezmoradi/gopamin --signer-workflow bezmoradi/gopamin/.github/workflows/release.yml
$ sh install.sh
```

Each release's notes show how to check an archive the same way.

### Other ways to install

With Go (see [Prerequisites](#prerequisites)), build and install the latest release:

```text
$ go install github.com/bezmoradi/gopamin@latest
```

This installs it into `$(go env GOPATH)/bin`. If `gopamin version` then says `gopamin: command not found`, add that folder to `PATH` in `~/.bashrc` or `~/.zshrc`, and reload the file with `source`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

To build from source, see [Development](#development).

## Guides

Each new project includes a `README.md` file in the root path which provides you with some guides on how to use that specific project.

## Development

To build from source, with Go 1.27.3 or later, in Gopamin's checkout:

```text
$ go build -o gopamin github.com/bezmoradi/gopamin
$ go build ./... && go vet ./... && go test ./... && test -z "$(gofmt -l .)"
```

`scripts/release.sh <version> <dir>` builds a release folder from a clean checkout: the four archives, the notices, `install.sh` and `SHA256SUMS`. [CLAUDE.md](CLAUDE.md#testing) says how to check a change by generating real projects.

### Releasing

Only the owner (bezmoradi) can release: the release workflow refuses a tag pushed by anyone else, one not on `main`, or one that doesn't equal `VERSION`.

1. Set `VERSION` in `internal/commands/constants.go` to the new version (`vX.Y.Z`), and merge that to `main`.
2. Tag `main` and push the tag:

   ```text
   $ git fetch origin && git tag -a vX.Y.Z -m "Gopamin X.Y.Z" origin/main && git push origin vX.Y.Z
   ```

3. Watch the release workflow:

   ```text
   $ gh run watch "$(gh run list --workflow release.yml --limit 1 --json databaseId --jq '.[0].databaseId')"
   ```

4. Check the draft with `gh release view vX.Y.Z`: it lists the four archives, the NOTICES file, `install.sh` and `SHA256SUMS`.
5. Publish it with `gh release edit vX.Y.Z --draft=false` (or on the releases page). From then on the one-command install installs it and older versions announce it.

## Author

This project is maintained by [Bez Moradi](https://github.com/bezmoradi)

## License

Gopamin is licensed under [MIT](https://github.com/bezmoradi/gopamin/blob/main/LICENSE)
