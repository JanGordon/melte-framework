package compile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"rogchap.com/v8go"
)

type Context struct {
	absPath string
	v8Ctx   *v8go.Context
}

var CCount = 0
var Contexts []Context

func getContext(path string) *v8go.Context {
	var ctx *v8go.Context
	path, _ = filepath.Split(path)
	a, _ := filepath.Abs(filepath.Join(path, "out.html"))
	for c := 0; c < len(Contexts); c++ {
		// fmt.Printf("Conctext %s to match %s \n", Contexts[c].absPath, a)
		if Contexts[c].absPath == a {
			ctx = Contexts[c].v8Ctx
		}
	}
	if ctx == nil {
		fmt.Println("No context found: creating new at " + a)
		ctx = v8go.NewContext()
		Contexts = append(Contexts, Context{a, ctx})
	} else {
		fmt.Println("Found context")
	}

	return ctx
}
func ReplaceComponentWithHTML(root html.Node, findLayouts bool, pagePath string, c *v8go.Context) html.Node {
	CCount++
	var ctx = getContext(pagePath)

	replace(&root, pagePath, ctx)
	if findLayouts {
		//fmt.Println("Finding layout for ", pagePath)
		dir := filepath.Dir(filepath.Join(pagePath))
		files, err := os.ReadDir(dir)
	out:
		for {
			for _, f := range files {
				// fmt.Println(f.Name())
				if !f.IsDir() && strings.HasPrefix(f.Name(), "layout") {
					//fmt.Println("Layout template found: ", filepath.Join(dir, f.Name()))
					// first arg is template
					// second is slotInsert
					file, err := os.ReadFile(filepath.Join(dir, f.Name()))
					if err != nil {
						panic(err)
					}

					tempRender(pagePath, &root)
					newRootOut, err := os.ReadFile(pagePath)
					//fmt.Println(string(newRootOut))
					//fmt.Println("HHHHHHHHHHHHHHH")
					parsed := ParseHTMLFragmentFromString(string(newRootOut), pagePath)
					TreePrinter(&parsed)
					//fmt.Println("---------")
					root = ReplaceLayoutWithHTML(ParseHTMLFragmentFromString(string(file), pagePath), string(newRootOut), pagePath)
					break out
				}

			}
			if strings.HasSuffix(dir, "/routes") {
				fmt.Println("Couldn't find layout for : ", pagePath)
				break out
			}
			dir = filepath.Dir(filepath.Join(dir))
			files, err = os.ReadDir(dir)
			if err != nil {
				panic(err)
			}
		}
	}

	return root
}

func ReplaceLayoutWithHTML(root html.Node, slotInsert string, pagePath string) html.Node {
	CCount++
	var ctx = getContext(pagePath)
	newRoot := root
	replaceSlot((&root), slotInsert, pagePath, &newRoot, false, ctx)

	return root
}

func ReplaceCustomComponentWithHTML(root []*html.Node, pagePath string) []*html.Node {
	var ctx = getContext(pagePath)
	for _, child := range root {
		// OutScript := fmt.Sprintf(`const SELF = document.querySelector("[melte-id='%s']")`, child.Data+fmt.Sprintf("%d", CCount))

		// for _, a := range child.Attr {
		// 	if a.Key == "ssr" {
		// 		fmt.Println("Found script", child.Attr)

		// 		ctx.RunScript(child.FirstChild.Data, fmt.Sprintf("%v.js", CCount))
		// 		// v, _ := ctx.RunScript(fmt.Sprintf("result%v", CCount), "value.js") // return a value in JavaScript back to Go

		// 		fmt.Printf("Run scirpt : ", child.FirstChild.Data)
		// 	}

		// }
		replace(child, pagePath, ctx)
		CCount++

		// if err = html.Render(writeFile, root[child]); err != nil {
		// 	panic(err)
		// }
	}
	return root
}

func tempRender(path string, root *html.Node) {
	tempRootOut, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	tempRootOut.Write([]byte{})
	// have to render to out.html temporarliy because no way to convert *html.Node to proper []*html.Node

	for child := root.FirstChild; child != nil; child = child.NextSibling {

		// only render internal elemnt not head and stuff
		//fmt.Println("Rednedering ", child.Data, " to out.html temporarily")
		if child.Type == html.ElementNode {
			child.Attr = append(child.Attr, html.Attribute{
				Key: "tempRendered",
				Val: "1",
			})
		}

		if err = html.Render(tempRootOut, child); err != nil {
			panic(err)
		}
	}
	defer tempRootOut.Close()
}

var Scripts []html.Node
var HeadScripts []html.Node
var ScriptIDs []string

func replace(n *html.Node, pagePath string, ctx *v8go.Context) {
	CCount++
	if n.Type == html.ElementNode {
		//checkForMelteDef(n)
		// if n.Data == "script" {

		// 	// fmt.Println("Found script", n.FirstChild.Data)
		// 	OutScript := fmt.Sprintf(`const SELF = document.querySelector("[melte-id='%s']");`, n.Data+fmt.Sprintf("%d", CCount)) + "\n"
		// 	// We need to move the script to end and add module tag

		// 	scriptComponent := &html.Node{
		// 		Data:     "script",
		// 		Type:     html.ElementNode,
		// 		DataAtom: atom.Script,
		// 		Attr:     n.Attr,
		// 	}
		// 	newScript := &html.Node{
		// 		Data: OutScript + n.FirstChild.Data,
		// 		Type: html.TextNode,
		// 	}
		// 	scriptComponent.AppendChild(newScript)

		// 	scriptComponent.Attr = append(scriptComponent.Attr, html.Attribute{
		// 		Key: "melte-docpos",
		// 		Val: pagePath,
		// 	})

		// 	// n.RemoveChild(n.FirstChild)
		// 	// component[node].AppendChild(scriptComponent)
		// 	// if err != nil {
		// 	// 	panic(err)
		// 	// }
		// 	fmt.Println("Data adding to SCripts: ", scriptComponent.FirstChild.Data, " the end")
		// 	Scripts = append(Scripts, *scriptComponent)
		// 	ScriptIDs = append(ScriptIDs, fmt.Sprintf("out-%s%d.js", n.Data, CCount))
		// 	if n.Parent != nil {
		// 		n.Parent.RemoveChild(n)

		// 	}

		// }
		wd, err := os.Getwd()
		if err != nil {
			panic("failed to get working directory")
		}
		// if n.Data == "script" {
		// 	// OutScript := fmt.Sprintf(`const SELF = document.querySelector("[melte-id='%s']");`, n.Data+fmt.Sprintf("%d", CCount))
		// 	fmt.Println("Found script", n.Attr)
		// 	for _, a := range n.Attr {
		// 		if a.Key == "ssr" {

		// 			ctx.RunScript(n.FirstChild.Data, "main.js") //fmt.Sprintf("%v.js", CCount)
		// 			v, _ := ctx.RunScript("test", "value.js")   // return a value in JavaScript back to Go

		// 			fmt.Printf("Run scirpt : %s , %s", n.FirstChild.Data, v)
		// 		}

		// 	}
		// }

		// seeing if custom component exists

		_, err = os.ReadFile(filepath.Join(wd, "components", n.Data+".melte"))
		if err == nil {
			component := ReplaceCustomComponentWithHTML(ParseHTMLAsComponent(filepath.Join(wd, "components", n.Data+".melte")), filepath.Join(wd, "components")) // adds components scripts to Scripts

			n.Attr = append(n.Attr, html.Attribute{
				Key: "melte-id",
				Val: n.Data + fmt.Sprintf("%d", CCount),
			})
			// out:
			for _, child := range component {
				if child.Data == "script" {
					OutScript := fmt.Sprintf(`const SELF = document.querySelector("[melte-id='%s']")`, n.Data+fmt.Sprintf("%d", CCount))
					fmt.Println("hi", child.Attr)

					// this doesnt loop over scripts in html file: fix

					//fmt.Println("Found script")
					// We need to move the script to end and add module tag

					scriptComponent := &html.Node{
						Data:     "script",
						Type:     html.ElementNode,
						DataAtom: atom.Script,
						Attr:     child.Attr,
					}
					for child := child.FirstChild; child != nil; child = child.NextSibling {
						newScript := &html.Node{
							Data: OutScript + child.Data,
							Type: html.TextNode,
						}
						scriptComponent.AppendChild(newScript)
					}
					scriptComponent.Attr = append(scriptComponent.Attr, html.Attribute{
						Key: "melte-docpos",
						Val: pagePath,
					})

					child.RemoveChild(child.FirstChild)
					// component[node].AppendChild(scriptComponent)
					if err != nil {
						panic(err)
					}
					Scripts = append(Scripts, *scriptComponent)
					ScriptIDs = append(ScriptIDs, fmt.Sprintf("out-%s%d.js", n.Data, CCount))
					if child.Parent != nil {
						child.Parent.RemoveChild(child)

					}
				} else {
					n.AppendChild(child)
				}
			}
		}

	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		replace(child, pagePath, ctx)
	}
}

func replaceSlot(n *html.Node, slotInsert string, pagePath string, rootCopy *html.Node, cont bool, ctx *v8go.Context) {
	// maybe adding slot to itslef and then adding thme evry loop
	CCount++
	done := cont
	if n.Type == html.ElementNode {
		//checkForMelteDef(n)
		wd, err := os.Getwd()
		if err != nil {
			panic("failed to get working directory")
		}

		// seeing if custom component exists
		_, err = os.ReadFile(filepath.Join(wd, "components", n.Data+".melte"))
		// TreePrinter(n)

		// if n.Data == "script" {
		// 	// OutScript := fmt.Sprintf(`const SELF = document.querySelector("[melte-id='%s']");`, n.Data+fmt.Sprintf("%d", CCount))
		// 	fmt.Println("Found script", n.Attr)
		// 	for _, a := range n.Attr {
		// 		if a.Key == "ssr" {

		// 			ctx.RunScript(n.FirstChild.Data, "main.js") //fmt.Sprintf("%v.js", CCount)
		// 			v, _ := ctx.RunScript("test", "value.js")   // return a value in JavaScript back to Go

		// 			fmt.Printf("Run scirpt : %s , %s", n.FirstChild.Data, v)
		// 		}

		// 	}
		// }
		if err == nil {

			component := ReplaceCustomComponentWithHTML(ParseHTMLAsComponent(filepath.Join(wd, "components", n.Data+".melte")), filepath.Join(wd, "components")) // adds components scripts to Scripts

			n.Attr = append(n.Attr, html.Attribute{
				Key: "melte-id",
				Val: n.Data + fmt.Sprintf("%d", CCount),
			})
			for _, child := range component {
				if child.Data == "script" {
					OutScript := fmt.Sprintf(`const SELF = document.querySelector("[melte-id='%s']")`, n.Data+fmt.Sprintf("%d", CCount))
					// We need to move the script to end and add module tag
					fmt.Println("hi", child.Attr)

					for _, a := range child.Attr {
						fmt.Println(a.Key)
						if a.Key == "ssr" {
							ctx.RunScript(OutScript+child.FirstChild.Data, fmt.Sprintf("%v.js", CCount))
							// v, _ := ctx.RunScript(fmt.Sprintf("result%v", CCount), "value.js") // return a value in JavaScript back to Go

							fmt.Printf("Run scirpt : ", OutScript+child.FirstChild.Data)
						}

					}

					scriptComponent := &html.Node{
						Data:     "script",
						Type:     html.ElementNode,
						DataAtom: atom.Script,
						Attr:     child.Attr,
					}
					for child := child.FirstChild; child != nil; child = child.NextSibling {
						newScript := &html.Node{
							Data: OutScript + child.Data,
							Type: html.TextNode,
						}
						scriptComponent.AppendChild(newScript)
					}
					scriptComponent.Attr = append(scriptComponent.Attr, html.Attribute{
						Key: "melte-docpos",
						Val: pagePath,
					})

					child.RemoveChild(child.FirstChild)
					// component[node].AppendChild(scriptComponent)
					if err != nil {
						panic(err)
					}
					Scripts = append(Scripts, *scriptComponent)
					ScriptIDs = append(ScriptIDs, fmt.Sprintf("out-%s%d.js", n.Data, CCount))

				} else {
					n.AppendChild(child)
				}
			}
		}
		if n.Data == "slot" && !isChildOf(n, "slot") && !done {
			// us epagepath
			// file, _ := os.ReadFile(filepath.Join(wd, "components/slotTester.melte"))

			// need to use specific child of slotInsert
			// component := ReplaceCustomComponentWithHTML(ParseHTMLStringAsComponent(string(file)), pagePath) // adds components scripts to Scripts
			component := ReplaceCustomComponentWithHTML(ParseHTMLStringAsComponent(slotInsert, pagePath), pagePath) // adds components scripts to Scripts
			//fmt.Println("----------------------------")

			n.Attr = append(n.Attr, html.Attribute{
				Key: "melte-id",
				Val: n.Data + fmt.Sprintf("%d", CCount),
			})
			// for e := range component {
			// 	TreePrinter(component[e])

			// }
			for _, child := range component {
				// TreePrinter(child)
				//fmt.Println("Looping over elements in slot : ", n.Data)
				if child.Data == "script" {
					OutScript := fmt.Sprintf(`const SELF = document.querySelector("[melte-id='%s']" )`, n.Data+fmt.Sprintf("%d", CCount))
					// We need to move the script to end and add module tag

					// if there is no child of the script then set Data to reading of src
					// to do: add the attributes back on to script

					scriptComponent := &html.Node{
						Data:     "script",
						Type:     html.ElementNode,
						DataAtom: atom.Script,
						Attr:     child.Attr,
					}
					for child := child.FirstChild; child != nil; child = child.NextSibling {
						newScript := &html.Node{
							Data: OutScript + child.Data,
							Type: html.TextNode,
						}
						scriptComponent.AppendChild(newScript)
					}
					scriptComponent.Attr = append(scriptComponent.Attr, html.Attribute{
						Key: "melte-docpos",
						Val: pagePath,
					})

					//child.Parent.RemoveChild(child)
					// component[node].AppendChild(scriptComponent)
					// if err != nil {
					// 	panic(err)
					// }
					Scripts = append(Scripts, *scriptComponent)
					ScriptIDs = append(ScriptIDs, fmt.Sprintf("out-%s%d.js", n.Data, CCount))

				} else {
					n.AppendChild(child)
				}

				// n.Data = "melte-null-slot"
				// n.Parent.RemoveChild(n)
				done = true

				// }
				//stop loop once one slot found
				// in futur make this happen once entire doc parsed
			}
		}
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		// if done {
		// 	break
		// }
		replaceSlot(child, slotInsert, pagePath, rootCopy, done, ctx)
		if child.Data == "slot" {
			break
		}
	}
}

func ParseHTMLFragmentFromPath(path string) html.Node {
	//do what old replacehtml did
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	file = []byte(checkHTMLFile(string(file), path, getContext(path)))
	root, err := html.Parse(strings.NewReader(string(file)))
	if err != nil {
		panic(err)
	}
	return *root
}

func ParseHTMLFragmentFromString(file string, path string) html.Node {
	//do what old replacehtml did
	root, err := html.Parse(strings.NewReader(checkHTMLFile(file, path, getContext(path))))
	if err != nil {
		panic(err)
	}
	return *root
}

// func ParseHTMLNodeToChildren(node *html.Node) []*html.Node {
// 	var newNodes []*html.Node
// 	for child := node.FirstChild; child != nil; child = child.NextSibling {
// 		newNodes = append(newNodes, child)
// 	}
// 	return newNodes
// }

func ParseHTMLStringAsComponent(root string, path string) []*html.Node {
	root = checkHTMLFile(root, path, getContext(path))
	rootList, err := html.ParseFragment(strings.NewReader(root), &html.Node{
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

func ParseHTMLAsComponent(path string) []*html.Node {
	//do what old replacehtml did
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	file = []byte(checkHTMLFile(string(file), path, getContext(path)))
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

func isChildOf(child *html.Node, parentName string) bool {
	for parent := child.Parent; parent != nil; parent = parent.Parent {
		if parent.Data == parentName {
			return true
		}
	}
	return false
}

// func checkForMelteDef(n *html.Node) {
// 	if n.Data == "script" && n.Type == html.ElementNode {
// 		lines := strings.Split(n.FirstChild.Data, "\n")
// 		for lineIndex, line := range lines {
// 			l := strings.Trim(line, " ")
// 			if strings.HasPrefix(l, "//=") {
// 				l = strings.Trim(l, "//=")
// 				if strings.HasPrefix(l, "keep state:") {
// 					l = strings.TrimSpace(strings.Trim(l, "keep state:"))
// 					if strings.HasPrefix(l, "js") {
// 						scriptNode := &html.Node{
// 							Data:     "script",
// 							DataAtom: atom.Script,
// 							Type:     html.ElementNode,
// 						}
// 						scriptNode.AppendChild(&html.Node{
// 							Data: lines[lineIndex+1],
// 							Type: html.TextNode,
// 						})
// 						HeadScripts = append(HeadScripts, *scriptNode)
// 					} else if strings.HasPrefix(l, "url") {
// 						// jsDict := "{}"
// 						// let js modify url when var chnage
// 						// when url with ?variable=10 router should serve js with variable embedded in js
// 						// and if possible update the html fragments with reactive html in them before serving
// 					}
// 				}
// 			}
// 		}

// 	} else {
// 		fmt.Println("not a melte def script")
// 	}
// }
