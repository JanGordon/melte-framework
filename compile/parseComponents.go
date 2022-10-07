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

func ReplaceComponentWithHTML(root html.Node) html.Node {
	CCount++
	child := root.FirstChild
	lastChild := root.LastChild
	for {
		replace(child)
		if child != lastChild {
			child = child.NextSibling
		} else {
			break
		}
		// if err = html.Render(writeFile, root[child]); err != nil {
		// 	panic(err)
		// }
	}
	CCount = 0
	return root
}

func ReplaceCustomComponentWithHTML(root []*html.Node) []*html.Node {
	for _, child := range root {
		replace(child)
		CCount++

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
		_, err = os.ReadFile(filepath.Join(wd, "components", n.Data+".melte"))
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
			component := ReplaceCustomComponentWithHTML(ParseHTMLAsComponent(filepath.Join(wd, "components", n.Data+".melte"))) // adds components scripts to Scripts
			//fmt.Println("Replacing component with : ", component[0].Data)
			// component, err := html.ParseFragment(strings.NewReader(string(f)), &html.Node{
			// 	Type:     html.ElementNode,
			// 	Data:     "div",
			// 	DataAtom: atom.Div,
			// })

			if err != nil {
				panic(err)
			}
			n.Attr = append(n.Attr, html.Attribute{
				Key: "melte-id",
				Val: n.Data + fmt.Sprintf("%d", CCount),
			})
			fmt.Println("Component : ", component)
			for _, child := range component {
				if child.Data == "script" {
					OutScript := fmt.Sprintf(`const SELF = document.querySelector("[melte-id='%s']")`, n.Data+fmt.Sprintf("%d", CCount))
					// We need to move the script to end and add module tag

					scriptComponent := &html.Node{
						Data: OutScript + child.FirstChild.Data,
						Type: html.TextNode,
					}
					child.RemoveChild(child.FirstChild)
					// component[node].AppendChild(scriptComponent)
					if err != nil {
						panic(err)
					}
					Scripts = append(Scripts, *scriptComponent)
					ScriptIDs = append(ScriptIDs, fmt.Sprintf("out-%s%d.js", n.Data, CCount))

				} else {
					// this is happening twice for some reason
					fmt.Println(child.Data, ": adding component node")
					n.AppendChild(child)
				}
				// if err = html.Render(writeFile, root[child]); err != nil {
				// 	panic(err)
				// }
			}

			// node := html.Node()
		}
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		replace(child)
	}
}

func ParseHTMLFragmentFromPath(path string) html.Node {
	//do what old replacehtml did
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	root, err := html.Parse(strings.NewReader(string(file)))
	if err != nil {
		panic(err)
	}
	return *root
}

func ParseHTMLAsComponent(path string) []*html.Node {
	//do what old replacehtml did
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	rootList, err := html.ParseFragment(strings.NewReader(string(file)), &html.Node{
		Data:     "div",
		DataAtom: atom.Div,
		Type:     html.ElementNode,
	})
	if err != nil {
		panic("failed to parse component")
	}
	// root := &html.Node{
	// 	Data:     "div",
	// 	DataAtom: atom.Div,
	// 	Type:     html.ElementNode,
	// }
	// for _, n := range rootList {

	// }
	return rootList
}
