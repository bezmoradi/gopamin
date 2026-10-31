package scaffolder

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/bezmoradi/gopamin/internal/templates"
	"github.com/fatih/color"
)

func New(projectType, platform, broker, name, database, logger, ci, observability, auth, openapi string) {
	alphanumericName := replaceNonAlphanumeric(name)
	moduleName := replaceNonAlphanumeric(name, "/")
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Unable to get your current working directory:", err)
		os.Exit(1)
	}

	projectPath := filepath.Join(currentDir, alphanumericName)

	err = os.Mkdir(alphanumericName, 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	p := Project{
		Database:      database,
		Platform:      platform,
		Broker:        broker,
		Name:          moduleName,
		Path:          projectPath,
		ProjectType:   projectType,
		Logger:        logger,
		GoVersion:     generatedGoVersion,
		CI:            ci,
		Observability: observability,
		Auth:          auth,
		OpenAPI:       openapi,
	}

	if err := generate(&p); err != nil {
		fmt.Printf("Error generating project: %v\n", err)
		_ = os.RemoveAll(projectPath)
		os.Exit(1)
	}

	// `go get` records every dependency as `// indirect`; run `go mod tidy`
	// once, after all files and go-gets, so directly-imported modules are promoted
	// to direct requirements and anything unused is dropped. Non-fatal: the project
	// is already complete and its go.mod is valid, so a tidy hiccup (e.g. a transient
	// network failure) should surface a warning, not discard the generated project.
	if err := goModTidy(p.Path); err != nil {
		color.Yellow("warning: go mod tidy failed: %v", err)
	}

	fmt.Printf("%v "+BUILD_SUCCESS_MESSAGE+"\n", p.Name)
}

// generate writes every file for p, then initializes git and the Go module and
// fetches the modules the files need, in one pinned `go get`.
func generate(p *Project) error {
	g := &generator{p: p}
	g.need("github.com/joho/godotenv")
	buildersMap[p.ProjectType](g)
	buildCommon(g)
	if g.err != nil {
		return g.err
	}

	if err := initGit(p.Path); err != nil {
		return err
	}
	if err := initGoMod(p.Name, p.Path); err != nil {
		return err
	}
	return goGetPackages(p.Path, append(g.modules, securityFloor...))
}

func IsProjectNameTaken(name string) bool {
	name = replaceNonAlphanumeric(name)

	if _, err := os.Stat(name); err == nil {
		dirEntries, err := os.ReadDir(name)
		if err != nil {
			log.Println("Could not read the directory")
			return true
		}
		if len(dirEntries) > 0 {
			return true
		}
	}

	return false
}

// fileGenerator renders the templates named by keys, concatenated in order, into
// the output path of the last key (concatenated templates share one path), and
// gofmt's Go output.
func fileGenerator(keys []string, p *Project) error {
	var concatenatedContent strings.Builder
	var fileName string

	for _, key := range keys {
		content, name, err := templates.Lookup(key)
		if err != nil {
			return err
		}
		fileName = name
		concatenatedContent.Write(content)
		concatenatedContent.WriteString("\n\n")
	}

	if err := os.MkdirAll(filepath.Dir(filepath.Join(p.Path, fileName)), 0755); err != nil {
		return err
	}

	tmpl, err := template.New(fileName).Parse(concatenatedContent.String())
	if err != nil {
		return fmt.Errorf("%s: %w", fileName, err)
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, p); err != nil {
		return fmt.Errorf("%s: %w", fileName, err)
	}

	content := rendered.Bytes()
	if strings.HasSuffix(fileName, ".go") {
		formatted, err := format.Source(content)
		if err != nil {
			// Don't fail generation on a formatting error; surface it so the
			// offending template can be fixed, and write the file unformatted.
			fmt.Printf("warning: gofmt failed for %s: %v\n", fileName, err)
		} else {
			content = formatted
		}
	}

	if err := os.WriteFile(filepath.Join(p.Path, fileName), content, 0644); err != nil {
		return err
	}

	color.Green(fileName + " created")
	return nil
}

func executeCommand(name string, args []string, dir string) error {
	var stdout, stderr bytes.Buffer

	command := exec.Command(name, args...)
	command.Dir = dir
	command.Stdout = &stdout
	command.Stderr = &stderr

	if err := command.Run(); err != nil {
		if detail := strings.TrimSpace(stderr.String()); detail != "" {
			return fmt.Errorf("%s %s: %w\n%s", name, strings.Join(args, " "), err, detail)
		}
		return fmt.Errorf("%s %s: %w", name, strings.Join(args, " "), err)
	}

	return nil
}

// generatedGoVersion is the `go` directive stamped into every project Gopamin
// scaffolds, so generated projects target a consistent Go version regardless of
// the user's local toolchain. Keep it in lockstep with the `go` directive in
// this repo's go.mod (and the README) whenever the supported Go version is bumped.
const generatedGoVersion = "1.27.3"

func initGoMod(projectName string, appDir string) error {
	if err := executeCommand("go", []string{"mod", "init", projectName}, appDir); err != nil {
		return err
	}

	if err := executeCommand("go", []string{"mod", "edit", "-go=" + generatedGoVersion}, appDir); err != nil {
		return err
	}

	return nil
}

func initGit(appDir string) error {
	return executeCommand("git", []string{"init"}, appDir)
}

func goModTidy(appDir string) error {
	return executeCommand("go", []string{"mod", "tidy"}, appDir)
}

// goGetPackages adds modules to the generated go.mod at their pinned versions
// (see modulePins), in one `go get` so the module graph is resolved once.
func goGetPackages(appDir string, modules []string) error {
	if len(modules) == 0 {
		return nil
	}
	args, err := pinned(modules)
	if err != nil {
		return err
	}
	return executeCommand("go", append([]string{"get"}, args...), appDir)
}

func replaceNonAlphanumeric(input string, exclude ...string) string {
	var nonAlphanumericRegex *regexp.Regexp
	if len(exclude) > 0 {
		nonAlphanumericRegex = regexp.MustCompile("[^a-zA-Z0-9._" + exclude[0] + "]")
	} else {
		nonAlphanumericRegex = regexp.MustCompile("[^a-zA-Z0-9._]")
	}
	return nonAlphanumericRegex.ReplaceAllString(strings.TrimSpace(strings.ToLower(input)), "-")
}
