// a layout should be name layout-layoutname.html and should be able to acces the current route using it and change aspects based on that
package compile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func populateLayout(page []*html.Node, pagePath string, writeFile *os.File) []*html.Node {
	dir := pagePath
	fmt.Println("The path is :", dir)
out:
	for {
		if dir == "routes" {
			fmt.Println("Couldn't find layout for : ", pagePath)
			break
		}
		dir = filepath.Dir(filepath.Join(dir))
		fmt.Println("The nwe path is :", dir)
		files, err := os.ReadDir(dir)
		if err != nil {
			panic(err)
		}
		for _, f := range files {
			if !f.IsDir() && strings.HasPrefix(f.Name(), "layout") {
				tFile, err := os.ReadFile(pagePath)
				fmt.Println("Layout template found: ", filepath.Join(dir, f.Name()))
				template, err := html.ParseFragment(strings.NewReader(string(tFile)), &html.Node{
					Type:     html.ElementNode,
					Data:     "div",
					DataAtom: atom.Div,
				})
				if err != nil {
					panic(err)
				}

				for _, child := range template {
					foundSlot := replaceSlot(child, page)
					if foundSlot {
						fmt.Println("replaced slot with code")
						break
					}
				}

				BuildPage(ReplaceComponentWithHTML(template), pagePath, dir, false, true, false)

				break out
			}
		}
	}

	return page
}

func replaceSlot(n *html.Node, page []*html.Node) bool {
	if n.Data == "slot" {
		for _, child := range page {
			newChild := child
			n.Parent.AppendChild(newChild)
		}
		n.Parent.RemoveChild(n)
		return true
	}
	return false
}
