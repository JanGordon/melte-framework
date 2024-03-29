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
	// c := getContext(inPath)
	// if ctx needed here then uncomment this
	os.Truncate(outPath, 0)
	writeFile, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE, 0600)

	//newPage, err := html.Parse(strings.NewReader(string(f)))
	if err != nil {
		panic(err)
	}
	importLines := ""
	scriptExceptImports := ""
	HeadScript := "var _"

	for script := range Scripts {

		if inlineJS {
			// cat be asked wht this
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
			s, i, h := transJS(scriptExceptImports, importLines, HeadScript, Scripts[script].FirstChild.Data, script)
			scriptExceptImports += fmt.Sprintf("{\n// script for %s\n %s}", fmt.Sprintf("%s", ScriptIDs[script]), s)
			importLines += i
			HeadScript = h
		}
		// newPage.AppendChild(scriptC)

	}
	for _, script := range ExternalScripts {
		data, err := os.ReadFile(script)
		if err != nil {
			panic(err)
		}
		s, i, h := transJS("", "", HeadScript, string(data), 0)
		HeadScript = h
		fmt.Println("Building external scripts")
		BuildScriptFile(i+"\n"+s, script+".melte-out.js")

	}
	HeadScriptNode := &html.Node{
		Data:     "script",
		DataAtom: atom.Script,
		Type:     html.ElementNode,
	}
	HeadScriptNode.AppendChild(&html.Node{
		Data: HeadScript,
		Type: html.TextNode,
	})
	root.LastChild.FirstChild.AppendChild(HeadScriptNode)
	cwd, err := os.Getwd()
	file := importLines + "\n" + scriptExceptImports
	// fmt.Println(file)
	BuildScriptFile(file, filepath.Join(outPathJS, "out.js"))
	//fmt.Println(outPathJS)
	scriptC := &html.Node{
		Data:     "script",
		Type:     html.ElementNode,
		DataAtom: atom.Script,
	}
	scriptC.Attr = append(scriptC.Attr, html.Attribute{
		Key: "src",
		Val: filepath.Join(strings.Replace(strings.Replace(outPathJS, cwd, "", 1), "routes", "", 1), "out.js"),
	})
	root.LastChild.AppendChild(scriptC)
	scriptFlamethrower := &html.Node{
		Data:     "script",
		Type:     html.ElementNode,
		DataAtom: atom.Script,
	}
	scriptFlamethrower.Attr = append(scriptFlamethrower.Attr, html.Attribute{
		Key: "src",
		Val: "/clientSideRouting/out.js",
	})
	scriptFlamethrower.Attr = append(scriptFlamethrower.Attr, html.Attribute{
		Key: "defer",
		Val: "",
	})
	root.LastChild.FirstChild.AppendChild(scriptFlamethrower)
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
		root.LastChild.FirstChild.AppendChild(scriptDev)

	}

	// adding all css together and to html
	fmt.Println("CSS", CSST)
	c := ""
	for _, i := range CSST {
		fmt.Println(i)

		c += i
	}
	styleNode := &html.Node{
		Data:     "style",
		Type:     html.ElementNode,
		DataAtom: atom.Style,
	}
	textNode := &html.Node{
		Data:     c,
		Type:     html.TextNode,
		DataAtom: 0,
	}
	root.LastChild.AppendChild(styleNode)
	styleNode.AppendChild(textNode)

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
	// HeadScript = ",var _"
	ScriptIDs = nil
	CCount = 0
	CSST = nil
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

func transJS(importRemovedLines string, importLines string, HeadScript string, scriptData string, scriptIndex int) (string, string, string) {

	// docPos := ""
	// for _, a := range Scripts[script].Attr {
	// 	if a.Key == "melte-docpos" {
	// 		docPos = a.Val
	// 	}
	// }
	// for _, a := range Scripts[script].Attr {
	// 	if a.Key == "src" {
	// 		scriptDataNew, _ := os.ReadFile(filepath.Join(docPos, a.Val))
	// 		scriptData = string(scriptDataNew)
	// 		// should also use aliases
	// 		//make sure this file is opened with write state
	// 		// ./ should be realtive to
	// 	}
	// }
	lines := strings.Split(scriptData, "\n")
	customCount := 0
	for lineIndex, line := range lines {
		l := strings.Trim(line, " ")
		if strings.HasPrefix(strings.TrimSpace(line), "import") {
			importLines += fmt.Sprintf("%s\n", line)
		} else if strings.HasPrefix(l, "//=") {
			l = strings.Trim(l, "//=")

			if strings.HasPrefix(l, "keep state: ") {
				l = strings.TrimSpace(strings.Replace(l, "keep state: ", "", 1))
				//fmt.Println("Found melte custom : ", l, lines[lineIndex+1])

				if strings.HasPrefix(l, "js") {
					decLine := strings.TrimSpace(lines[lineIndex+1])
					if strings.HasPrefix(decLine, "var") {
						fmt.Println("adding global component variable", strings.Replace(decLine, "var", ",", 1))
						HeadScript += strings.Replace(decLine, "var", ",", 1)
						lines[lineIndex+1] = ""
						// global and state preserevrd variable

					} else if strings.HasPrefix(decLine, "let") {
						varName := strings.Split(strings.Replace(decLine, "let", ",", 1), "=")
						HeadScript += ", " + strings.TrimSpace(strings.Replace(varName[0], ", ", "", 1)) + strings.Replace(strings.Replace(ScriptIDs[scriptIndex], "out-", "", 1), ".js", "", 1) + " = " + varName[1]
						//fmt.Println("adding this to head: " + "let " + strings.Replace(varName[0], ", ", "", 1) + " = " + strings.TrimSpace(strings.Replace(varName[0], ", ", "", 1)) + strings.Replace(strings.Replace(ScriptIDs[script], "out-", "", 1), ".js", "", 1))
						lines[lineIndex+1] = "let " + strings.Replace(varName[0], ", ", "", 1) + "=" + strings.TrimSpace(strings.Replace(varName[0], ", ", "", 1)) + strings.Replace(strings.Replace(ScriptIDs[scriptIndex], "out-", "", 1), ".js", "", 1)
						fmt.Println("added private component variable", lines[lineIndex+1])
						// private and state preseverd variable
					}
				} else if strings.HasPrefix(l, "url") {
					// jsDict := "{}"
					// let js modify url when var chnage
					// when url with ?variable=10 router should serve js with variable embedded in js
					// and if possible update the html fragments with reactive html in them before serving
				}
				customCount++
			}
		} else {
			//fmt.Println(line)
			importRemovedLines += line + "\n"
		}
	}
	fmt.Println(HeadScript)
	return importRemovedLines, importLines, HeadScript
}
