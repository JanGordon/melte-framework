// a layout should be name layout-layoutname.html and should be able to acces the current route using it and change aspects based on that
package compile

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func populateLayout(page html.Node, pagePath string) html.Node {
	dir := filepath.Dir(filepath.Join(pagePath))
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	fmt.Println("The path is :", dir)
	pageBytes, err := os.ReadFile(pagePath)
	pageString := string(pageBytes)
out:
	for {
		for _, f := range files {
			// fmt.Println(f.Name())
			if !f.IsDir() && strings.HasPrefix(f.Name(), "layout") {
				tFile, _ := os.ReadFile(filepath.Join(dir, f.Name()))
				fmt.Println("Layout template found: ", filepath.Join(dir, f.Name()))
				var re = regexp.MustCompile(`<slot></slot>`)
				fmt.Println(pageString)

				pageString = re.ReplaceAllString(string(tFile), pageString)
				fmt.Println(pageString)
				break out
			}

		}
		if strings.HasSuffix(dir, "/routes") {
			fmt.Println("Couldn't find layout for : ", pagePath)
			break
		}
		dir = filepath.Dir(filepath.Join(dir))
		files, err = os.ReadDir(dir)
		if err != nil {
			panic(err)
		}
	}
	template, err := html.Parse(strings.NewReader(pageString))
	if err != nil {
		panic(err)
	}
	return *template
}

// func replaceSlot(n *html.Node, page html.Node) bool {
// 	fmt.Println("checking if ", n.Data)
// 	if n.Data == "slot" {

// 		//removeAllchildren of slot paren

// 		for child := page.FirstChild; child != nil; child = child.NextSibling {
// 			fmt.Println("replacing slot with :", child.Data)
// 			n.Parent.InsertBefore(child, n)
// 		}
// 		n.Parent.RemoveChild(n)
// 		return true
// 	}
// 	fmt.Println("Searching next :", n.FirstChild.Data)
// 	// use node. InertBefor
// 	for child := n.FirstChild; child != nil; child = child.NextSibling {
// 		fmt.Println("searching children of :", child.Data)
// 		n := replaceSlot(child, page)
// 		fmt.Println("finding slot in ", child.Data)
// 		if n {
// 			return n
// 		} else {
// 			fmt.Println("no slot found in : ", child.Data)
// 		}
// 	}
// 	return false
// }

// func replaceSlot(n *html.Node, page *html.Node, fullTemplate *html.Node) *html.Node {
// 	fmt.Println("checking if ", n.Data)
// 	if n.Data == "slot" {

// 		//removeAllchildren of slot paren
// 		TreePrinter(page)
// 		page := ReplaceComponentWithHTML(*page, false, "")
// 		fmt.Println("Formatted:")
// 		TreePrinter(&page)
// 		//insertPage(&page, n)
// 		return fullTemplate
// 	}

// 	for child := n.FirstChild; child != nil; child = child.NextSibling {
// 		f := replaceSlot(child, page, fullTemplate)
// 		if f != nil {
// 			return f
// 		}
// 	}
// 	return nil
// }

func populateSlot(page html.Node) *html.Node {
	// return the slots parent containing slots children and original contents of body
	newSlot := &html.Node{
		Data:     "div",
		DataAtom: atom.Div,
	}
	newSlot.AppendChild(&page)

	return newSlot
}

func TreePrinter(n *html.Node) {
	fmt.Println("<", n.Data, ">")
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		TreePrinter(child)
	}
}

// func insertPage(page *html.Node, n *html.Node) {
// 	// insert all elents recursively into the slot.
// 	// pass element to insert before to this func
// 	fmt.Println(" pgaebody parent : ", page.FirstChild.Data)
// 	newSlotChild := &html.Node{
// 		Data:      page.FirstChild.LastChild.Data,
// 		DataAtom:  page.FirstChild.LastChild.DataAtom,
// 		Type:      page.FirstChild.LastChild.Type,
// 		Namespace: page.FirstChild.LastChild.Namespace,
// 		Attr:      page.FirstChild.LastChild.Attr,
// 	}
// 	page.FirstChild.LastChild.Parent = nil
// 	n.Parent.InsertBefore(newSlotChild, n)
// 	TreePrinter(n.Parent)
// 	fmt.Println("replacing slot with code : ", newSlotChild)
// 	for child := n.FirstChild; child != nil; child = child.NextSibling {
// 		insertPage(page, n)
// 	}
// }

// func ChildrenToArray(n *html.Node) []*html.Node {
// 	currArray := []*html.Node{}
// 	for child := n.FirstChild; child != nil; child = child.NextSibling {
// 		currArray = append(currArray, child)
// 		// currArray = append(currArray, ChildrenToArray(child)...)
// 	}
// 	return currArray
// }

// search through children of layout to find slot and then set the data of the slot to div and id app
// append all children of slot to the app div
