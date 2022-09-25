package compile

import (
	"os"

	"github.com/evanw/esbuild/pkg/api"
)

func TransformScript(script string) string {
	result := api.Transform(script, api.TransformOptions{
		Loader: api.LoaderTS,
	})

	return string(result.Code)
}

func BuildScriptFile(script string, outDir string) {
	os.WriteFile("in.ts", []byte(script), 0644)
	api.Build(api.BuildOptions{
		EntryPoints: []string{"in.ts"},
		Outfile:     outDir,
		Write:       true,
		Bundle:      true,
		Loader: map[string]api.Loader{
			".png": api.LoaderDataURL,
			".js":  api.LoaderJS,
		},
		MinifyWhitespace: true,

		Platform: api.PlatformNeutral,
	})
	os.Remove("in.ts")
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
