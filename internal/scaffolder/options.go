package scaffolder

import (
	"maps"
	"slices"
)

// The values each flag accepts, read from the recipe tables, so adding a recipe
// is all it takes for the CLI to accept the new value.

func ProjectTypes() []string { return sortedKeys(buildersMap) }

// Frameworks lists the -f values for projectType (none for hello-world/worker).
func Frameworks(projectType string) []string {
	switch projectType {
	case "api":
		return sortedKeys(apiRecipes)
	case "web-app":
		return sortedKeys(webRecipes)
	}
	return nil
}

func Brokers() []string { return sortedKeys(workerRecipes) }

// Databases lists the -d values; "memory" is the store api/web-app use without
// -d, not something to ask for.
func Databases() []string {
	return slices.DeleteFunc(sortedKeys(databaseRecipes), func(d string) bool { return d == "memory" })
}

func Loggers() []string { return sortedKeys(loggerRecipes) }

func CIs() []string { return sortedKeys(ciRecipes) }

func sortedKeys[V any](m map[string]V) []string {
	return slices.Sorted(maps.Keys(m))
}
