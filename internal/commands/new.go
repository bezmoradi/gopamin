package commands

import (
	"github.com/bezmoradi/gopamin/internal/scaffolder"
	"github.com/spf13/cobra"
)

var (
	database    string
	platform    string
	name        string
	projectType string
	logger      string
)

var newCmd = &cobra.Command{
	Example: `To create a simple hello world app, run the following command:
gopamin new -n HelloWorld -t hello-world
To create a hello world app with MySQL integration, run the following command:
 gopamin new -n HelloWorld -t hello-world -d mysql
To create a web application using the built-in http package and MySQL support, run the following command:
 gopamin new -n HelloWorld -t web-app -p http -d mysql
To create a RESTful API using the echo framework, run the following command:
 gopamin new -n HelloWorld -t api -p echo
To create a RESTful API using the built-in http package and MongoDB integration, run the following command:
 gopamin new -n HelloWorld -t api -p http -d mongodb
To create a Microservice using Kafka message broker, run the following command:
 gopamin new -n HelloWorld -t microservice -p kafka`,
	Use:   "new",
	Short: "Create a new project",
	Run: func(cmd *cobra.Command, args []string) {
		if argsValidator() {
			scaffolder.New(projectType, platform, name, database, logger)
		}
	},
}

func init() {
	newCmd.Flags().StringVarP(&name, "name", "n", "", `Name of project (Keep in mind the all white spaces will be replaced by dash). 
If you want to use space-separated words, place them inside double quotes like "my demo app" then it will be converted to "my-demo-app".`)

	newCmd.Flags().StringVarP(&projectType, "type", "t", "", `Type of the project. Available types are:
 - hello-world (A simple "Hello World" app without any extra functionalities. Database integration can also be added).
 - web-app (If chosen, you must also add the "-p" flag to specify the type of the third-party platform you want to use).
 - microservice (If chosen, you must also add the "-p" flag to specify the type of the third-party platform you want to use).
 - api (If chosen, you must also add the "-p" flag to specify the type of the third-party platform you want to use).`)

	newCmd.Flags().StringVarP(&platform, "platform", "p", "", `If the chosen "-t" is "api", "web-app", or "microservice", you must also add the "-p" flag to specify the underlying third-party platform.
Available values for the "api" type are:
 - echo
 - chi
 - gin
 - httprouter
 - gorilla
 - http (The build-in HTTP package will be used)
 - graphql
Available values for the "web-app" type are:
 - echo
 - chi
 - gin
 - httprouter
 - gorilla
 - http (The build-in HTTP package will be used)
Available values for the "microservice" type are:
- redis
- kafka
- rabbitmq`)

	newCmd.Flags().StringVarP(&database, "database", "d", "", `Type of the database. Available values are:
 - mysql
 - mariadb
 - cassandra
 - postgres
 - mongodb
 - sqlite
 - dynamodb
 - badgerdb
 - redis`)

	newCmd.Flags().StringVarP(&logger, "logger", "l", "", `Type of the logger. Available values are:
 - log
 - slog
 - logrus
 - zap`)

	rootCmd.AddCommand(newCmd)
}
