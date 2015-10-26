package model

type Action interface {
	Execute(resultsChannel chan HttpReqResult, sessionMap map[string]string)
}