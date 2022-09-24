package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:     "development",
	Aliases: []string{"dev"},
	Short:   "Runs a dev server wiht hot reload",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dev server running")
	},
}

func init() {
	rootCmd.AddCommand(devCmd)
}
