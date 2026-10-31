package scaffolder

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

func buildWebApp(g *generator) {
	recipe := webRecipes[g.p.Platform]
	g.env = append(g.env, "web-app-env")
	g.readme = append(g.readme, "web-app-readme")

	database := g.p.Database
	if database == "" {
		database = "memory"
	}
	buildDatabase(g, database)

	for _, key := range []string{
		"http-main",
		"http-server",
		recipe.routes,
		recipe.users,
		"health-handler",
		"middleware",
		"web-app-styles",
		"web-app-users-html-template",
		"web-app-user-html-template",
		"web-app-error-html-template",
		"web-app-embed",
	} {
		g.file(key)
	}
	g.need(recipe.packages...)
}
