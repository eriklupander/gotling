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
    "os"
    "bufio"
    "encoding/json"
)

var w *bufio.Writer
var f *os.File
var err error

var opened bool = false

func OpenResultsFile(fileName string) {
    if !opened {
        opened = true
    } else {
        return
    }
    f, err = os.Create(fileName)
    if err != nil {
    	os.Mkdir("results", 0777);
        os.Mkdir("results/log", 0777);
        f, err = os.Create(fileName)
        if err != nil {
            panic(err)
        }
    }
    w = bufio.NewWriter(f)
    _, err = w.WriteString(string("var logdata = '"))
}

func CloseResultsFile() {
    if opened {
        _, err = w.WriteString(string("';"))
        w.Flush()
        f.Close()
    }
    // Do nothing if not opened
}

func writeResult(httpResult *HttpReqResult) {
    jsonString, err := json.Marshal(httpResult)
    if err != nil {
        panic(err)
    }
    _, err = w.WriteString(string(jsonString))
    _, err = w.WriteString("|")

    if err != nil {
        panic(err)
    }
    w.Flush()

}