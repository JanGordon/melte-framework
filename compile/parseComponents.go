package compile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var CCount = 0

func ReplaceComponentWithHTML(root []*html.Node) []*html.Node {
	CCount++
	for child := range root {
		replace(root[child])

		// if err = html.Render(writeFile, root[child]); err != nil {
		// 	panic(err)
		// }
	}
	CCount = 0
	return root
}

var Scripts []html.Node
var ScriptIDs []string

func replace(n *html.Node) {
	if n.Type == html.ElementNode {
		wd, err := os.Getwd()
		if err != nil {
			panic("Failed to get working directory")
		}
		f, err := os.ReadFile(filepath.Join(wd, "components", n.Data+".melte"))
		// seeing if custom component exists
		if err == nil {
			// ComponentName := n.Data + "awdw"
			// OutScript := ""
			// scripts, _ := RemoveJS(n.Data + ".melte")
			// for scriptIndex := range scripts {
			// 	script := TransformScript(scripts[scriptIndex])
			// 	OutScript += script
			// }
			// fmt.Println(OutScript)
			//fi, err := os.ReadFile(filepath.Join(wd, "components", n.Data+".melte"))
			//writeFile, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				panic(err)
			}
			//component := ReplaceComponentWithHTML(f) // adds components scripts to Scripts
			//fmt.Println("Replacing component with : ", component[0].Data)
			component, err := html.ParseFragment(strings.NewReader(string(f)), &html.Node{
				Type:     html.ElementNode,
				Data:     "div",
				DataAtom: atom.Div,
			})
			fmt.Println("Inserting Component...")
			if err != nil {
				panic(err)
			}
			n.Attr = append(n.Attr, html.Attribute{
				Key: "melte-id",
				Val: n.Data + fmt.Sprintf("%d", CCount),
			})
			for node := range component {
				if component[node].Data == "script" {
					OutScript := fmt.Sprintf(`const SELF = document.querySelector("[melte-id='%s']")`, n.Data+fmt.Sprintf("%d", CCount))
					// We need to move the script to end and add module tag

					scriptComponent := &html.Node{
						Data: OutScript + component[node].FirstChild.Data,
						Type: html.TextNode,
					}

					component[node].RemoveChild(component[node].FirstChild)
					// component[node].AppendChild(scriptComponent)
					if err != nil {
						panic(err)
					}
					Scripts = append(Scripts, *scriptComponent)
					ScriptIDs = append(ScriptIDs, fmt.Sprintf("out-%s%d.js", n.Data, CCount))

				} else {
					n.AppendChild(component[node])
				}

			}

			// node := html.Node()
		}
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		replace(child)
	}
}

func ParseHTMLFragmentFromPath(path string) []*html.Node {
	//do what old replacehtml did
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	root, err := html.ParseFragment(strings.NewReader(string(file)), &html.Node{
		Type:     html.ElementNode,
		Data:     "div",
		DataAtom: atom.Div,
	})
	if err != nil {
		panic(err)
	}
	return ReplaceComponentWithHTML(root)
}
