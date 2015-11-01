package main
import (
	"github.com/eriklupander/goload/model"
)

type TcpAction struct {
	Address string `yaml:"address"`
	Payload string `yaml:"payload"`
}

func (t TcpAction) Execute(resultsChannel chan model.HttpReqResult, sessionMap map[string]string) {
	DoTcpRequest(t, resultsChannel, sessionMap)
}

func NewTcpAction(a map[interface{}]interface{}) TcpAction {

	// TODO validation
	return TcpAction{
		a["address"].(string),
		a["payload"].(string),
	}
}