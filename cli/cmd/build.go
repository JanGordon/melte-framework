package cmd

import (
	"fmt"

	"github.com/JanGordon/melte-framework/compile"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"b"},
	Short:   "Builds app to executable to ./build",
	Run: func(cmd *cobra.Command, args []string) {
		Build()
		fmt.Println("Building")
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

func Build() {
	compile.Build()
}
