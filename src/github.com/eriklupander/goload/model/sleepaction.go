package model
import (
	"fmt"
	"time"
)

type SleepAction struct {
	Duration int `yaml:"duration"`
}

func (s SleepAction) Execute(resultsChannel chan HttpReqResult, sessionMap map[string]string) {
	fmt.Println("SleepAction")
	time.Sleep(time.Duration(s.Duration) * time.Second)
}