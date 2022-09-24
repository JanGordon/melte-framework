package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:     "configure",
	Aliases: []string{"config"},
	Short:   "Displays config",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Melte v0.0.0")
		fmt.Println("Config: ")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
