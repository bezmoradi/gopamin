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

type Project struct {
	Database    string
	Platform    string
	Broker      string
	Name        string
	Path        string
	ProjectType string
	Logger      string
	// GoVersion is the pinned Go version (== generatedGoVersion) exposed to
	// templates so the Dockerfile's build image never drifts from the go.mod
	// `go` directive.
	GoVersion string
}

func New(projectType, platform, broker, name, database, logger string) {
	alphanumericName := replaceNonAlphanumeric(name)
	moduleName := replaceNonAlphanumeric(name, "/")
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Unable to get your current working directory")
		return
	}

	projectPath := filepath.Join(currentDir, alphanumericName)

	err = os.Mkdir(alphanumericName, 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	p := Project{
		Database:    database,
		Platform:    platform,
		Broker:      broker,
		Name:        moduleName,
		Path:        projectPath,
		ProjectType: projectType,
		Logger:      logger,
		GoVersion:   generatedGoVersion,
	}

	if err := generateProjectAgnosticFiles(&p); err != nil {
		fmt.Printf("Error generating project: %v\n", err)
		_ = os.RemoveAll(projectPath)
		os.Exit(1)
	}

	buildersMap[p.ProjectType](&p)
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

func generateProjectAgnosticFiles(p *Project) error {
	fileGenerator([]string{"gitignore"}, p)
	fileGenerator([]string{"license"}, p)
	fileGenerator([]string{"agents"}, p)

	if err := initGit(p.Path); err != nil {
		return err
	}
	if err := initGoMod(p.Name, p.Path); err != nil {
		return err
	}

	return goGetPackages(p.Path, []string{"github.com/joho/godotenv"})
}

func fileGenerator(fileTypes []string, p *Project) {
	templateMapper := templates.Mapper()
	var concatenatedContent strings.Builder
	var fileName string

	for _, fileType := range fileTypes {
		fileTemplate, name := templateMapper[fileType]()
		fileName = name
		concatenatedContent.WriteString(string(fileTemplate) + "\n\n")
	}

	dir := filepath.Dir(filepath.Join(p.Path, fileName))
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Printf("Error creating directory: %s\n", err)
			_ = os.RemoveAll(p.Path)
			os.Exit(1)
		}
	}

	tmpl := template.Must(
		template.New(fileName).Parse(concatenatedContent.String()),
	)

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, p); err != nil {
		fmt.Println(err)
		_ = os.RemoveAll(p.Path)
		os.Exit(1)
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
		fmt.Println(err)
		_ = os.RemoveAll(p.Path)
		os.Exit(1)
	}

	color.Green(fileName + " created")
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
const generatedGoVersion = "1.26.0"

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

func goGetPackages(appDir string, packages []string) error {
	for _, packageName := range packages {
		if err := executeCommand("go",
			[]string{"get", "-u", packageName},
			appDir); err != nil {
			return err
		}
	}

	return nil
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
