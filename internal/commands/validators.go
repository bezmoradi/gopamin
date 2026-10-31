package commands

import (
	"fmt"
	"slices"

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

	if !scaffolder.IsValidProjectName(name) {
		fmt.Println(`The -n value must start with a letter or a digit. For more help, type "gopamin new -h"`)
		return false
	}

	if scaffolder.IsProjectNameTaken(name) {
		fmt.Printf("%q already exists in the current directory\n", scaffolder.DirectoryName(name))
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

	if auth != "" {
		if projectType != "api" {
			fmt.Println(`The -a flag is only supported for projects of type "api"`)
			return false
		}
		if !authValidator() {
			fmt.Println(`The specified value for the -a flag is wrong. For more help, type "gopamin new -h"`)
			return false
		}
	}

	if openapi != "" {
		if projectType != "api" || framework == "graphql" {
			fmt.Println(`The -s flag is only supported for REST "api" projects (not graphql)`)
			return false
		}
		if !openapiValidator() {
			fmt.Println(`The specified value for the -s flag is wrong. For more help, type "gopamin new -h"`)
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
	return slices.Contains(scaffolder.ProjectTypes(), projectType)
}

func webAppFrameworkValidator() bool {
	return slices.Contains(scaffolder.Frameworks("web-app"), framework)
}

func apiFrameworkValidator() bool {
	return slices.Contains(scaffolder.Frameworks("api"), framework)
}

func brokerValidator() bool {
	return slices.Contains(scaffolder.Brokers(), broker)
}

func databaseValidator() bool {
	return slices.Contains(scaffolder.Databases(), database)
}

func loggerValidator() bool {
	return slices.Contains(scaffolder.Loggers(), logger)
}

func ciValidator() bool {
	return slices.Contains(scaffolder.CIs(), ci)
}

func observabilityValidator() bool {
	switch observability {
	case "otel":
		return true
	}

	return false
}

func authValidator() bool {
	switch auth {
	case "jwt":
		return true
	}

	return false
}

func openapiValidator() bool {
	switch openapi {
	case "swagger":
		return true
	}

	return false
}
