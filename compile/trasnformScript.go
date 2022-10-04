package compile

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

func TransformScript(script string) string {
	result := api.Transform(script, api.TransformOptions{
		Loader: api.LoaderTS,
	})

	return string(result.Code)
}

var componentPlugin = api.Plugin{
	Name: "component",
	Setup: func(build api.PluginBuild) {
		build.OnResolve(api.OnResolveOptions{Filter: `^Component/`},
			func(args api.OnResolveArgs) (api.OnResolveResult, error) {
				parts := strings.Split(args.Path, string(os.PathSeparator))
				endPath := ""
				for part := range parts {
					if part != 0 {
						endPath += parts[part]
					}
				}

				return api.OnResolveResult{
					Path: filepath.Join(args.ResolveDir, "components", endPath),
				}, nil
			})
	},
}

func BuildScriptFile(script string, outDir string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()
	os.WriteFile("in.ts", []byte(script), 0644)
	result := api.Build(api.BuildOptions{
		EntryPoints: []string{"in.ts"},
		Outfile:     outDir,
		Write:       true,
		Bundle:      true,
		Plugins:     []api.Plugin{componentPlugin},
		Loader: map[string]api.Loader{
			".png": api.LoaderDataURL,
			".js":  api.LoaderJS,
		},
		// Engines: []api.Engine{
		// 	{api.EngineChrome, "64"},
		// 	{api.EngineFirefox, "80"},
		// 	{api.EngineSafari, "11"},
		// 	{api.EngineEdge, "16"},
		// },
		// MinifyWhitespace: true,

		Platform: api.PlatformNeutral,
	})
	// os.Remove("in.ts")
	for _, err := range result.Errors {
		fmt.Printf("build error: %s found at %d:%d in file %s: \n %s\n suggested change: %s", err.Text, err.Location.Line, err.Location.Column, err.Location.File, err.Location.LineText, err.Location.Suggestion)
	}
	// scriptCode, err := os.ReadFile("out.js")
	// if err != nil {
	// 	panic("File removed during compilation")
	// }
	// regex := regexp.MustCompile("//require(//")
	// r2 := regex.ReplaceAllString("hello my name require(hello world)", "import ")
	// fmt.Println(r2)
	// // This should replace cjs require wiht esm imports

	// fmt.Println(result)
}
