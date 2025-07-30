package scaffolder

import "fmt"

type webAppHttpBuilder struct {
	project *Project
}

func (b *webAppHttpBuilder) build() {
	env := []string{"env", "web-app-env"}
	readme := []string{"readme", b.project.Logger + "-readme"}
	makefile := []string{"makefile"}

	if b.project.Database == "" {
		b.project.Database = "mock"
		fileGenerator([]string{"web-app-main"}, b.project)
		readme = append(readme, "web-app-readme")
	} else {
		fileGenerator([]string{"web-app-main-with-db"}, b.project)
		readme = append(readme, b.project.Database+"-readme", "web-app-readme-with-db")
		makefile = append(makefile, b.project.Database+"-makefile")
		env = append(env, b.project.Database+"-env")
	}

	dbSelector(b.project)
	loggerSelector(b.project)
	fileGenerator(env, b.project)
	fileGenerator(readme, b.project)
	fileGenerator(makefile, b.project)
	fileGenerator([]string{"configs"}, b.project)
	fileGenerator([]string{"configs-test"}, b.project)
	fileGenerator([]string{"tools"}, b.project)
	fileGenerator([]string{"tools-test"}, b.project)
	fileGenerator([]string{"web-app-server"}, b.project)
	fileGenerator([]string{"web-app-http-routes"}, b.project)
	fileGenerator([]string{"web-app-http-users"}, b.project)
	fileGenerator([]string{"router-inferface"}, b.project)
	fileGenerator([]string{"web-app-styles"}, b.project)
	fileGenerator([]string{"web-app-users-html-template"}, b.project)
	fileGenerator([]string{"web-app-user-html-template"}, b.project)

	fmt.Printf("%v "+BUILD_SUCCESS_MESSAGE+"\n", b.project.Name)
}

type webAppChiBuilder struct {
	project *Project
}

func (b *webAppChiBuilder) build() {
	env := []string{"env", "web-app-env"}
	readme := []string{"readme", b.project.Logger + "-readme"}
	makefile := []string{"makefile"}

	if b.project.Database == "" {
		b.project.Database = "mock"
		fileGenerator([]string{"web-app-main"}, b.project)
		readme = append(readme, "web-app-readme")
	} else {
		fileGenerator([]string{"web-app-main-with-db"}, b.project)
		readme = append(readme, b.project.Database+"-readme", "web-app-readme-with-db")
		makefile = append(makefile, b.project.Database+"-makefile")
		env = append(env, b.project.Database+"-env")
	}

	dbSelector(b.project)
	loggerSelector(b.project)
	fileGenerator(env, b.project)
	fileGenerator(readme, b.project)
	fileGenerator(makefile, b.project)
	fileGenerator([]string{"configs"}, b.project)
	fileGenerator([]string{"configs-test"}, b.project)
	fileGenerator([]string{"tools"}, b.project)
	fileGenerator([]string{"tools-test"}, b.project)
	fileGenerator([]string{"web-app-server"}, b.project)
	fileGenerator([]string{"web-app-chi-routes"}, b.project)
	fileGenerator([]string{"web-app-chi-users"}, b.project)
	fileGenerator([]string{"router-inferface"}, b.project)
	fileGenerator([]string{"web-app-styles"}, b.project)
	fileGenerator([]string{"web-app-users-html-template"}, b.project)
	fileGenerator([]string{"web-app-user-html-template"}, b.project)

	goGetPackages(b.project.Path, []string{"github.com/go-chi/chi/v5"})

	fmt.Printf("%v "+BUILD_SUCCESS_MESSAGE+"\n", b.project.Name)
}

type webAppEchoBuilder struct {
	project *Project
}

func (b *webAppEchoBuilder) build() {
	env := []string{"env", "web-app-env"}
	readme := []string{"readme", b.project.Logger + "-readme"}
	makefile := []string{"makefile"}

	if b.project.Database == "" {
		b.project.Database = "mock"
		fileGenerator([]string{"web-app-main"}, b.project)
		readme = append(readme, "web-app-readme")
	} else {
		fileGenerator([]string{"web-app-main-with-db"}, b.project)
		readme = append(readme, b.project.Database+"-readme", "web-app-readme-with-db")
		makefile = append(makefile, b.project.Database+"-makefile")
		env = append(env, b.project.Database+"-env")
	}

	dbSelector(b.project)
	loggerSelector(b.project)
	fileGenerator(env, b.project)
	fileGenerator(readme, b.project)
	fileGenerator(makefile, b.project)
	fileGenerator([]string{"configs"}, b.project)
	fileGenerator([]string{"configs-test"}, b.project)
	fileGenerator([]string{"tools"}, b.project)
	fileGenerator([]string{"tools-test"}, b.project)
	fileGenerator([]string{"web-app-server"}, b.project)
	fileGenerator([]string{"web-app-echo-routes"}, b.project)
	fileGenerator([]string{"web-app-echo-users"}, b.project)
	fileGenerator([]string{"router-inferface"}, b.project)
	fileGenerator([]string{"web-app-styles"}, b.project)
	fileGenerator([]string{"web-app-users-html-template"}, b.project)
	fileGenerator([]string{"web-app-user-html-template"}, b.project)

	goGetPackages(b.project.Path, []string{"github.com/labstack/echo/v4"})

	fmt.Printf("%v "+BUILD_SUCCESS_MESSAGE+"\n", b.project.Name)
}

type webAppGinBuilder struct {
	project *Project
}

func (b *webAppGinBuilder) build() {
	env := []string{"env", "web-app-env"}
	readme := []string{"readme", b.project.Logger + "-readme"}
	makefile := []string{"makefile"}

	if b.project.Database == "" {
		b.project.Database = "mock"
		fileGenerator([]string{"web-app-main"}, b.project)
		readme = append(readme, "web-app-readme")
	} else {
		fileGenerator([]string{"web-app-main-with-db"}, b.project)
		readme = append(readme, b.project.Database+"-readme", "web-app-readme-with-db")
		makefile = append(makefile, b.project.Database+"-makefile")
		env = append(env, b.project.Database+"-env")
	}

	dbSelector(b.project)
	loggerSelector(b.project)
	fileGenerator(env, b.project)
	fileGenerator(readme, b.project)
	fileGenerator(makefile, b.project)
	fileGenerator([]string{"configs"}, b.project)
	fileGenerator([]string{"configs-test"}, b.project)
	fileGenerator([]string{"tools"}, b.project)
	fileGenerator([]string{"tools-test"}, b.project)
	fileGenerator([]string{"web-app-server"}, b.project)
	fileGenerator([]string{"web-app-gin-routes"}, b.project)
	fileGenerator([]string{"web-app-gin-users"}, b.project)
	fileGenerator([]string{"router-inferface"}, b.project)
	fileGenerator([]string{"web-app-styles"}, b.project)
	fileGenerator([]string{"web-app-users-html-template"}, b.project)
	fileGenerator([]string{"web-app-user-html-template"}, b.project)

	goGetPackages(b.project.Path, []string{"github.com/gin-gonic/gin"})

	fmt.Printf("%v "+BUILD_SUCCESS_MESSAGE+"\n", b.project.Name)
}

type webAppGorillaBuilder struct {
	project *Project
}

func (b *webAppGorillaBuilder) build() {
	env := []string{"env", "web-app-env"}
	readme := []string{"readme", b.project.Logger + "-readme"}
	makefile := []string{"makefile"}

	if b.project.Database == "" {
		b.project.Database = "mock"
		fileGenerator([]string{"web-app-main"}, b.project)
		readme = append(readme, "web-app-readme")
	} else {
		fileGenerator([]string{"web-app-main-with-db"}, b.project)
		readme = append(readme, b.project.Database+"-readme", "web-app-readme-with-db")
		makefile = append(makefile, b.project.Database+"-makefile")
		env = append(env, b.project.Database+"-env")
	}

	dbSelector(b.project)
	loggerSelector(b.project)
	fileGenerator(env, b.project)
	fileGenerator(readme, b.project)
	fileGenerator(makefile, b.project)
	fileGenerator([]string{"configs"}, b.project)
	fileGenerator([]string{"configs-test"}, b.project)
	fileGenerator([]string{"tools"}, b.project)
	fileGenerator([]string{"tools-test"}, b.project)
	fileGenerator([]string{"web-app-server"}, b.project)
	fileGenerator([]string{"web-app-gorilla-routes"}, b.project)
	fileGenerator([]string{"web-app-gorilla-users"}, b.project)
	fileGenerator([]string{"router-inferface"}, b.project)
	fileGenerator([]string{"web-app-styles"}, b.project)
	fileGenerator([]string{"web-app-users-html-template"}, b.project)
	fileGenerator([]string{"web-app-user-html-template"}, b.project)

	goGetPackages(b.project.Path, []string{"github.com/gorilla/mux"})

	fmt.Printf("%v "+BUILD_SUCCESS_MESSAGE+"\n", b.project.Name)
}

type webAppHttprouterBuilder struct {
	project *Project
}

func (b *webAppHttprouterBuilder) build() {
	env := []string{"env", "web-app-env"}
	readme := []string{"readme", b.project.Logger + "-readme"}
	makefile := []string{"makefile"}

	if b.project.Database == "" {
		b.project.Database = "mock"
		fileGenerator([]string{"web-app-main"}, b.project)
		readme = append(readme, "web-app-readme")
	} else {
		fileGenerator([]string{"web-app-main-with-db"}, b.project)
		readme = append(readme, b.project.Database+"-readme", "web-app-readme-with-db")
		makefile = append(makefile, b.project.Database+"-makefile")
		env = append(env, b.project.Database+"-env")
	}

	dbSelector(b.project)
	loggerSelector(b.project)
	fileGenerator(env, b.project)
	fileGenerator(readme, b.project)
	fileGenerator(makefile, b.project)
	fileGenerator([]string{"configs"}, b.project)
	fileGenerator([]string{"configs-test"}, b.project)
	fileGenerator([]string{"tools"}, b.project)
	fileGenerator([]string{"tools-test"}, b.project)
	fileGenerator([]string{"web-app-server"}, b.project)
	fileGenerator([]string{"web-app-httprouter-routes"}, b.project)
	fileGenerator([]string{"web-app-httprouter-users"}, b.project)
	fileGenerator([]string{"router-inferface"}, b.project)
	fileGenerator([]string{"web-app-styles"}, b.project)
	fileGenerator([]string{"web-app-users-html-template"}, b.project)
	fileGenerator([]string{"web-app-user-html-template"}, b.project)

	goGetPackages(b.project.Path, []string{"github.com/julienschmidt/httprouter"})

	fmt.Printf("%v "+BUILD_SUCCESS_MESSAGE+"\n", b.project.Name)
}

type webAppBuilderFactory func(p *Project) boilerplateBuilder

var webAppBuilderMap = map[string]webAppBuilderFactory{
	"http": func(p *Project) boilerplateBuilder {
		return &webAppHttpBuilder{p}
	},
	"chi": func(p *Project) boilerplateBuilder {
		return &webAppChiBuilder{p}
	},
	"echo": func(p *Project) boilerplateBuilder {
		return &webAppEchoBuilder{p}
	},
	"gin": func(p *Project) boilerplateBuilder {
		return &webAppGinBuilder{p}
	},
	"gorilla": func(p *Project) boilerplateBuilder {
		return &webAppGorillaBuilder{p}
	},
	"httprouter": func(p *Project) boilerplateBuilder {
		return &webAppHttprouterBuilder{p}
	},
}
