package model
import (
	"time"
)

type SleepAction struct {
	Duration int `yaml:"duration"`
}

func (s SleepAction) Execute(resultsChannel chan HttpReqResult, sessionMap map[string]string) {
	time.Sleep(time.Duration(s.Duration) * time.Second)
}

func NewSleepAction(a map[interface{}]interface{}) SleepAction {
	return SleepAction{a["duration"].(int)}
}