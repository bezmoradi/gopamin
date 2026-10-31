package scaffolder

// apiRecipe captures what varies between API platforms: the README section, the
// platform-specific files, and the modules to fetch.
type apiRecipe struct {
	readme   string
	files    []string
	packages []string
}

func restAPIRecipe(routes, users string, packages []string) apiRecipe {
	return apiRecipe{
		readme:   "api-readme",
		files:    []string{routes, users, "api-errors", "api-response", "api-user"},
		packages: packages,
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
		readme:   "graphql-readme",
		files:    []string{"graphql-schema"},
		packages: []string{"github.com/graphql-go/graphql", "github.com/graphql-go/handler"},
	},
}

func buildAPI(g *generator) {
	recipe := apiRecipes[g.p.Platform]
	g.env = append(g.env, "api-env")
	g.readme = append(g.readme, recipe.readme)

	database := g.p.Database
	if database == "" {
		database = "mock"
	}
	buildDatabase(g, database)

	g.file("http-main")
	g.file("http-server")
	for _, key := range recipe.files {
		g.file(key)
	}
	g.file("health-handler")
	g.file("middleware")
	g.need(recipe.packages...)
}
