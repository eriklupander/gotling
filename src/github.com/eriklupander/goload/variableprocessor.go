package main
import (
"regexp"
"strings"
	"net/url"
)

var re = regexp.MustCompile("\\$\\{([a-zA-Z0-9]{0,})\\}")

func SubstParams(sessionMap map[string]string, textData string) string {
	if strings.ContainsAny(textData, "${") {
		res := re.FindAllStringSubmatch(textData, -1)
		for _, v := range res {
			textData = strings.Replace(textData, "${" + v[1] + "}", url.QueryEscape(sessionMap[v[1]]), 1)
		}
		return textData
	} else {
		return textData
	}
	return textData
}
