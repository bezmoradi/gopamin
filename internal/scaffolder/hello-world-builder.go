package scaffolder

import "fmt"

type helloWorldProjectBuilder struct {
	project *Project
}

func (b *helloWorldProjectBuilder) build() {
	readme := []string{"readme", b.project.Logger + "-readme"}
	makefile := []string{"makefile"}
	env := []string{"env"}

	if b.project.Database == "" {
		fileGenerator([]string{"hello-world-main"}, b.project)
		fileGenerator([]string{"hello-world-main-test"}, b.project)
	} else {
		fileGenerator([]string{"hello-world-main-with-db"}, b.project)
		readme = append(readme, b.project.Database+"-readme", "hello-world-readme-with-db")
		makefile = append(makefile, b.project.Database+"-makefile")
		env = append(env, b.project.Database+"-env")
		dbSelector(b.project)
	}

	loggerSelector(b.project)
	fileGenerator(env, b.project)
	fileGenerator(readme, b.project)
	fileGenerator(makefile, b.project)
	fileGenerator([]string{"configs"}, b.project)
	fileGenerator([]string{"configs-test"}, b.project)
	fileGenerator([]string{"tools"}, b.project)
	fileGenerator([]string{"tools-test"}, b.project)

	fmt.Printf("%v "+BUILD_SUCCESS_MESSAGE+"\n", b.project.Name)
}

type helloWorldBuilderFactory func(p *Project) boilerplateBuilder

var helloWorldBuilderMap = map[string]helloWorldBuilderFactory{
	"hello-world": func(p *Project) boilerplateBuilder {
		return &helloWorldProjectBuilder{p}
	},
}
