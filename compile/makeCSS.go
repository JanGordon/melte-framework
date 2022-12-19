package compile

import (
	"fmt"
	"strings"

	"github.com/vanng822/css"
)

func makeCSS(currentCSS []string, nodeType string, cssT string, nodeID string) []string {
	fmt.Println("Making CSS...")
	for _, style := range currentCSS {
		if strings.HasPrefix(style, fmt.Sprintf("/* %v */", nodeType)) {
			return currentCSS
			// to check if style for this component already exists in this context
		}
	}
	ss := css.Parse(cssT)
	rules := ss.GetCSSRuleList()
	newCSSString := ""
	for _, rule := range rules {
		rule.Style.Selector = css.NewCSSValue(fmt.Sprintf("[melte-id='%v'] %v", nodeID, rule.Style.Selector.Text()))
		newCSSString += rule.Style.Text()
	}

	currentCSS = append(currentCSS, fmt.Sprintf("/* %v */\n%v", nodeType, newCSSString))
	return currentCSS
}
