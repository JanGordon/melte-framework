package compile

import (
	"fmt"

	"github.com/evanw/esbuild/pkg/api"
)

func main() {
	result := api.Transform("let x: number = 1", api.TransformOptions{
		Loader: api.LoaderTS,
	})

	if len(result.Errors) == 0 {
		fmt.Printf("%s", result.Code)
	}
}

func BuildFile(filePath string, outPath string) {
	// result := api.Build(api.BuildOptions{
	// 	EntryPoints:       []string{filePath},
	// 	Bundle:            true,
	// 	MinifyWhitespace:  true,
	// 	MinifyIdentifiers: true,
	// 	MinifySyntax:      true,
	// 	Loader: map[string]api.Loader{
	// 		".html": api.LoaderFile,
	// 		".svg":  api.LoaderText,
	// 	},
	// 	Write:  true,
	// 	Outdir: "out/ybox",
	// })
	// fmt.Println(result)

	// if len(result.Errors) > 0 {
	// 	os.Exit(1)
	// }

	html := ReplaceComponentWithHTML(filePath)
	// Loop over every script in:
	BuildPage(html, outPath, false)

	// OutScript := `const SELF = document.querySelector("[melte-id='']")`
	// scripts, html := RemoveJS(filePath)
	// for scriptIndex := range scripts {
	// 	script := TransformScript(scripts[scriptIndex])
	// 	OutScript += script
	// }
}
