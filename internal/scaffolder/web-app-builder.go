package scaffolder

import "fmt"

// webRecipe captures what varies between web-app platforms.
type webRecipe struct {
	routes   string
	users    string
	packages []string
}

var webRecipes = map[string]webRecipe{
	"http":       {"web-app-http-routes", "web-app-http-users", nil},
	"chi":        {"web-app-chi-routes", "web-app-chi-users", []string{"github.com/go-chi/chi/v5"}},
	"echo":       {"web-app-echo-routes", "web-app-echo-users", []string{"github.com/labstack/echo/v5"}},
	"gin":        {"web-app-gin-routes", "web-app-gin-users", []string{"github.com/gin-gonic/gin"}},
	"gorilla":    {"web-app-gorilla-routes", "web-app-gorilla-users", []string{"github.com/gorilla/mux"}},
	"httprouter": {"web-app-httprouter-routes", "web-app-httprouter-users", []string{"github.com/julienschmidt/httprouter"}},
}

func buildWebApp(p *Project) {
	recipe := webRecipes[p.Platform]
	env := []string{"env", "web-app-env"}
	readme := []string{"readme", p.Logger + "-readme"}
	makefile := []string{"makefile"}

	database := p.Database
	if database == "" {
		database = "mock"
		fileGenerator([]string{"web-app-main"}, p)
		readme = append(readme, "web-app-readme")
	} else {
		fileGenerator([]string{"web-app-main-with-db"}, p)
		readme = append(readme, database+"-readme", "web-app-readme-with-db")
		env = append(env, database+"-env")
	}

	if hasCompose(p) {
		makefile = append(makefile, "docker-makefile")
	}
	if isRelational(p.Database) {
		makefile = append(makefile, "migrate-makefile")
	}
	if p.Broker != "" {
		env = append(env, p.Broker+"-microservice-env")
	}
	if hasObservability(p) {
		env = append(env, "otel-env")
	}

	buildDatabase(p, database)
	buildLogger(p)
	buildObservability(p)
	if p.Broker != "" {
		fileGenerator([]string{p.Broker + "-microservice-broker"}, p)
		goGetPackages(p.Path, workerRecipes[p.Broker])
	}
	fileGenerator(env, p)
	fileGenerator(readme, p)
	fileGenerator(makefile, p)
	if hasCompose(p) {
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
	fileGenerator([]string{"web-app-server"}, p)
	fileGenerator([]string{recipe.routes}, p)
	fileGenerator([]string{recipe.users}, p)
	fileGenerator([]string{"health-handler"}, p)
	fileGenerator([]string{"middleware"}, p)
	fileGenerator([]string{"router-interface"}, p)
	fileGenerator([]string{"web-app-styles"}, p)
	fileGenerator([]string{"web-app-users-html-template"}, p)
	fileGenerator([]string{"web-app-user-html-template"}, p)
	fileGenerator([]string{"web-app-embed"}, p)
	goGetPackages(p.Path, recipe.packages)

	fmt.Printf("%v "+BUILD_SUCCESS_MESSAGE+"\n", p.Name)
}
