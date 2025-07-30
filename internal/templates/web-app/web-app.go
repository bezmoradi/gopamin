package templates

import (
	_ "embed"
)

//go:embed files/assets/users.tmpl
var webAppUsersHtmlTemplate []byte

func WebAppUsersHtmlTemplateTemplate() ([]byte, string) {
	return webAppUsersHtmlTemplate, "assets/users.html"
}

//go:embed files/assets/user.tmpl
var webAppUserHtmlTemplate []byte

func WebAppUserHtmlTemplateTemplate() ([]byte, string) {
	return webAppUserHtmlTemplate, "assets/user.html"
}

//go:embed files/assets/styles.tmpl
var webAppStyles []byte

func WebAppStylesTemplate() ([]byte, string) {
	return webAppStyles, "assets/styles.css"
}

//go:embed files/env.tmpl
var webAppEnv []byte

func WebAppEnvTemplate() ([]byte, string) {
	return webAppEnv, ".env"
}

//go:embed files/readme/readme.tmpl
var webAppReadme []byte

func WebAppReadmeTemplate() ([]byte, string) {
	return webAppReadme, "README.md"
}

//go:embed files/readme/readme-with-db.tmpl
var readme []byte

func WebAppReadmeWithDbTemplate() ([]byte, string) {
	return readme, "README.md"
}

//go:embed files/cmd/server/server.tmpl
var webAppServer []byte

func WebAppServerTemplate() ([]byte, string) {
	return webAppServer, "cmd/server/server.go"
}

//go:embed files/internal/adapters/handlers/http/routes.tmpl
var webAppHttpRoutes []byte

func WebAppHttpRoutesTemplate() ([]byte, string) {
	return webAppHttpRoutes, "internal/adapters/handlers/http/routes.go"
}

//go:embed files/internal/adapters/handlers/http/users.tmpl
var webAppHttpUsers []byte

func WebAppHttpUsersTemplate() ([]byte, string) {
	return webAppHttpUsers, "internal/adapters/handlers/http/users.go"
}

//go:embed files/internal/adapters/handlers/chi/routes.tmpl
var webAppChiRoutes []byte

func WebAppChiRoutesTemplate() ([]byte, string) {
	return webAppChiRoutes, "internal/adapters/handlers/chi/routes.go"
}

//go:embed files/internal/adapters/handlers/chi/users.tmpl
var webAppHttpChiUsers []byte

func WebAppChiUsersTemplate() ([]byte, string) {
	return webAppHttpChiUsers, "internal/adapters/handlers/chi/users.go"
}

//go:embed files/internal/adapters/handlers/echo/routes.tmpl
var webAppEchoRoutes []byte

func WebAppEchoRoutesTemplate() ([]byte, string) {
	return webAppEchoRoutes, "internal/adapters/handlers/echo/routes.go"
}

//go:embed files/internal/adapters/handlers/echo/users.tmpl
var webAppHttpEchoUsers []byte

func WebAppEchoUsersTemplate() ([]byte, string) {
	return webAppHttpEchoUsers, "internal/adapters/handlers/echo/users.go"
}

//go:embed files/internal/adapters/handlers/gin/routes.tmpl
var webAppGinRoutes []byte

func WebAppGinRoutesTemplate() ([]byte, string) {
	return webAppGinRoutes, "internal/adapters/handlers/gin/routes.go"
}

//go:embed files/internal/adapters/handlers/gin/users.tmpl
var webAppHttpGinUsers []byte

func WebAppGinUsersTemplate() ([]byte, string) {
	return webAppHttpGinUsers, "internal/adapters/handlers/gin/users.go"
}

//go:embed files/internal/adapters/handlers/gorilla/routes.tmpl
var webAppGorillaRoutes []byte

func WebAppGorillaRoutesTemplate() ([]byte, string) {
	return webAppGorillaRoutes, "internal/adapters/handlers/gorilla/routes.go"
}

//go:embed files/internal/adapters/handlers/gorilla/users.tmpl
var webAppHttpGorillaUsers []byte

func WebAppGorillaUsersTemplate() ([]byte, string) {
	return webAppHttpGorillaUsers, "internal/adapters/handlers/gorilla/users.go"
}

//go:embed files/internal/adapters/handlers/httprouter/routes.tmpl
var webAppHttprouterRoutes []byte

func WebAppHttprouterRoutesTemplate() ([]byte, string) {
	return webAppHttprouterRoutes, "internal/adapters/handlers/httprouter/routes.go"
}

//go:embed files/internal/adapters/handlers/httprouter/users.tmpl
var webAppHttpHttprouterUsers []byte

func WebAppHttprouterUsersTemplate() ([]byte, string) {
	return webAppHttpHttprouterUsers, "internal/adapters/handlers/httprouter/users.go"
}

//go:embed files/cmd/main.tmpl
var webAppMain []byte

func WebAppMainTemplate() ([]byte, string) {
	return webAppMain, "cmd/main.go"
}

//go:embed files/cmd/main-with-db.tmpl
var webAppMainWithDb []byte

func WebAppMainWithDbTemplate() ([]byte, string) {
	return webAppMainWithDb, "cmd/main.go"
}
