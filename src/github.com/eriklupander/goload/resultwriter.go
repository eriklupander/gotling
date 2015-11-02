package main
import (
    "os"
    "bufio"
    "github.com/eriklupander/goload/model"
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
}

func CloseResultsFile() {
    if opened {
        w.Flush()
        f.Close()
    }
    // Do nothing if not opened
}

func writeResult(httpResult *model.HttpReqResult) {
    jsonString, err := json.Marshal(httpResult)
    if err != nil {
        panic(err)
    }
    _, err = w.WriteString(string(jsonString))
    _, err = w.WriteString("\n")

    if err != nil {
        panic(err)
    }
    w.Flush()

}