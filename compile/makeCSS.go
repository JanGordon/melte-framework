package compile

import (
	"fmt"
	"strings"
)

func makeCSS(currentCSS []string, nodeType string, css string) []string {
	for _, style := range currentCSS {
		if strings.HasPrefix(style, fmt.Sprintf("/* %v */", nodeType)) {
			return currentCSS
			// to check if style for this component already exists in this context
		}
	}
	currentCSS = append(currentCSS, fmt.Sprintf("/* %v */\n%v", nodeType, css))
	return currentCSS
}
