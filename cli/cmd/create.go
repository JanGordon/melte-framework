package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "creates a new melte project",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Created new project at: ")
	},
}

func init() {
	// the commadn should build scaffold and cd into the project
	rootCmd.AddCommand(createCmd)
}
