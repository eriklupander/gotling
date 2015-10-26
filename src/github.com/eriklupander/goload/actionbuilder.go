package main
import "github.com/eriklupander/goload/model"

func buildActionList(t *model.TestDef) []interface{} {
	actions := make([]interface{}, len(t.Actions), len(t.Actions))
	for _, element := range t.Actions {
		for key, value := range element {

			switch key {
			case "sleep":
				action := value.(map[interface{}]interface{})
				actions = append(actions, model.SleepAction{action["duration"].(int)})
				break
			case "http":
				action := value.(map[interface{}]interface{})
				var responseHandler HttpResponseHandler
				if action["response"] != nil {
					response := action["response"].(map[interface{}]interface{})
					responseHandler.Jsonpath = response["jsonpath"].(string)
					responseHandler.Variable = response["variable"].(string)
					responseHandler.Index = response["index"].(string)
				}
				actions = append(actions, HttpAction{
					action["method"].(string),
					action["url"].(string),
					getBody(action),
					action["accept"].(string),
					responseHandler,
				})

				break
			}
		}
	}
	return actions
}

func getBody(action map[interface{}]interface{}) string {
	var body string = ""
	if action["body"] != nil {
		body = action["body"].(string)
	}
	return body
}