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
package action

import (
	"fmt"
	"github.com/eriklupander/gotling/internal/pkg/result"
	"time"
)

type SleepAction struct {
	Duration time.Duration `yaml:"duration"`
}

func (s SleepAction) Execute(resultsChannel chan result.HttpReqResult, sessionMap map[string]string) {
	time.Sleep(s.Duration)
}

func NewSleepAction(a map[interface{}]interface{}) SleepAction {
	switch val := a["duration"].(type) {
	case int:
		return SleepAction{Duration: time.Second * time.Duration(val)}
	case string:
		dur, err := time.ParseDuration(val)
		if err != nil {
			fmt.Printf("Error trying to parse duration '%v' from string representation into Go duration format. Error: %v\n", val, err.Error())
			panic(err.Error())
		}
		return SleepAction{Duration:dur}
	case time.Duration:
		return SleepAction{Duration: val}
	default:
		fmt.Printf("unsupported Sleep value type. Supported is int or string (golang time.Duration), was %T\n", val)
		panic("unsupported sleep value")
	}
}
