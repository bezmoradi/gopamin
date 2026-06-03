## Table of Contents

- [Introduction](#introduction)
- [An Intro to The Clean Architecture](#an-intro-to-the-clean-architecture)
- [Prerequisites](#prerequisites)
- [Supported Flags](#supported-flags)
- [Recipes](#recipes)
- [Installation](#installation)
- [Update](#update)
- [Guides](#guides)
- [Author](#author)
- [License](#license)

## Introduction

Gopamin is a CLI which creates new projects based on ideas promoted by [Standard Go Project Layout](https://github.com/golang-standards/project-layout) which is a **well-accepted** architecture by the Go community (It's not an official standard defined by the core Go team; however, it is a set of common historical and emerging project layout patterns in the Go ecosystem).

## An Intro to The Clean Architecture

All boilerplates created by Gopamin are based on [The Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html). Simply put, it is a software architectural pattern introduced by Robert C. Martin in order to create maintainable, scalable, and testable software systems by decoupling the application logic from the infrastructure details (Although there are some minor differences, this architectural pattern is referred to by other names such as Ports & Adapters Architecture, Hexagonal Architecture, and Onion Architecture).

## Prerequisites

The minimum required tools for using the Gopamin CLI tool is Golang v1.26.0 or higher which can be downloaded from [Go All Releases](https://go.dev/dl). To have full development setup though, other tools are also recommended to be installed on your local machine:

- **[Git](https://git-scm.com/)**: By default each new project created by this tool initializes a Git repo; that's why you need to make sure Git is installed on your machine.
- **[Docker](https://www.docker.com)**: If you choose to create a new project with database integration, a `docker-compose.yml` will be included in the root of the project for running the database of your choice.
- **[GNU Make](https://www.gnu.org/software/make)**: This a tool which controls the generation of executables and other types of files. By default, each new project includes a `Makefile` for running some most-used commands like running an application (This tool is installed by default on Mac and some distributions of Linux). To check whether this tool is installed on your machine, open terminal and run `make --version` (If you do not have this tool installed on your machine though, still you can use this tools without any limitations).

## Supported Flags

A project is defined by **orthogonal capabilities**: the **type** (`-t`) picks the app's primary interface (HTML, JSON/GraphQL, or message-driven), while a **framework** (`-f`), a **broker** (`-b`), and a **database** (`-d`) are optional capabilities you attach to it. You can mix them freely.

| Flag | Description | Required | Options |
|------|-------------|----------|---------|
| `-n`, `--name` | Project name, or a full module path | Always | Any string; characters outside `A-Z a-z 0-9 . - _` are replaced with `-`. A module path such as `github.com/you/repo` also sets the `go.mod` module. |
| `-t`, `--type` | The app's primary interface | Always | `hello-world`, `web-app`, `api`, `worker` |
| `-f`, `--framework` | HTTP framework | `api` & `web-app` only | `echo`, `chi`, `gin`, `httprouter`, `gorilla`, `http` — plus `graphql` for `api` |
| `-b`, `--broker` | Message broker — publishes an event on every write and runs a consumer loop | Required for `worker`; optional for `api` / `web-app` | `kafka`, `rabbitmq`, `redis` |
| `-d`, `--database` | Database | Optional (for any type) | `mysql`, `mariadb`, `postgres`, `mongodb`, `cassandra`, `sqlite`, `dynamodb`, `badgerdb` |
| `-l`, `--logger` | Logging library | Always | `log`, `slog`, `logrus`, `zap` |

A few things worth knowing:

- A `worker` is a **headless, broker-driven service** — it has no HTTP server (so it takes no `-f`) and is triggered by messages off its broker.
- Adding `-b` to an `api` or `web-app` gives you one binary that runs the HTTP server **and** a broker consumer side by side, publishing an event on every write. In other words, "an API that is also a worker" is `-t api -f <framework> -b <broker>`.
- When a containerized database and/or a broker is selected, a single merged `docker-compose.yml` (plus `docker-up` / `docker-down` make targets) is generated so every backend runs together for local development.
- Redis is a broker/cache, offered only via `-b` — it is intentionally **not** a database option.

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
| `gopamin new -n app -t worker -b rabbitmq -d postgres -l slog` | Headless worker: consume a message → persist it |
| `gopamin new -n app -t worker -b kafka -l logrus` | Headless worker with no database |

## Installation

To install the Gopamin tool, run the following command inside terminal:

```text
$ go install github.com/bezmoradi/gopamin@latest
```

To make sure it's installed correctly in your `$GOPATH`, run the following command:

```text
$ gopamin version
```

If this prints a version number, you're all set. If you instead see `gopamin: command not found`, you need to configure `PATH`. Add it as shown below, then re-run `gopamin version`. For Bash, add this line to `~/.bashrc`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

Then reload:

```bash
source ~/.bashrc
```

For Zsh, add it to `~/.zshrc`:

```bash
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc

```

Then reload:

```bash
source ~/.zshrc
```

If the installation process goes well, from now on you can run the `gopamin` command from anywhere on your file system.

## Update

Since this tool is still in development, bugs are being identified and addressed while new features are continuously integrated. To ensure users benefit from these bug fixes and additions, it's essential to have the latest version installed. Hence, a version check mechanism has been implemented within the tool. This mechanism automatically compares the installed version on your device with the latest release whenever the tool is accessed. If a disparity is detected between the two versions, you'll receive a prompt to update to the latest version before proceeding. As an example we have:

```text
The newest version of the Gopamin CLI is 1.0.13 but the installed version on your system is v1.0.10.
To get the latest features and likely bugfixes, please install the latest version by running 'go install github.com/bezmoradi/gopamin@1.0.13'
```

## Guides

Each new project includes a `README.md` file in the root path which provides you with some guides on how to use that specific project.

## Author

This project is maintained by [Bez Moradi](https://github.com/bezmoradi)

## License

Gopamin is licensed under [MIT](https://github.com/bezmoradi/gopamin/blob/master/LICENSE)
