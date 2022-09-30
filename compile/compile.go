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

func BuildFile() {
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
	// routesToMake, err := filepath.Glob("route/*.html")
	// if err != nil {
	// 	panic(err)
	// }
	// mux := http.NewServeMux()
	// rh := http.RedirectHandler("http://example.org", 307)
	// for _, file := range routesToMake {
	// 	fmt.Println(strings.TrimSuffix(file, filepath.Ext(file)))
	// 	mux.Handle("/"+strings.TrimSuffix(file, filepath.Ext(file)), rh)
	// }
	// http.ListenAndServe(":3000", mux)
	html := ReplaceComponentWithHTML(ParseHTMLFragmentFromPath("test.html"))
	// Loop over every script in:
	BuildPage(html, "out.html", "./", false, false, true)

	// OutScript := `const SELF = document.querySelector("[melte-id='']")`
	// scripts, html := RemoveJS(filePath)
	// for scriptIndex := range scripts {
	// 	script := TransformScript(scripts[scriptIndex])
	// 	OutScript += script
	// }
}

func DevServer() {

}
