package main

type TcpAction struct {
	Address string `yaml:"address"`
	Payload string `yaml:"payload"`
}

func (t TcpAction) Execute(resultsChannel chan HttpReqResult, sessionMap map[string]string) {
	DoTcpRequest(t, resultsChannel, sessionMap)
}

func NewTcpAction(a map[interface{}]interface{}) TcpAction {

	// TODO validation
	return TcpAction{
		a["address"].(string),
		a["payload"].(string),
	}
}