package features

import (
	"fmt"
	"path/filepath"
	"strings"

	v8 "rogchap.com/v8go"
)

type Token struct {
	open         bool
	startBracket int
	depth        int
	endBracket   int
	startLine    int
	endLine      int
}

type ServerFunc struct {
	Route string
	F     func(string, string) string
}

func FormatForLoop(file string, split []string, path string, lines []string, t Token, ServerFunctions *[]ServerFunc) string {
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
					AddServeFunc(ServerFunctions, dir+"out.html", func(s string, file string) string {
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
	return file
}

func AddServeFunc(ServerFunctions *[]ServerFunc, route string, f func(string, string) string) {
	*ServerFunctions = append(*ServerFunctions, ServerFunc{route, f})
}
