package main
import (
	"log"
)

type HttpAction struct {
	Method string `yaml:"method"`
	Url string `yaml:"url"`
	Body string `yaml:"body"`
	Accept string `yaml:"accept"`
	Title string `yaml:"title"`
	ResponseHandler HttpResponseHandler `yaml:"response"`
}

func (h HttpAction) Execute(resultsChannel chan HttpReqResult, sessionMap map[string]string) {
	DoHttpRequest(h, resultsChannel, sessionMap)
}

type HttpResponseHandler struct {
	Jsonpath string `yaml:"jsonpath"`
	Variable string `yaml:"variable"`
	Index string `yaml:"index"`
}

func NewHttpAction(a map[interface{}]interface{}) HttpAction {
    var valid bool = true
    if a["accept"] == nil || a["accept"] != "json" {
        log.Println("Error: HttpAction only accepts 'json' as Accept")
        valid = false
    }
    if a["url"] == "" || a["url"] == nil {
        log.Println("Error: HttpAction must define a URL")
        valid = false
    }
    if a["method"] == nil || (a["method"] != "GET" && a["method"] != "POST" && a["method"] != "PUT" && a["method"] != "DELETE")  {
        log.Println("Error: HttpAction must specify a HTTP method: GET, POST, PUT or DELETE")
        valid = false
    }
    if a["title"] == nil || a["title"] == "" {
        log.Println("Error: HttpAction must define a title")
        valid = false
    }

    if a["response"] != nil {
        r := a["response"].(map[interface{}]interface{})
        if r["index"] == nil || r["index"] == "" || (r["index"] != "first" && r["index"] != "last" && r["index"] != "random") {
            log.Println("Error: HttpAction ResponseHandler must define an Index of either of: first, last or random")
            valid = false
        }
        if r["jsonpath"] == nil || r["jsonpath"] == "" {
            log.Println("Error: HttpAction ResponseHandler must define a Jsonpath")
            valid = false
        }
        if r["variable"] == nil ||  r["variable"] == "" {
            log.Println("Error: HttpAction ResponseHandler must define a Variable")
            valid = false
        }
    }

    if !valid {
        log.Fatalf("Your YAML defintion contains an invalid HttpAction, see errors listed above.")
    }
    var responseHandler HttpResponseHandler
    if a["response"] != nil {
        response := a["response"].(map[interface{}]interface{})
        responseHandler.Jsonpath = response["jsonpath"].(string)
        responseHandler.Variable = response["variable"].(string)
        responseHandler.Index = response["index"].(string)
    }
    httpAction := HttpAction{
        a["method"].(string),
        a["url"].(string),
        getBody(a),
        a["accept"].(string),
        a["title"].(string),
        responseHandler}

    return httpAction

}