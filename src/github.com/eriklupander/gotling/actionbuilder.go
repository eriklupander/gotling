/**
The MIT License (MIT)

Copyright (c) 2015 ErikL

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package main
import (
	"log"
"os"
"io/ioutil"
)

func buildActionList(t *TestDef) ([]Action, bool) {
	var valid bool = true
	actions := make([]Action, len(t.Actions), len(t.Actions))
	for _, element := range t.Actions {
		for key, value := range element {
			var action Action
            actionMap := value.(map[interface{}]interface{})
			switch key {
				case "sleep":
					action = NewSleepAction(actionMap)
					break
				case "http":
					action = NewHttpAction(actionMap)
					break
				case "tcp":
					action = NewTcpAction(actionMap)
				case "udp":
					action = NewUdpAction(actionMap)
					break
				default:
					valid = false
					log.Fatal("Unknown action type encountered: " + key)
					break
			}
			if valid {
				actions = append(actions, action)
			}
		}
	}
	return actions, valid
}

func getBody(action map[interface{}]interface{}) string {
	//var body string = ""
	if action["body"] != nil {
		return action["body"].(string)
	} else {
		return "";
	}
}

func getTemplate(action map[interface{}]interface{}) string {
	if action["template"] != nil {
		var templateFile = action["template"].(string)
		dir, _ := os.Getwd()
		templateData, _ := ioutil.ReadFile(dir + "/templates/" + templateFile)
		return string(templateData)
	} else {
		return ""
	}
}