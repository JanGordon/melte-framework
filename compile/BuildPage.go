package compile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func BuildPage(root html.Node, outPath string, outPathJS string, inlineJS bool, dev bool, findLayouts bool) {
	// This function should build a full html page from the list of Scripts and the component
	//fmt.Println("Building the page: out.html and all scripts")
	os.Truncate(outPath, 0)
	writeFile, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE, 0600)

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
			root.AppendChild(scriptC)
			//fmt.Println("Adding Script", Scripts)
		} else {
			importRemovedLines := ""
			lines := strings.Split(Scripts[script].Data, "\n")
			for i, line := range lines {
				if strings.HasPrefix(strings.TrimSpace(line), "import") {
					importLines += fmt.Sprintf("%s\n", line)
				} else if strings.HasPrefix(strings.TrimSpace(line), "//@melte-custom:") {
					args := strings.Split(string(strings.Replace(string(strings.TrimSpace(line)), "//@melte-custom: ", "", 1)), ",")
					// found a custom declaration
					// get type
					if strings.TrimSpace((args[0])) == "var" {
						//varDeclaration := lines[i+1]
						// put in head if should'nt be reloaded on page change
						lines[i+1] = ""
						if stringInSlice("server", args) {

						}
					}
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
	fmt.Println(outPathJS)
	scriptC := &html.Node{
		Data:     "script",
		Type:     html.ElementNode,
		DataAtom: atom.Script,
	}
	scriptC.Attr = append(scriptC.Attr, html.Attribute{
		Key: "src",
		Val: filepath.Join(strings.Replace(outPathJS, "routes", "", 1), "out.js"),
	})
	root.AppendChild(scriptC)
	scriptFlamethrower := &html.Node{
		Data:     "script",
		Type:     html.ElementNode,
		DataAtom: atom.Script,
	}
	scriptFlamethrower.Attr = append(scriptFlamethrower.Attr, html.Attribute{
		Key: "type",
		Val: "module",
	})
	scriptFlamethrower.Attr = append(scriptFlamethrower.Attr, html.Attribute{
		Key: "src",
		Val: "/clientSideRouting/src.js",
	})
	root.AppendChild(scriptFlamethrower)
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
		root.AppendChild(scriptDev)

	}

	//fmt.Println("Adding Script\n", file, root)
	child := root.FirstChild
	lastChild := root.LastChild
	for {
		if err = html.Render(writeFile, child); err != nil {
			panic(err)
		}
		if child != lastChild {
			child = child.NextSibling
		} else {
			break
		}
	}
	Scripts = nil
	ScriptIDs = nil
	CCount = 0
	//fmt.Println(len(Scripts))
	//f, err := os.ReadFile(filepath.Join(outPathJS, "out.js"))
	//fmt.Println(string(f))
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
