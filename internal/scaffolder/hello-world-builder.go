package scaffolder

func buildHelloWorld(g *generator) {
	if g.p.Database == "" {
		g.file("hello-world-main")
		g.file("hello-world-main-test")
		return
	}
	g.readme = append(g.readme, "hello-world-readme-with-db")
	buildDatabase(g, g.p.Database)
	g.file("hello-world-main-with-db")
}
