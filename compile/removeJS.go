package compile

import (
	"bytes"
	"errors"
	"io"
	"os"

	"golang.org/x/net/html"
)

func script(doc *html.Node) ([]*html.Node, error) {
	var script []*html.Node
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "script" {
			script = append(script, node)
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	if script != nil {
		return script, nil
	}
	return nil, errors.New("missing <body> in the node tree")
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

func RemoveJS(filepath string) ([]string, string) {
	f, err := os.Open(filepath)
	if err != nil {
		os.Exit(1)
	}
	doc, _ := html.Parse(f)
	sn, err := script(doc)
	if err != nil {
		os.Exit(1)
	}
	fileText, err := os.ReadFile(filepath)
	if err != nil {
		os.Exit(1)
	}
	var scripts []string
	var html string
	for scriptNode := range sn {
		script := renderNode(sn[scriptNode])
		html = string(bytes.Replace([]byte(fileText), []byte(script), []byte(""), -1))

		script = string(bytes.Replace([]byte(script), []byte("<script>"), []byte(""), -1))
		script = string(bytes.Replace([]byte(script), []byte("</script>"), []byte(""), -1))
		scripts = append(scripts, script)
	}
	return scripts, html
}
