package compile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func BuildPage(root []*html.Node, outPath string, outPathJS string, inlineJS bool, dev bool, findLayouts bool) {
	// This function should build a full html page from the list of Scripts and the component
	//fmt.Println("Building the page: out.html and all scripts")
	writeFile, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE, 0600)
	writeFile.Write([]byte(""))
	//newPage, err := html.Parse(strings.NewReader(string(f)))
	if err != nil {
		panic(err)
	}
	importLines := ""
	scriptExceptImports := ""
	if findLayouts {
		root = populateLayout(root, outPath, writeFile)

	}
	for script := range Scripts {

		if inlineJS {
			scriptC := &html.Node{
				Data:     "script",
				Type:     html.ElementNode,
				DataAtom: atom.Script,
			}
			Scripts[script].Data = TransformScript(Scripts[script].Data)
			scriptC.AppendChild(&Scripts[script])
			root = append(root, scriptC)
			//fmt.Println("Adding Script", Scripts)
		} else {
			importRemovedLines := ""
			lines := strings.Split(Scripts[script].Data, "\n")
			for _, line := range lines {
				if strings.HasPrefix(strings.TrimSpace(line), "import") {
					importLines += fmt.Sprintf("%s\n", line)
				} else {
					importRemovedLines += line + "\n"
				}
			}
			scriptExceptImports = fmt.Sprintf("{\n// script for %s\n %s}", fmt.Sprintf("out%s.js", ScriptIDs[script]), importRemovedLines)
		}

		// newPage.AppendChild(scriptC)

	}
	file := importLines + "\n" + scriptExceptImports
	BuildScriptFile(file, filepath.Join(outPathJS, "out.js"))

	scriptC := &html.Node{
		Data:     "script",
		Type:     html.ElementNode,
		DataAtom: atom.Script,
	}
	scriptC.Attr = append(scriptC.Attr, html.Attribute{
		Key: "src",
		Val: filepath.Join(strings.Replace(outPathJS, "routes", "", 1), "out.js"),
	})
	root = append(root, scriptC)
	if dev {
		scriptDev := &html.Node{
			Data:     "script",
			Type:     html.ElementNode,
			DataAtom: atom.Script,
		}

		scriptDev.Attr = append(scriptDev.Attr, html.Attribute{
			Key: "src",
			Val: "/hotReload/WebSocket.js",
		})
		root = append(root, scriptDev)

	}

	//fmt.Println("Adding Script\n", file, root)
	for child := range root {
		if err = html.Render(writeFile, root[child]); err != nil {
			panic(err)
		}
	}
	Scripts = nil
	ScriptIDs = nil
	CCount = 0
	//fmt.Println(len(Scripts))
	//f, err := os.ReadFile(filepath.Join(outPathJS, "out.js"))
	//fmt.Println(string(f))
}
