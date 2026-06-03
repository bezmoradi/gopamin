## Table of Contents

-   [Introduction](#introduction)
-   [An Intro to The Clean Architecture](#an-intro-to-the-clean-architecture)
-   [Prerequisites](#prerequisites)
-   [Supported Application Types](#supported-application-types)
-   [Supported Message Brokers](#supported-message-brokers)
-   [Supported Databases](#supported-databases)
-   [Supported Logging Libraries](#supported-logging-libraries)
-   [Installation](#installation)
-   [Update](#update)
-   [Usage](#usage)
-   [Guides](#guides)
-   [Author](#author)
-   [License](#license)

## Introduction

Gopamin is a CLI which creates new projects based on ideas promoted by [Standard Go Project Layout](https://github.com/golang-standards/project-layout) which is a **well-accepted** architecture by the Go community (It's not an official standard defined by the core Go team; however, it is a set of common historical and emerging project layout patterns in the Go ecosystem).

## An Intro to The Clean Architecture

All boilerplates created by Gopamin are based on [The Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html). Simply put, it is a software architectural pattern introduced by Robert C. Martin in order to create maintainable, scalable, and testable software systems by decoupling the application logic from the infrastructure details (Although there are some minor differences, this architectural pattern is referred to by other names such as Ports & Adapters Architecture, Hexagonal Architecture, and Onion Architecture).

## Prerequisites

The minimum required tools for using the Gopamin CLI tool is Golang v1.26.0 or higher which can be downloaded from [Go All Releases](https://go.dev/dl). To have full development setup though, other tools are also recommended to be installed on your local machine:

-   **[Git](https://git-scm.com/)**: By default each new project created by this tool initializes a Git repo; that's why you need to make sure Git is installed on your machine.
-   **[Docker](https://www.docker.com)**: If you choose to create a new project with database integration, a `docker-compose.yml` will be included in the root of the project for running the database of your choice.
-   **[GNU Make](https://www.gnu.org/software/make)**: This a tool which controls the generation of executables and other types of files. By default, each new project includes a `Makefile` for running some most-used commands like running an application (This tool is installed by default on Mac and some distributions of Linux). To check whether this tool is installed on your machine, open terminal and run `make --version` (If you do not have this tool installed on your machine though, still you can use this tools without any limitations).

## Supported Application Types

You can create a range of different Golang applications using the Gopamin tool, from a simple Hello World app to an event-driven worker. A project is defined by orthogonal capabilities: its **type** (`-t`), an optional HTTP **framework** (`-f`), an optional message **broker** (`-b`), and an optional **database** (`-d`). Supported application types are as follows:

-   Hello World
-   Web Application (requires an HTTP framework via `-f`)
-   RESTful API (requires an HTTP framework via `-f`)
    -   Echo
    -   Chi
    -   Gin
    -   Httprouter
    -   Gorilla
    -   HTTP (the built-in `net/http` package will be used)
-   GraphQL API (an `api` project created with `-f graphql`)
-   Worker (a broker-driven service with no HTTP server; requires a broker via `-b`)

> A "microservice" is not a separate type — it's a deployment style. Any of the above (a REST API, a GraphQL API, or a worker) can be a microservice, and any of them can carry a database. An `api`/`web-app` given a `-b` broker emits an event on every write **and** runs a consumer loop alongside its HTTP server.

## Supported Message Brokers

By passing the `-b` flag you add messaging. It is **required** for the `worker` type and **optional** for `api` and `web-app` projects. Supported brokers are as follows:

-   Kafka
-   RabbitMQ
-   Redis

## Supported Databases

By passing the `-d` flag you can add a database to any project type — it is always optional. Supported databases are as follows:

-   MySQL
-   MariaDB
-   Cassandra
-   PostgreSQL
-   MongoDB
-   SQLite
-   DynamoDB
-   BadgerDB

When a containerized database and/or a broker is selected, a single merged `docker-compose.yml` (plus generic `docker-up` / `docker-down` make targets) is generated so you can run every backend together for local development. (Redis is intentionally **not** a database option — it is a cache/broker, offered only via `-b`.)

## Supported Logging Libraries

No matter what type of project you want to scaffold, the `-l` flag must be passed in order to choose the logger type. Supported loggers are as follows:

-   Log
-   Slog
-   Logrus
-   Zap

## Installation

To install the Gopamin tool, run the following command inside terminal:

```text
$ go install github.com/bezmoradi/gopamin@latest
```

To make sure it's installed correctly in your `$GOPATH`, run the following command:

```text
$ gopamin version
```

If this prints a version number, you're all set. If you instead see `gopamin: command not found`, Go installed the binary into its bin directory (`$(go env GOPATH)/bin`), which isn't on your `PATH` yet. Add it as shown below, then re-run `gopamin version`.

### For Bash

Add this line to `~/.bashrc`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

Then reload:

```bash
source ~/.bashrc
```

### For Zsh

Add it to `~/.zshrc`:

```bash
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
source ~/.zshrc
```

If the installation process goes well, from now on you can run the `gopamin` command from anywhere on your file system.

## Update

Since this tool is still in development, bugs are being identified and addressed while new features are continuously integrated. To ensure users benefit from these bug fixes and additions, it's essential to have the latest version installed. Hence, a version check mechanism has been implemented within the tool. This mechanism automatically compares the installed version on your device with the latest release whenever the tool is accessed. If a disparity is detected between the two versions, you'll receive a prompt to update to the latest version before proceeding. As an example we have:

```text
The newest version of the Gopamin CLI is 1.0.13 but the installed version on your system is v1.0.10. 
To get the latest features and likely bugfixes, please install the latest version by running 'go install github.com/bezmoradi/gopamin@1.0.13'
```

## Usage

If you got the current version correctly in the previous step, now you can use this tools to scaffold new Golang projects.

### The `-h` Flag

The `-h` flag which is short for `--help` can be used to show you a guide on how to create a new project:

```text
$ gopamin new -h
```

The above command shows the help on how to use different flags to scaffold different types of projects. In the following sections, these flags will be introduced.

### The `-n` Flag

The `-n` flag which is short for `--name` should be used for choosing a name for your project. For example:

```text
$ gopamin new -n my-hello-world-app -t hello-world -l log
```

Only ASCII letters, digits, and the characters `.`, `-`, and `_` are accepted and any other character will be replaced by `-`. You can also pick your repository name as follows:

```text
$ gopamin new -n github.com/user-name/repo_name -t hello-world -l log
```

In this case, the project folder will be named `github.com-user-name-repo_name` but the name of the module inside the `go.mod` will be as follows:

```text
module github.com/user-name/repo_name

go 1.26.0
```

### The `-t` Flag

The `-t` flag, short for `--type`, chooses the project type to scaffold. For example, to create a simple Hello World app we have:

```text
$ gopamin new -n my-hello-world-app -t hello-world -l log
```

The supported values for the `-t` flag are `hello-world`, `web-app`, `api`, and `worker`.

### The `-f` Flag

The `-f` flag, short for `--framework`, selects the HTTP framework and is **required** for the `web-app` and `api` types. Supported values for `web-app` are `echo`, `chi`, `gin`, `httprouter`, `gorilla`, and `http`. For `api`, the same values plus `graphql` are supported. For example, to create an API that uses the Echo framework under the hood we have:

```text
$ gopamin new -n my-rest-api -t api -f echo -l log
```

Or in order to create a GraphQL API with MySQL integration we have:

```text
$ gopamin new -n my-graphql-api -t api -f graphql -d mysql -l log
```

### The `-b` Flag

The `-b` flag, short for `--broker`, selects a message broker. It is **required** for the `worker` type and **optional** for `api` and `web-app`. Supported values are `kafka`, `rabbitmq`, and `redis`. For example, to create a broker-driven worker with Kafka we have:

```text
$ gopamin new -n my-worker -t worker -b kafka -l log
```

When you add a broker to an `api` or `web-app`, the generated service produces an event on every write and runs a consumer loop alongside the HTTP server — and it can still own a database:

```text
$ gopamin new -n my-service -t api -f echo -b kafka -d postgres -l slog
```

### The `-d` Flag

The `-d` flag, short for `--database`, adds a database and is optional for every type. Available values are `mysql`, `cassandra`, `mariadb`, `postgres`, `mongodb`, `sqlite`, `dynamodb`, and `badgerdb`. For example, to create a web application with MySQL integration we have:

```text
$ gopamin new -n my-web-app -t web-app -f http -d mysql -l log
```

### The `-l` Flag

The `-l` flag, short for `--logger`, sets the logger and is required for every type. Available values are `log`, `slog`, `logrus`, and `zap`. For example, to create a worker with Logrus support we have:

```text
$ gopamin new -n my-worker -t worker -b rabbitmq -l logrus
```

## Guides

Each new project includes a `README.md` file in the root path which provides you with some guides on how to use that specific project.

## Author

This project is maintained by [Bez Moradi](https://github.com/bezmoradi)

## License

Gopamin is licensed under [MIT](https://github.com/bezmoradi/gopamin/blob/master/LICENSE)
