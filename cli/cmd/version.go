package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "shows the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Melte v0.0.2")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
