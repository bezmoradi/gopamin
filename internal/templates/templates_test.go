package templates

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
	"text/template"
)

// sampleProject mirrors the fields scaffolder.Project exposes to templates. It
// is declared locally to avoid an import cycle with the scaffolder package.
type sampleProject struct {
	Name        string
	Platform    string
	Broker      string
	Database    string
	ProjectType string
	Logger      string
	Path        string
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
