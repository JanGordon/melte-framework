package compile

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
	v8 "rogchap.com/v8go"
)

func GetJSForCSR(path string) string {

	return "console.log('Failed to client side render')"
}

func RunInitialScripts(path string, file string) {
	ctx := getContext(path)
	doc, err := html.Parse(strings.NewReader(file))
	if err != nil {
		panic(err)
	}
	runJS(doc.FirstChild, ctx)
}

func runJS(n *html.Node, ctx *v8.Context) {
	if n.Data == "script" {
		ctx.RunScript(n.FirstChild.Data, fmt.Sprintf("v%v.js", CCount))
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		runJS(child, ctx)
	}
	CCount++
}
