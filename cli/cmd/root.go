package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "melte",
	Short: "melte - a blazingly fast full-stack framework for js, ts, go and html written in go",
	Long: `melte is heavily inspired by svelte (hence the name) and nextjs.
	
melte writes code with you without any of the boilerplate that come with most modern frameworks. It has an intuitive syntax that mimicks that of svelte and an easy to use file based router. It also has a blazingly fast dev server along with one designed for deployment`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing melte '%s'", err)
		os.Exit(1)
	}
}
