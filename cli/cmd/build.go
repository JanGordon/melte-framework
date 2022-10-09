package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/JanGordon/melte-framework/compile"
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
	// js, html := compile.RemoveJS("ybox/m.melte")
	// fmt.Println(js)
	// write1, _ := os.Create("ybox/temp.js")
	// write1.Write([]byte(js))
	// write2, _ := os.Create("ybox/temp.html")
	// write2.Write([]byte(html))
	// compile.BuildFile("app.js")
	// os.Remove("temp.js")

	compile.BuildFile()

}
