package scaffolder

import (
	"bytes"
	"fmt"
	"go/format"
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
	alphanumericName := DirectoryName(name)
	moduleName := ModulePath(name)
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

	// Random passwords for the compose database and broker, so no project ships
	// a known credential on a published port.
	p.Secrets.DatabasePassword = randomHex(16)
	p.Secrets.DatabaseRootPassword = randomHex(16)
	p.Secrets.BrokerPassword = randomHex(16)

	// -a jwt: a random signing secret and demo password, so no two projects share
	// one and none ships a known default. Only the password's hash is written.
	var demoPassword string
	if p.HasAuth() {
		demoPassword = randomToken(16)
		hash, err := hashPassword(demoPassword)
		if err != nil {
			fmt.Printf("Error generating project: %v\n", err)
			_ = os.RemoveAll(projectPath)
			os.Exit(1)
		}
		p.Secrets.JWTSecret = randomToken(32)
		p.Secrets.AuthPasswordHash = hash
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

	fmt.Printf("%v "+buildSuccessMessage+"\n", p.Name)
	if demoPassword != "" {
		color.Cyan("Demo login for POST /login: username \"admin\", password %q.", demoPassword)
		color.Cyan("This is the only time the password is shown; .env keeps only its hash. Run `make hash-password` to set a new one.")
	}
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

// IsProjectNameTaken reports whether the project's directory already exists, as
// anything, even an empty directory: New creates it with Mkdir, which would fail.
func IsProjectNameTaken(name string) bool {
	_, err := os.Lstat(DirectoryName(name))
	return err == nil
}

// IsValidProjectName reports whether name yields a usable directory and module
// path: one that starts with a letter or digit, so neither `go mod init` nor
// another command reads it as a flag.
func IsValidProjectName(name string) bool {
	module := ModulePath(name)
	if module == "" {
		return false
	}
	first := module[0]
	return ('a' <= first && first <= 'z') || ('A' <= first && first <= 'Z') || ('0' <= first && first <= '9')
}

// ModulePath is the go.mod module path for name: whitespace and characters a
// module path can't hold become '-'. Case is kept, so github.com/Acme/svc stays
// as written.
func ModulePath(name string) string {
	return replaceNonAlphanumeric(name, "/")
}

// DirectoryName is the directory created for name: the module path with '/'
// replaced too, lowercased.
func DirectoryName(name string) string {
	return strings.ToLower(replaceNonAlphanumeric(name))
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
	return nonAlphanumericRegex.ReplaceAllString(strings.TrimSpace(input), "-")
}
