package main
import (
	"log"
)

func buildActionList(t *TestDef) ([]Action, bool) {
	var valid bool = true
	actions := make([]Action, len(t.Actions), len(t.Actions))
	for _, element := range t.Actions {
		for key, value := range element {
            actionMap := value.(map[interface{}]interface{})
			switch key {
				case "sleep":

					sleepAction := NewSleepAction(actionMap)

					actions = append(actions, sleepAction)
					break
				case "http":

					httpAction := NewHttpAction(actionMap)
					actions = append(actions, httpAction)

					break
				case "tcp":
                    tcpAction := NewTcpAction(actionMap)
                    actions = append(actions, tcpAction)
				break
				default:
					valid = false
					log.Fatal("Unknown action type encountered: " + key)
					break
			}
		}
	}
	return actions, valid
}

func getBody(action map[interface{}]interface{}) string {
	var body string = ""
	if action["body"] != nil {
		body = action["body"].(string)
	}
	return body
}