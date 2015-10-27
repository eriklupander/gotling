package main
import (
	"github.com/eriklupander/goload/model"
)

type HttpAction struct {
	Method string `yaml:"method"`
	Url string `yaml:"url"`
	Body string `yaml:"body"`
	Accept string `yaml:"accept"`
	Title string `yaml:"title"`
	ResponseHandler HttpResponseHandler `yaml:"response"`
}

func (h HttpAction) Execute(resultsChannel chan model.HttpReqResult, sessionMap map[string]string) {

	DoHttpRequest(h, resultsChannel, sessionMap)
}

type HttpResponseHandler struct {
	Jsonpath string `yaml:"jsonpath"`
	Variable string `yaml:"variable"`
	Index string `yaml:"index"`
}