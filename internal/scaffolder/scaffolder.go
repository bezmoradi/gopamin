package scaffolder

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/fatih/color"
	"github.com/bezmoradi/gopamin/internal/templates"
)

type Project struct {
	Database    string
	Platform    string
	Name        string
	Path        string
	ProjectType string
	Logger      string
}

func New(projectType, platform, name, database, logger string) {
	alphanumericName := replaceNonAlphanumeric(name)
	moduleName := replaceNonAlphanumeric(name, "/")
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Unable to get your current working directory")
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
		Name:        moduleName,
		Path:        projectPath,
		ProjectType: projectType,
		Logger:      logger,
	}

	generateProjectAgnosticFiles(&p)

	builderFactory := buildersMap[p.ProjectType]
	builder := builderFactory(&p)
	director := &director{builder: builder}
	director.construct()
}

func IsProjectNameTaken(name string) bool {
	name = replaceNonAlphanumeric(name)

	if _, err := os.Stat(name); err == nil {
		dirEntries, err := os.ReadDir(name)
		if err != nil {
			log.Println("Could not read the directory")
		}
		if len(dirEntries) > 0 {
			return true
		}
	}

	return false
}

func generateProjectAgnosticFiles(p *Project) {
	fileGenerator([]string{"gitignore"}, p)
	fileGenerator([]string{"license"}, p)

	initGit(p.Path)
	initGoMod(p.Name, p.Path)
	goGetPackages(p.Path, []string{"github.com/joho/godotenv"})
}

func fileGenerator(fileTypes []string, p *Project) {
	templateMapper := templates.Mapper()
	var concatenatedContent string
	var fileName string

	for _, fileType := range fileTypes {
		fileTemplate, name := templateMapper[fileType]()
		fileName = name
		concatenatedContent += string(fileTemplate) + "\n\n"
	}

	dir := filepath.Dir(filepath.Join(p.Path, fileName))
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Printf("Error creating directory: %s\n", err)

			os.Exit(1)
		}
	}

	file, err := os.Create(filepath.Join(p.Path, fileName))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	template := template.Must(
		template.New(fileName).Parse(
			string(concatenatedContent),
		),
	)

	err = template.Execute(file, p)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	color.Green(fileName + " created")
}

func executeCommand(name string, args []string, dir string) error {
	var out bytes.Buffer

	command := exec.Command(name, args...)
	command.Dir = dir
	command.Stdout = &out

	if err := command.Run(); err != nil {
		return err
	}

	return nil
}

func initGoMod(projectName string, appDir string) error {
	if err := executeCommand("go", []string{"mod", "init", projectName}, appDir); err != nil {
		return err
	}

	return nil
}

func initGit(appDir string) {
	err := executeCommand("git", []string{"init"}, appDir)
	if err != nil {
		fmt.Printf("Error initializing git repo: %v", err)
		os.Exit(1)
	}
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
