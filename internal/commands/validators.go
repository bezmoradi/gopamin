package commands

import (
	"fmt"

	"github.com/bezmoradi/gopamin/internal/scaffolder"
)

func argsValidator() bool {
	switch {
	case name == "":
		fmt.Println(`The -n flag is required. For more help, type "gopamin new -h"`)
		return false
	case logger == "":
		fmt.Println(`The -l flag is required. For more help, type "gopamin new -h"`)
		return false
	case projectType == "":
		fmt.Println(`The -t flag is required. For more help, type "gopamin new -h"`)
		return false
	}

	if !validateType() {
		fmt.Println(`The specified value for the -t flag is wrong. For more help, type "gopamin new -h"`)
		return false
	}

	if ok, message := flagsAllowedForType(); !ok {
		fmt.Println(message)
		return false
	}

	if scaffolder.IsProjectNameTaken(name) {
		fmt.Println("This project name is already taken in the current directory")
		return false
	}

	if projectType == "api" && !apiFrameworkValidator() {
		fmt.Println(`The specified value for the -f flag for "api" projects is wrong. For more help, type "gopamin new -h"`)
		return false
	}

	if projectType == "web-app" && !webAppFrameworkValidator() {
		fmt.Println(`The specified value for the -f flag for "web-app" projects is wrong. For more help, type "gopamin new -h"`)
		return false
	}

	if broker != "" && !brokerValidator() {
		fmt.Println(`The specified value for the -b flag is wrong. For more help, type "gopamin new -h"`)
		return false
	}

	if database != "" && !databaseValidator() {
		fmt.Println(`The specified value for the -d flag is wrong. For more help, type "gopamin new -h"`)
		return false
	}

	if !loggerValidator() {
		fmt.Println(`The specified value for the -l flag is wrong. For more help, type "gopamin new -h"`)
		return false
	}

	if ci != "" && !ciValidator() {
		fmt.Println(`The specified value for the -c flag is wrong. For more help, type "gopamin new -h"`)
		return false
	}

	if observability != "" {
		if projectType == "hello-world" {
			fmt.Println(`The -o flag is not allowed for projects of type "hello-world"`)
			return false
		}
		if !observabilityValidator() {
			fmt.Println(`The specified value for the -o flag is wrong. For more help, type "gopamin new -h"`)
			return false
		}
	}

	return true
}

// flagsAllowedForType enforces which capability flags each project type requires
// or forbids: -f (HTTP framework) for api/web-app, -b (broker) for worker.
func flagsAllowedForType() (bool, string) {
	switch projectType {
	case "hello-world":
		if framework != "" {
			return false, `The -f flag is not allowed for projects of type "hello-world"`
		}
		if broker != "" {
			return false, `The -b flag is not allowed for projects of type "hello-world"`
		}
	case "api", "web-app":
		if framework == "" {
			return false, `The -f flag is required for projects of type "api" and "web-app"`
		}
	case "worker":
		if framework != "" {
			return false, `The -f flag is not allowed for projects of type "worker"`
		}
		if broker == "" {
			return false, `The -b flag is required for projects of type "worker"`
		}
	}

	return true, ""
}

func validateType() bool {
	switch projectType {
	case "hello-world", "web-app", "api", "worker":
		return true
	}

	return false
}

func webAppFrameworkValidator() bool {
	switch framework {
	case "echo", "chi", "gin", "http", "gorilla", "httprouter":
		return true
	}

	return false
}

func apiFrameworkValidator() bool {
	switch framework {
	case "echo", "chi", "gin", "http", "gorilla", "graphql", "httprouter":
		return true
	}

	return false
}

func brokerValidator() bool {
	switch broker {
	case "kafka", "rabbitmq", "redis":
		return true
	}

	return false
}

func databaseValidator() bool {
	switch database {
	case "mysql", "mariadb", "postgres", "mongodb", "sqlite", "dynamodb", "badgerdb", "cassandra":
		return true
	}

	return false
}

func loggerValidator() bool {
	switch logger {
	case "log", "slog", "logrus", "zap":
		return true
	}

	return false
}

func ciValidator() bool {
	switch ci {
	case "github", "gitlab":
		return true
	}

	return false
}

func observabilityValidator() bool {
	switch observability {
	case "otel":
		return true
	}

	return false
}
