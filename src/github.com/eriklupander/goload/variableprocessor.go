package main
import (
"regexp"
"strings"
	"fmt"
	"net/url"
)

// Move this to some utility file
var re = regexp.MustCompile("\\$\\{([a-zA-Z0-9]{0,})\\}")

func SubstParams(sessionMap map[string]string, textData string) string {
	if strings.ContainsAny(textData, "${") {
		fmt.Println("Do " + textData)
		res := re.FindAllStringSubmatch(textData, -1)
		for _, v := range res {
			textData = strings.Replace(textData, "${" + v[1] + "}", url.QueryEscape(sessionMap[v[1]]), 1)
		}
		fmt.Println("Return " + textData)
		return textData
	} else {
		return textData
	}
	return textData
}
