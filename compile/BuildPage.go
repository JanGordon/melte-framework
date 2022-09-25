package compile

import (
	"os"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func BuildPage(root []*html.Node, outPath string, inlineJS bool) {
	// This function should build a full html page from the list of Scripts and the component
	writeFile, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE, 0600)
	//f, err := os.ReadFile(filepath)
	//newPage, err := html.Parse(strings.NewReader(string(f)))
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

		if inlineJS {
			Scripts[script].Data = TransformScript(Scripts[script].Data)
			scriptC.AppendChild(&Scripts[script])
		} else {
			BuildScriptFile(Scripts[script].Data, "out.js")
			scriptC.Attr = append(scriptC.Attr, html.Attribute{
				Key: "src",
				Val: "out.js",
			})
		}
		root = append(root, scriptC)
		// newPage.AppendChild(scriptC)

	}
	for child := range root {
		if err = html.Render(writeFile, root[child]); err != nil {
			panic(err)
		}
	}

}
