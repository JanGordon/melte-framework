package compile

import (
	"fmt"
	"path/filepath"
	"strings"

	v8 "rogchap.com/v8go"
)

type ServerFunc struct {
	Route string
	F     func(string, string) string
}

type Token struct {
	open         bool
	startBracket int
	depth        int
	endBracket   int
	startLine    int
	endLine      int
}

var ServerFunctions []ServerFunc

func checkHTMLFile(file string, path string) string {
	// chars := []rune(file)
	tokenDepth := 0
	tokens := []Token{}
	lines := strings.Split(file, "\n")
	currentCharNum := 0
	for currentLine, line := range lines {
		for c := 0; c < len(line); c++ {

			char := string(line[c])
			//to prevent error of chars[-1]
			if c > 0 {
				if char == "{" && string(line[c-1]) == "{" {
					tokens = append(tokens, Token{
						open:         true,
						startBracket: currentCharNum + c,
						depth:        tokenDepth,
						startLine:    currentLine,
					})
					fmt.Println("Opeing token")
					tokenDepth++
				} else if char == "}" && string(line[c-1]) == "}" {
					fmt.Println("closing, ", char, c, currentLine, currentCharNum+c)

					for p := len(tokens) - 1; p >= 0; p-- {

						if tokens[p].open {
							tokens[p].open = false
							tokens[p].endBracket = currentCharNum + c
							tokens[p].endLine = currentLine
							break
						}
					}
				}
			}
		}
		currentCharNum += len(line) + 1
	}
	// newChars := []rune(file)
	for p := len(tokens) - 1; p >= 0; p-- {
		t := tokens[p]
		//fmt.Println(t.open)
		if t.open {
			panic(fmt.Errorf("unclosed {{ at: %v", t.startLine))
		} else {
			split := strings.Split(lines[t.startLine], "{{")
			if strings.HasPrefix(split[1], "for") {

			}
			option := strings.Split(split[1], "||")
			js := option[0]
			if len(option) > 1 {
				reloadOn := strings.Split(option[1], ",")
				for _, event := range reloadOn {
					event = strings.TrimSpace(event)
					if event == "onbuild" {
						// run js
						// need to add mutation object to update whne changed
						ctx := v8.NewContext()
						fmt.Println("JS:", js)

						jsForReload := "var result = '';" + js + "{result += '" + strings.Join((lines[(t.startLine+1):(t.endLine)]), "") + "'}"
						fmt.Printf("var result = '';" + js + "{result += '" + strings.Join((lines[(t.startLine+1):(t.endLine)]), "") + "'}")

						ctx.RunScript(jsForReload, "main.js") // any functions previously added to the context can be called
						// ctx.RunScript("var result = 'hello'; for (let i of [{hello : 'g'},2,3,4]){result += '    <h1></h1>'}", "main.js") // any functions previously added to the context can be called
						val, _ := ctx.RunScript("result", "value.js") // return a value in JavaScript back to Go
						fmt.Println("Fixed loop: ", val, " //")
						fmt.Printf("After brackets : %v || \n", file[t.endBracket+2:])
						file = file[:t.startBracket-1] + fmt.Sprint(val) + file[t.endBracket+1:]

					} else if strings.HasPrefix(event, "onload") {
						if strings.HasSuffix(event, "server") {
							fmt.Println("The path", path)
							dir, _ := filepath.Split(path)
							fmt.Println("Startbracket :", t.endBracket)
							UpdateForLoop(dir, js, t, lines)

						}
					}
				}
			} else {
				// run js
				// need to add mutation object to update whne changed
				ctx := v8.NewContext()
				jsForReload := "var result = '';" + js + "{result += `" + strings.Join((lines[(t.startLine+1):(t.endLine)]), "") + "`}"
				fmt.Printf("var result = '';" + js + "{result += `" + strings.Join((lines[(t.startLine+1):(t.endLine)]), "") + "`}")

				ctx.RunScript(jsForReload, "main.js") // any functions previously added to the context can be called
				// ctx.RunScript("var result = 'hello'; for (let i of [{hello : 'g'},2,3,4]){result += '    <h1></h1>'}", "main.js") // any functions previously added to the context can be called
				val, _ := ctx.RunScript("result", "value.js") // return a value in JavaScript back to Go
				fmt.Println("Fixed loop: ", val, " //")
				fmt.Printf("After brackets : %v || \n", file[t.endBracket+2:])
				file = file[:t.startBracket-1] + fmt.Sprint(val) + file[t.endBracket+1:]

			}
			// file = strings.ReplaceAll(file, strings.Join((lines[(t.startLine):(t.endLine)]), "\n")+"\n}}", fmt.Sprint(val))

		}
	}
	return file
}

func AddServeFunc(route string, f func(string, string) string) {
	ServerFunctions = append(ServerFunctions, ServerFunc{route, f})
}

func UpdateForLoop(dir string, js string, t Token, lines []string) {
	AddServeFunc(dir+"out.html", func(s string, file string) string {
		ctx := v8.NewContext()
		jsForReload := "var result = '';" + js + "{result += '" + strings.Join((lines[(t.startLine+1):(t.endLine)]), "") + "'}"
		fmt.Printf("var result = '';" + js + "{result += '" + strings.Join((lines[(t.startLine+1):(t.endLine)]), "") + "'}")

		ctx.RunScript(jsForReload, "main.js") // any functions previously added to the context can be called
		// ctx.RunScript("var result = 'hello'; for (let i of [{hello : 'g'},2,3,4]){result += '    <h1></h1>'}", "main.js") // any functions previously added to the context can be called
		val, _ := ctx.RunScript("result", "value.js") // return a value in JavaScript back to Go
		fmt.Println("Fixed loop: ", t.endBracket, " //")
		fmt.Printf("After brackets : %v || \n", file[t.endBracket+2:])
		file = file[:t.startBracket-1] + fmt.Sprint(val) + file[t.endBracket+1:]

		return file
	})
}
