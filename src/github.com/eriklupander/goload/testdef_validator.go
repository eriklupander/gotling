package main
import (
    "log"
    "github.com/eriklupander/goload/model"
    "reflect"
)

// TODO refactor this so it runs before parsing the actions

func ValidateTestDefinition(t *model.TestDef) (bool) {
    var valid bool = true
    if t.Iterations == 0 {
        log.Fatalf("Iterations not set, must be > 0")
        valid = false
    }
    if t.Rampup < 0 {
        log.Fatalf("Rampup not defined. must be > -1")
        valid = false
    }
    if t.Users == 0 {
        log.Fatalf("Users must be > 0")
        valid = false
    }
    return valid
}

func ValidateActions(actions []interface{}) (bool) {
    var valid bool = true
    for _, action := range actions {

        if action != nil {
            actionType := reflect.TypeOf(action).String()
            switch actionType {
            case "model.SleepAction":
                sleepAction := action.(model.SleepAction)
                if sleepAction.Duration == 0 {
                    log.Fatalf("Sleep action must define Duration > 0 in seconds")
                    valid = false
                }
                break
            case "main.HttpAction":
                httpAction := action.(HttpAction)
                if httpAction.Accept != "json" {
                    log.Fatalf("HttpAction only accepts 'json' as Accept")
                    valid = false
                }
                if httpAction.Url == "" {
                    log.Fatalf("HttpAction must define a URL")
                    valid = false
                }
                if httpAction.Method != "GET" && httpAction.Method != "POST" && httpAction.Method != "PUT" && httpAction.Method != "DELETE"  {
                    log.Fatalf("HttpAction must specify a HTTP method: GET, POST, PUT or DELETE")
                    valid = false
                }
                if httpAction.Title == "" {
                    log.Fatalf("HttpAction must define a title")
                    valid = false
                }

                if httpAction.ResponseHandler.Index == "" {
                    log.Fatalf("HttpAction ResponseHandler must define an Index of either of: first, last or random")
                    valid = false
                }
                if httpAction.ResponseHandler.Jsonpath == "" {
                    log.Fatalf("HttpAction ResponseHandler must define a Jsonpath")
                    valid = false
                }
                if httpAction.ResponseHandler.Variable == "" {
                    log.Fatalf("HttpAction ResponseHandler must define a Variable")
                    valid = false
                }
                break
            }
        }

    }

    return valid
}