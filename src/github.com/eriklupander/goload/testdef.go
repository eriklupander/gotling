package main
import "fmt"

const FIRST = "first"
const LAST = "last"
const RANDOM = "random"

//type Action struct {
//	Type string `yaml:"type"`
//	Properties map[string]string `yaml:"properties"`
//}

type Action interface {
    Execute()
}

type SleepAction struct {
	Duration int `yaml:"duration"`
}

type HttpReqAction struct {
	Method string `yaml:"method"`
	Url string `yaml:"url"`
    Body string `yaml:"body"`
	Accept string `yaml:"accept"`
	ResponseHandler HttpResponseHandler `yaml:"response"`
}

func (h HttpReqAction) Execute() {
    fmt.Println("HttpReqAction")
}

func (s SleepAction) Execute() {
    fmt.Println("SleepAction")
}

var _ Action = (*HttpReqAction)(nil)
var _ Action = (*SleepAction)(nil)

type HttpResponseHandler struct {
    Jsonpath string `yaml:"jsonpath"`
    Variable string `yaml:"variable"`
    Index string `yaml:"index"`
}

type TestDef struct {
	Iterations int `yaml:"iterations"`
	Users int `yaml:"users"`
	Rampup int `yaml:"rampup"`
	Actions []map[string]interface{} `yaml:"actions"`
}
