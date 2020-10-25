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
	"github.com/eriklupander/gotling/internal/pkg/result"
)

type UdpAction struct {
	Address string `yaml:"address"`
	Payload string `yaml:"payload"`
	Title   string `yaml:"title"`
}

func (t UdpAction) Execute(resultsChannel chan result.HttpReqResult, sessionMap map[string]string) {
	DoUdpRequest(t, resultsChannel, sessionMap)
}

func NewUdpAction(a map[interface{}]interface{}) UdpAction {
	payload, ok := a["payload"].(string)
	if !ok {
		return UdpAction{}
	}
	address, ok := a["address"].(string)
	if !ok {
		return UdpAction{}
	}
	title, ok := a["title"].(string)
	if !ok {
		return UdpAction{}
	}
	return UdpAction{
		address,
		payload,
		title,
	}
}

/*

 */
