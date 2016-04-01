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
)

type HttpAction struct {
    Method          string              `yaml:"method"`
    Url             string              `yaml:"url"`
    Body            string              `yaml:"body"`
    Template        string              `yaml:"template"`
    Accept          string              `yaml:"accept"`
    ContentType     string              `yaml:"contentType"`
    Title           string              `yaml:"title"`
    ResponseHandler HttpResponseHandler `yaml:"response"`
    StoreCookie     string              `yaml:"storeCookie"`
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
        log.Println("Error: HttpAction must define a URL.")
        valid = false
    }
    if a["method"] == nil || (a["method"] != "GET" && a["method"] != "POST" && a["method"] != "PUT" && a["method"] != "DELETE") {
        log.Println("Error: HttpAction must specify a HTTP method: GET, POST, PUT or DELETE")
        valid = false
    }
    if a["title"] == nil || a["title"] == "" {
        log.Println("Error: HttpAction must define a title.")
        valid = false
    }

    if a["body"] != nil && a["template"] != nil {
        log.Println("Error: A HttpAction can not define both a 'body' and a 'template'.")
        valid = false
    }

    if a["response"] != nil {
        r := a["response"].(map[interface{}]interface{})
        if r["index"] == nil || r["index"] == "" || (r["index"] != "first" && r["index"] != "last" && r["index"] != "random") {
            log.Println("Error: HttpAction ResponseHandler must define an Index of either of: first, last or random.")
            valid = false
        }
        if (r["jsonpath"] == nil || r["jsonpath"] == "") && (r["xmlpath"] == nil || r["xmlpath"] == "") {
            log.Println("Error: HttpAction ResponseHandler must define a Jsonpath or a Xmlpath.")
            valid = false
        }
        if (r["jsonpath"] != nil && r["jsonpath"] != "") && (r["xmlpath"] != nil && r["xmlpath"] != "") {
            log.Println("Error: HttpAction ResponseHandler can only define either a Jsonpath OR a Xmlpath.")
            valid = false
        }

        // TODO perhaps compile Xmlpath expressions so we can validate early?

        if r["variable"] == nil || r["variable"] == "" {
            log.Println("Error: HttpAction ResponseHandler must define a Variable.")
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

    var contentType string
    if a["contentType"] != nil && len(a["contentType"].(string)) > 0 {
        contentType = a["contentType"].(string)
    }

    var storeCookie string
    if a["storeCookie"] != nil && a["storeCookie"].(string) != "" {
        storeCookie = a["storeCookie"].(string)
    }

    httpAction := HttpAction{
        a["method"].(string),
        a["url"].(string),
        getBody(a),
        getTemplate(a),
        accept,
        contentType,
        a["title"].(string),
        responseHandler,
        storeCookie,
    }

    return httpAction
}
