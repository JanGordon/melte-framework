package compile

import (
	"os"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func BuildPage(root []*html.Node, filepath string) {
	// This function should build a full html page from the list of Scripts and the component
	writeFile, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0600)
	f, err := os.ReadFile(filepath)
	newPage, err := html.Parse(strings.NewReader(string(f)))
	if err != nil {
		panic(err)
	}
	for script := range Scripts {
		scriptC := &html.Node{
			Data:     "script",
			Type:     html.ElementNode,
			DataAtom: atom.Script,
		}
		scriptC.Attr = append(scriptC.Attr, html.Attribute{
			Key: "type",
			Val: "module",
		})

		newPage.AppendChild(scriptC)
		scriptC.AppendChild(&Scripts[script])
	}
	if err = html.Render(writeFile, newPage); err != nil {
		panic(err)
	}

}
