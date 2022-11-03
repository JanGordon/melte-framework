package cmd

import (
	"fmt"
	"strconv"

	"github.com/JanGordon/melte-framework/dev"
	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:     "development",
	Aliases: []string{"dev"},
	Args:    cobra.ExactArgs(1),
	Short:   "Runs a dev server wiht hot reload",
	Run: func(cmd *cobra.Command, args []string) {
		portNum, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(fmt.Errorf("please enter a valid number for port"))
		}
		dev.Run(portNum)
		fmt.Println("dev server running")
	},
}

func init() {
	rootCmd.AddCommand(devCmd)
}
