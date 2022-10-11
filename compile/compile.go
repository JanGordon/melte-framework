package compile

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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

func Build() {
	cwd, err := os.Getwd()
	if err != nil {
		panic("failed to get working directory")
	}
	err = filepath.WalkDir(cwd+"/routes", buildRoute)
	if err != nil {
		panic("error reading routes folder")
	}

}

func buildRoute(path string, di fs.DirEntry, err error) error {
	dir, filename := filepath.Split(path)
	if filepath.Ext(path) == ".html" && filename != "out.html" && !strings.HasPrefix(filename, "layout") {
		BuildPage(ReplaceComponentWithHTML(ParseHTMLFragmentFromPath(path)), dir+"out.html", dir, false, true, false)
	}
	return nil
}
