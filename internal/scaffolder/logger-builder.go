package scaffolder

type logBuilder struct {
	project *Project
}

func (m *logBuilder) build() {
	fileGenerator([]string{"log-logger"}, m.project)
	fileGenerator([]string{"logger-interface"}, m.project)
}

type slogBuilder struct {
	project *Project
}

func (m *slogBuilder) build() {
	fileGenerator([]string{"slog-logger"}, m.project)
	fileGenerator([]string{"logger-interface"}, m.project)
}

type logrusBuilder struct {
	project *Project
}

func (m *logrusBuilder) build() {
	fileGenerator([]string{"logrus-logger"}, m.project)
	fileGenerator([]string{"logger-interface"}, m.project)
	goGetPackages(m.project.Path, []string{"github.com/sirupsen/logrus"})
}

type zapBuilder struct {
	project *Project
}

func (m *zapBuilder) build() {
	fileGenerator([]string{"zap-logger"}, m.project)
	fileGenerator([]string{"logger-interface"}, m.project)
	goGetPackages(m.project.Path, []string{"go.uber.org/zap"})
}

type loggerBuilderFactory func(p *Project) boilerplateBuilder

var loggerBuilderMap = map[string]loggerBuilderFactory{
	"log": func(p *Project) boilerplateBuilder {
		return &logBuilder{p}
	},
	"slog": func(p *Project) boilerplateBuilder {
		return &slogBuilder{p}
	},
	"logrus": func(p *Project) boilerplateBuilder {
		return &logrusBuilder{p}
	},
	"zap": func(p *Project) boilerplateBuilder {
		return &zapBuilder{p}
	},
}
