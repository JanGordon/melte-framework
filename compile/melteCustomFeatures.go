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
	identifier   string
}

var ServerFunctions []ServerFunc

func checkHTMLFile(file string, path string, ctx *v8.Context) string {
	// chars := []rune(file)
	RunInitialScripts(path, file)
	ctx.RunScript("let test = 'hello';", "rand.js")
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
						identifier:   "{{",
					})
					fmt.Println("adding for loop", currentCharNum+c)

					tokenDepth++
				} else if char == "}" && string(line[c-1]) == "}" {

					for p := len(tokens) - 1; p >= 0; p-- {

						if tokens[p].open && tokens[p].identifier == "{{" {
							tokens[p].open = false
							tokens[p].endBracket = currentCharNum + c
							tokens[p].endLine = currentLine
							break
						}
					}
				} else if char == "!" && string(line[c-1]) == "{" {
					tokens = append(tokens, Token{
						open:         true,
						startBracket: currentCharNum + c,
						depth:        tokenDepth,
						startLine:    currentLine,
						identifier:   "{!",
					})
					tokenDepth++
				} else if char == "}" && string(line[c-1]) == "!" {
					for p := len(tokens) - 1; p >= 0; p-- {

						if tokens[p].open && tokens[p].identifier == "{!" {
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
	resultNum := 0
	offset := 0
	lineOffset := 0
	for p := len(tokens) - 1; p >= 0; p-- {
		t := tokens[p]
		if t.open {
			panic(fmt.Errorf("unclosed {{ at: %v", t.startLine))
		} else {

			if t.identifier == "{{" {
				fmt.Println("running for loop", t.startLine)
				split := strings.Split(lines[t.startLine], "{{")
				// if strings.HasPrefix(split[1], "for") {

				// }

				option := strings.Split(split[1], "||")
				js := option[0]
				if len(option) > 1 {
					reloadOn := strings.Split(option[1], ",")
					for _, event := range reloadOn {
						event = strings.TrimSpace(event)
						if event == "onbuild" {
							// run js
							// need to add mutation object to update whne changed

							jsForReload := fmt.Sprintf("var result%v = '';", resultNum) + js + fmt.Sprintf("{result%v += `", resultNum) + strings.Join((lines[(t.startLine+1):(t.endLine)]), "") + "`}"
							ctx.RunScript(jsForReload, "main.js") // any functions previously added to the context can be called
							// ctx.RunScript("var result = 'hello'; for (let i of [{hello : 'g'},2,3,4]){result += '    <h1></h1>'}", "main.js") // any functions previously added to the context can be called
							val, _ := ctx.RunScript(fmt.Sprintf("result%v", resultNum), "value.js") // return a value in JavaScript back to Go
							ofLen := len(file[t.startBracket-1+offset : t.endBracket+1+offset])
							lineOfLen := strings.Count(file[t.startBracket-1+offset:t.endBracket+1+offset], "\n")
							file = file[:t.startBracket-1] + fmt.Sprint(val) + file[t.endBracket+1:]

							offset = len(fmt.Sprint(val)) - ofLen
							lineOffset = strings.Count(fmt.Sprint(val), "\n") - lineOfLen
							computeOffsets(tokens, offset, lineOffset, t)
						} else if strings.HasPrefix(event, "onload") {
							if strings.HasSuffix(event, "server") {
								dir, _ := filepath.Split(path)
								UpdateForLoop(dir, js, t, lines)

							}
						}
					}
				} else {
					// run js
					// need to add mutation object to update whne changed
					jsForReload := fmt.Sprintf("var result%v = '';", resultNum) + js + fmt.Sprintf("{result%v += `", resultNum) + strings.Join((lines[(t.startLine+1):(t.endLine)]), "") + "`}"
					ctx.RunScript(jsForReload, "main.js") // any functions previously added to the context can be called
					// ctx.RunScript("var result = 'hello'; for (let i of [{hello : 'g'},2,3,4]){result += '    <h1></h1>'}", "main.js") // any functions previously added to the context can be called
					val, _ := ctx.RunScript(fmt.Sprintf("result%v", resultNum), "value.js") // return a value in JavaScript back to Go
					fmt.Println(offset)
					ofLen := t.endBracket + 1 - t.startBracket - 1
					lineOfLen := strings.Count(file[t.startBracket-1:t.endBracket+1], "\n")
					file = file[:t.startBracket-1] + fmt.Sprint(val) + file[t.endBracket+1:]
					offset = len(fmt.Sprint(val)) - ofLen
					lineOffset = strings.Count(fmt.Sprint(val), "\n") - lineOfLen
					fmt.Println(tokens[1].endBracket)
					computeOffsets(tokens, offset, lineOffset, t)
					fmt.Println(offset)
					// indexes of tokens become messed up
				}
			} else if t.identifier == "{!" {
				js := fmt.Sprintf("var result%v = ", resultNum) + string(file[t.startBracket+1:t.endBracket-1])
				ctx.RunScript(js, fmt.Sprintf("main%v.js", CCount))
				fmt.Println("Found inline js :", js)
				val, _ := ctx.RunScript(fmt.Sprintf("result%v", resultNum), "value.js")
				ofLen := len(file[t.startBracket-1+offset : t.endBracket+1+offset])
				lineOfLen := strings.Count(file[t.startBracket-1+offset:t.endBracket+1+offset], "\n")
				file = file[:t.startBracket-1] + fmt.Sprintf("<melte-reload js='%s'>", js) + fmt.Sprint(val) + "</melte-reload>" + file[t.endBracket+1:]
				offset = len(fmt.Sprintf("<melte-reload js='%s'>", js)+fmt.Sprint(val)+"</melte-reload>") - ofLen
				lineOffset = strings.Count(fmt.Sprintf("<melte-reload js='%s'>", js)+fmt.Sprint(val)+"</melte-reload>", "\n") - lineOfLen
				computeOffsets(tokens, offset, lineOffset, t)
				// offset is only needed on compoents below point of addition
				fmt.Printf("oflen = %v, ofc = %v", ofLen, len(fmt.Sprintf("<melte-reload js='%s'>", js)+fmt.Sprint(val)+"</melte-reload>"))
			}

			// file = strings.ReplaceAll(file, strings.Join((lines[(t.startLine):(t.endLine)]), "\n")+"\n}}", fmt.Sprint(val))

		}
		resultNum++
	}
	return file
}

func computeOffsets(tokens []Token, offset int, lineOffset int, currentToken Token) {
	for i, t := range tokens {
		if t.startBracket > currentToken.endBracket {
			fmt.Println(t.startBracket)
			t.startBracket = t.startBracket + offset
			fmt.Println(t.startBracket)
			t.startLine = t.startLine + lineOffset
			fmt.Println("Adding offset of :", offset, "to", t.identifier, i, "due to", currentToken.identifier)
		}
		if t.endBracket > currentToken.endBracket {
			fmt.Println(t.endBracket)

			t.endBracket = t.endBracket + offset
			fmt.Println(t.endBracket)
			t.endLine = t.endLine + lineOffset
			fmt.Println("Adding offset of :", offset, "to end of", t.identifier, i, "due to", currentToken.identifier)

		}
	}
}

func AddServeFunc(route string, f func(string, string) string) {
	ServerFunctions = append(ServerFunctions, ServerFunc{route, f})
}

func UpdateForLoop(dir string, js string, t Token, lines []string) {
	AddServeFunc(dir+"out.html", func(s string, file string) string {
		ctx := v8.NewContext()
		jsForReload := "var result = '';" + js + "{result += '" + strings.Join((lines[(t.startLine+1):(t.endLine)]), "") + "'}"
		ctx.RunScript(jsForReload, "main.js") // any functions previously added to the context can be called
		// ctx.RunScript("var result = 'hello'; for (let i of [{hello : 'g'},2,3,4]){result += '    <h1></h1>'}", "main.js") // any functions previously added to the context can be called
		val, _ := ctx.RunScript("result", "value.js") // return a value in JavaScript back to Go
		file = file[:t.startBracket-1] + fmt.Sprint(val) + file[t.endBracket+1:]

		return file
	})
}
