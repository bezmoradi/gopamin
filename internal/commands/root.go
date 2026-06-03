package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gopamin",
	Short: "Gopamin CLI Tool",
	Long: `
	   ______                            _
	  / ____/___  ____  ____ _____ ___  (_)___
	 / / __/ __ \/ __ \/ __  / __  __ \/ / __ \
	/ /_/ / /_/ / /_/ / /_/ / / / / / / / / / /
	\____/\____/ .___/\__,_/_/ /_/ /_/_/_/ /_/
	          /_/

The CLI tool for scaffolding Go projects`,
}

func init() {
	// mark completion hidden
	completion := completionCommand()
	completion.Hidden = true
	rootCmd.AddCommand(completion)
}

func completionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "completion",
		Short: "Generate the autocompletion script for the specified shell",
	}
}

func Execute() {
	// The update check is advisory only: it never blocks a command, and a slow
	// or unreachable network simply skips it (see versionChecker).
	if isValid, newVersion := versionChecker(); !isValid {
		fmt.Printf(`A newer version of the Gopamin CLI is available (%v); you have %v.

%v

To update, run 'go install github.com/bezmoradi/gopamin@%v'.`+"\n\n", newVersion, VERSION, UPDATE_MESSAGE, newVersion)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
