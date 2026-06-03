package scaffolder

var loggerRecipes = map[string]loggerRecipe{
	"log":    {"log-logger", nil},
	"slog":   {"slog-logger", nil},
	"logrus": {"logrus-logger", []string{"github.com/sirupsen/logrus"}},
	"zap":    {"zap-logger", []string{"go.uber.org/zap"}},
}
