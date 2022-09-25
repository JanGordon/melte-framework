package compile

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func ReplaceComponentWithHTML(filepath string) []*html.Node {
	f, err := os.ReadFile(filepath)
	//writeFile, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	root, err := html.ParseFragment(strings.NewReader(string(f)), &html.Node{
		Type:     html.ElementNode,
		Data:     "div",
		DataAtom: atom.Div,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting Build")
	for child := range root {
		replace(root[child])

		// if err = html.Render(writeFile, root[child]); err != nil {
		// 	panic(err)
		// }
	}
	fmt.Println(root)
	return root
}

var Scripts []html.Node

func replace(n *html.Node) {
	if n.Type == html.ElementNode {
		f, err := os.ReadFile(n.Data + ".melte")
		fmt.Println(n.Data)
		if err == nil {
			// ComponentName := n.Data + "awdw"
			// OutScript := ""
			// scripts, _ := RemoveJS(n.Data + ".melte")
			// for scriptIndex := range scripts {
			// 	script := TransformScript(scripts[scriptIndex])
			// 	OutScript += script
			// }
			// fmt.Println(OutScript)
			ReplaceComponentWithHTML(n.Data + ".melte")
			component, err := html.ParseFragment(strings.NewReader(string(f)), &html.Node{
				Type:     html.ElementNode,
				Data:     "div",
				DataAtom: atom.Div,
			})
			fmt.Println("Inserting Component...", component)
			if err != nil {
				panic(err)
			}
			n.Attr = append(n.Attr, html.Attribute{
				Key: "melte-id",
				Val: n.Data + "-awdw",
			})

			for node := range component {
				if component[node].Data == "script" {
					OutScript := fmt.Sprintf(`const SELF = document.querySelector("[melte-id='%s']")`, n.Data+"-awdw")
					// We need to move the script to end and add module tag

					scriptComponent := &html.Node{
						Data: OutScript + component[node].FirstChild.Data,
						Type: html.TextNode,
					}

					component[node].RemoveChild(component[node].FirstChild)
					// component[node].AppendChild(scriptComponent)
					fmt.Println("Inserting Component...", component)
					if err != nil {
						panic(err)
					}
					fmt.Println("Found Script")
					Scripts = append(Scripts, *scriptComponent)

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
