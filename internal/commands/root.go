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
	// Keep cobra's working `completion` command, just out of the help listing.
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
