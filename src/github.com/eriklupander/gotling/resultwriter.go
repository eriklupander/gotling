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
        panic(err)
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