package scaffolder

import "fmt"

// apiRecipe captures what varies between API platforms: the cmd/main and README
// variants (with and without a database), the platform-specific files, and the
// modules to fetch.
type apiRecipe struct {
	mainNoDB     string
	mainWithDB   string
	readmeNoDB   string
	readmeWithDB string
	files        []string
	packages     []string
}

func restAPIRecipe(routes, users string, packages []string) apiRecipe {
	return apiRecipe{
		mainNoDB:     "api-main",
		mainWithDB:   "api-main-with-db",
		readmeNoDB:   "api-readme",
		readmeWithDB: "api-readme-with-db",
		files:        []string{"api-server", routes, users, "api-errors", "api-response", "router-interface"},
		packages:     packages,
	}
}

var apiRecipes = map[string]apiRecipe{
	"echo":       restAPIRecipe("api-echo-routes", "api-echo-users", []string{"github.com/labstack/echo/v5"}),
	"chi":        restAPIRecipe("api-chi-routes", "api-chi-users", []string{"github.com/go-chi/chi/v5"}),
	"gin":        restAPIRecipe("api-gin-routes", "api-gin-users", []string{"github.com/gin-gonic/gin"}),
	"httprouter": restAPIRecipe("api-httprouter-routes", "api-httprouter-users", []string{"github.com/julienschmidt/httprouter"}),
	"gorilla":    restAPIRecipe("api-gorilla-routes", "api-gorilla-users", []string{"github.com/gorilla/mux"}),
	"http":       restAPIRecipe("api-http-routes", "api-http-users", nil),
	"graphql": {
		mainNoDB:     "graphql-main",
		mainWithDB:   "graphql-main-with-db",
		readmeNoDB:   "graphql-readme",
		readmeWithDB: "graphql-readme-with-db",
		files:        []string{"graphql-server", "graphql-schema"},
		packages:     []string{"github.com/graphql-go/graphql", "github.com/graphql-go/handler"},
	},
}

func buildAPI(p *Project) {
	recipe := apiRecipes[p.Platform]
	env := []string{"env", "api-env"}
	readme := []string{"readme", p.Logger + "-readme"}
	makefile := []string{"makefile"}

	database := p.Database
	if database == "" {
		database = "mock"
		fileGenerator([]string{recipe.mainNoDB}, p)
		readme = append(readme, recipe.readmeNoDB)
	} else {
		fileGenerator([]string{recipe.mainWithDB}, p)
		readme = append(readme, database+"-readme", recipe.readmeWithDB)
		env = append(env, database+"-env")
	}

	if hasContainerizedBackend(p) {
		makefile = append(makefile, "docker-makefile")
	}
	if p.Broker != "" {
		env = append(env, p.Broker+"-microservice-env")
	}

	buildDatabase(p, database)
	buildLogger(p)
	if p.Broker != "" {
		fileGenerator([]string{p.Broker + "-microservice-broker"}, p)
		goGetPackages(p.Path, workerRecipes[p.Broker])
	}
	fileGenerator(env, p)
	fileGenerator(readme, p)
	fileGenerator(makefile, p)
	if hasContainerizedBackend(p) {
		fileGenerator([]string{"docker-compose"}, p)
	}
	fileGenerator([]string{"configs"}, p)
	fileGenerator([]string{"configs-test"}, p)
	fileGenerator([]string{"tools"}, p)
	fileGenerator([]string{"tools-test"}, p)
	fileGenerator([]string{"arch-lint"}, p)
	fileGenerator([]string{"golangci"}, p)
	fileGenerator([]string{"dockerfile"}, p)
	fileGenerator([]string{"dockerignore"}, p)
	for _, key := range recipe.files {
		fileGenerator([]string{key}, p)
	}
	fileGenerator([]string{"health-handler"}, p)
	goGetPackages(p.Path, recipe.packages)

	fmt.Printf("%v "+BUILD_SUCCESS_MESSAGE+"\n", p.Name)
}
