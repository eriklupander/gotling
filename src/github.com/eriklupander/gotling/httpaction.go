package main

import (
    "log"
)

type HttpAction struct {
    Method          string              `yaml:"method"`
    Url             string              `yaml:"url"`
    Body            string              `yaml:"body"`
    Accept          string              `yaml:"accept"`
    Title           string              `yaml:"title"`
    ResponseHandler HttpResponseHandler `yaml:"response"`
}

func (h HttpAction) Execute(resultsChannel chan HttpReqResult, sessionMap map[string]string) {
    DoHttpRequest(h, resultsChannel, sessionMap)
}

type HttpResponseHandler struct {
    Jsonpath string `yaml:"jsonpath"`
    Xmlpath string `yaml:"xmlpath"`
    Variable string `yaml:"variable"`
    Index    string `yaml:"index"`
}

func NewHttpAction(a map[interface{}]interface{}) HttpAction {
    var valid bool = true
    if a["url"] == "" || a["url"] == nil {
        log.Println("Error: HttpAction must define a URL")
        valid = false
    }
    if a["method"] == nil || (a["method"] != "GET" && a["method"] != "POST" && a["method"] != "PUT" && a["method"] != "DELETE") {
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
        if (r["jsonpath"] == nil || r["jsonpath"] == "") && (r["xmlpath"] == nil || r["xmlpath"] == "") {
            log.Println("Error: HttpAction ResponseHandler must define a Jsonpath or a Xmlpath")
            valid = false
        }
        if (r["jsonpath"] != nil && r["jsonpath"] != "") && (r["xmlpath"] != nil && r["xmlpath"] != "") {
            log.Println("Error: HttpAction ResponseHandler can only define either a Jsonpath OR a Xmlpath")
            valid = false
        }

        // TODO perhaps compile Xmlpath expressions so we can validate early?

        if r["variable"] == nil || r["variable"] == "" {
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

        if response["jsonpath"] != nil && response["jsonpath"] != "" {
            responseHandler.Jsonpath = response["jsonpath"].(string)
        }
        if response["xmlpath"] != nil && response["xmlpath"] != "" {
            responseHandler.Xmlpath = response["xmlpath"].(string)
        }

        responseHandler.Variable = response["variable"].(string)
        responseHandler.Index = response["index"].(string)
    }

    accept := "text/html,application/json,application/xhtml+xml,application/xml,text/plain"
    if a["accept"] != nil && len(a["accept"].(string)) > 0 {
        accept = a["accept"].(string)
    }

    httpAction := HttpAction{
        a["method"].(string),
        a["url"].(string),
        getBody(a),
        accept,
        a["title"].(string),
        responseHandler}

    return httpAction
}
