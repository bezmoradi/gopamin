package scaffolder

import "fmt"

func buildHelloWorld(p *Project) {
	readme := []string{"readme", p.Logger + "-readme"}
	makefile := []string{"makefile"}
	env := []string{"env"}

	if p.Database == "" {
		fileGenerator([]string{"hello-world-main"}, p)
		fileGenerator([]string{"hello-world-main-test"}, p)
	} else {
		fileGenerator([]string{"hello-world-main-with-db"}, p)
		readme = append(readme, p.Database+"-readme", "hello-world-readme-with-db")
		env = append(env, p.Database+"-env")
		buildDatabase(p, p.Database)
	}

	if hasContainerizedBackend(p) {
		makefile = append(makefile, "docker-makefile")
	}

	buildLogger(p)
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

	fmt.Printf("%v "+BUILD_SUCCESS_MESSAGE+"\n", p.Name)
}
