package templates

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
	"text/template"

	"gopkg.in/yaml.v3"
)

// sampleProject mirrors the fields scaffolder.Project exposes to templates. It
// is declared locally to avoid an import cycle with the scaffolder package.
type sampleProject struct {
	Name          string
	Platform      string
	Broker        string
	Database      string
	ProjectType   string
	Logger        string
	Path          string
	GoVersion     string
	CI            string
	Observability string
	Auth          string
}

// TestAllTemplatesRenderAndParse is a hermetic smoke test (no network, no shell
// out): every template in the registry must parse and execute, and every
// template that renders to a .go file must produce syntactically valid Go. This
// catches template-directive mistakes and Go syntax errors before release.
func TestAllTemplatesRenderAndParse(t *testing.T) {
	// Every field is a valid Go identifier / import path so that any template —
	// whatever placeholders it substitutes — renders to syntactically valid Go.
	// (The real generator only renders a template with its matching project
	// shape; here we just need parse-validity, not semantic accuracy.)
	p := sampleProject{
		Name:        "example.com/app",
		Platform:    "echo",
		Broker:      "kafka",
		Database:    "mysql",
		ProjectType: "api",
		Logger:      "zap",
		Path:        "/tmp/app",
		GoVersion:   "1.26.0",
	}

	for key, gen := range Mapper() {
		content, filename := gen()
		tmpl, err := template.New(filename).Parse(string(content))
		if err != nil {
			t.Errorf("template %q (%s): parse error: %v", key, filename, err)
			continue
		}

		var sb strings.Builder
		if err := tmpl.Execute(&sb, p); err != nil {
			t.Errorf("template %q (%s): execute error: %v", key, filename, err)
			continue
		}
		if strings.HasSuffix(filename, ".go") {
			if _, err := parser.ParseFile(token.NewFileSet(), filename, sb.String(), parser.AllErrors); err != nil {
				t.Errorf("template %q (%s): rendered Go does not parse: %v", key, filename, err)
			}
		}
	}
}

// TestAgentsTemplateRendersPerType renders AGENTS.md across every project shape
// (its content is heavily conditional on type/platform/broker/database) and asserts
// each renders without error to non-empty Markdown whose type-specific section is
// present. Guards the many {{if}} branches the single-shape smoke test can't reach.
func TestAgentsTemplateRendersPerType(t *testing.T) {
	content, filename := Mapper()["agents"]()
	tmpl, err := template.New(filename).Parse(string(content))
	if err != nil {
		t.Fatalf("agents template parse error: %v", err)
	}

	cases := []struct {
		p    sampleProject
		want string // a type-specific heading that must appear
	}{
		{sampleProject{ProjectType: "api", Platform: "echo", Database: "postgres", Broker: "kafka", Logger: "slog", Name: "x"}, "## HTTP API"},
		{sampleProject{ProjectType: "api", Platform: "graphql", Database: "mysql", Logger: "zap", Name: "x"}, "## GraphQL API"},
		{sampleProject{ProjectType: "web-app", Platform: "gin", Database: "mongodb", Logger: "logrus", Name: "x"}, "## Web UI"},
		{sampleProject{ProjectType: "worker", Broker: "redis", Logger: "slog", Name: "x"}, "## Broker workflow"},
		{sampleProject{ProjectType: "worker", Broker: "rabbitmq", Database: "postgres", Logger: "log", Name: "x"}, "## Broker workflow"},
		{sampleProject{ProjectType: "hello-world", Logger: "log", Name: "x"}, "## What this project is"},
		{sampleProject{ProjectType: "hello-world", Database: "sqlite", Logger: "slog", Name: "x"}, "## Running"},
	}

	for _, c := range cases {
		name := c.p.ProjectType + "/" + c.p.Platform + "/" + c.p.Broker + "/" + c.p.Database
		var sb strings.Builder
		if err := tmpl.Execute(&sb, c.p); err != nil {
			t.Errorf("%s: execute error: %v", name, err)
			continue
		}
		out := sb.String()
		if len(strings.TrimSpace(out)) == 0 {
			t.Errorf("%s: rendered AGENTS.md is empty", name)
		}
		if !strings.Contains(out, c.want) {
			t.Errorf("%s: rendered AGENTS.md missing expected section %q", name, c.want)
		}
	}
}

// archConfig is the subset of the go-arch-lint schema the smoke test validates.
type archConfig struct {
	Version    int                       `yaml:"version"`
	Components map[string]map[string]any `yaml:"components"`
	Deps       map[string]struct {
		MayDependOn []string `yaml:"mayDependOn"`
	} `yaml:"deps"`
}

// TestArchLintTemplateRendersValidConfig renders the .go-arch-lint.yml template
// across every project shape that emits it (each exercises a different set of the
// component-presence conditionals) and asserts, hermetically, that the result is
// valid YAML whose dependency rules never reference a component that wasn't
// declared. A dangling reference would make go-arch-lint reject the config, and an
// undeclared (empty-glob) component would make it exit 1 on a clean project — both
// invisible to a Go-syntax check, so they are asserted here.
func TestArchLintTemplateRendersValidConfig(t *testing.T) {
	content, filename := Mapper()["arch-lint"]()
	tmpl, err := template.New(filename).Parse(string(content))
	if err != nil {
		t.Fatalf("arch-lint template parse error: %v", err)
	}

	shapes := []sampleProject{
		{ProjectType: "api", Platform: "echo", Broker: "kafka", Database: "postgres"}, // all components
		{ProjectType: "api", Platform: "graphql", Database: "mysql"},                  // no api-response, no brokers
		{ProjectType: "api", Platform: "chi"},                                         // no db (mock), no brokers
		{ProjectType: "web-app", Platform: "gin", Database: "mongodb"},                // web, no api-response
		{ProjectType: "web-app", Platform: "echo", Broker: "redis"},                   // web + brokers, no db
		{ProjectType: "worker", Broker: "rabbitmq", Database: "postgres"},             // core + brokers, no handlers
		{ProjectType: "worker", Broker: "kafka"},                                      // no services/repositories
		{ProjectType: "hello-world", Database: "sqlite"},                              // db-only core
		{ProjectType: "hello-world"},                                                  // plain: no domain/services
		{ProjectType: "api", Platform: "echo", Database: "postgres", Observability: "otel"}, // observability: component + cmd/cmd-server edges
		{ProjectType: "worker", Broker: "kafka", Observability: "otel"},                     // observability on a worker (cmd edge, no cmd-server)
	}

	// Components present in every shape that emits the config. (The `shared` package —
	// logger/router — is always present; the `user` aggregate is absent from a plain
	// hello-world, so its components are not in this set.)
	always := []string{"shared", "loggers", "configs", "tools", "cmd"}

	for _, p := range shapes {
		name := p.ProjectType + "/" + p.Platform + "/" + p.Broker + "/" + p.Database
		var sb strings.Builder
		if err := tmpl.Execute(&sb, p); err != nil {
			t.Errorf("%s: execute error: %v", name, err)
			continue
		}

		var cfg archConfig
		if err := yaml.Unmarshal([]byte(sb.String()), &cfg); err != nil {
			t.Errorf("%s: rendered config is not valid YAML: %v\n%s", name, err, sb.String())
			continue
		}
		if cfg.Version != 3 {
			t.Errorf("%s: expected version 3, got %d", name, cfg.Version)
		}
		if len(cfg.Components) == 0 {
			t.Errorf("%s: no components declared", name)
		}
		for _, c := range always {
			if _, ok := cfg.Components[c]; !ok {
				t.Errorf("%s: expected component %q to always be present", name, c)
			}
		}
		// Every dependency rule — and every component it points at — must refer to
		// a declared component (no dangling reference from a conditional branch).
		for comp, rule := range cfg.Deps {
			if _, ok := cfg.Components[comp]; !ok {
				t.Errorf("%s: deps entry %q has no matching component", name, comp)
			}
			for _, dep := range rule.MayDependOn {
				if _, ok := cfg.Components[dep]; !ok {
					t.Errorf("%s: component %q mayDependOn undeclared component %q", name, comp, dep)
				}
			}
		}
	}
}
