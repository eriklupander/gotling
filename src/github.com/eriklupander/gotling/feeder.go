package main
import (
    "os"
    "bufio"
    "fmt"
    "log"
    "strings"
    "sync"
)

var data []map[string]string
var index = 0

var l sync.Mutex

// "public" synchronized channel for delivering feeder data
var FeedChannel chan map[string]string

//
func NextFromFeeder() {

    if data != nil && len(data) > 0 {

        // Push data into the FeedChannel
       // fmt.Printf("Current index: %d of total size: %d\n", index, len(data))
        FeedChannel <- data[index]

        // Cycle, does this need to be synchronized?
        l.Lock()
        if index < len(data) - 1 {
            index += 1
        } else {
            index = 0
        }
        l.Unlock()
    }

}

func Csv(filename string, separator string) {
    dir, _ := os.Getwd()
    file, _ := os.Open(dir + "/data/" + filename)

    scanner := bufio.NewScanner(file)
    var lines int = 0

    data = make([]map[string]string, 0, 0)

    // Scan the first line, should contain headers.
    scanner.Scan()
    headers := strings.Split(scanner.Text(), separator)

    for scanner.Scan() {
        line := strings.Split(scanner.Text(), separator)
        item := make(map[string]string)
        for n := 0; n < len(headers); n++ {
            item[headers[n]] = line[n]
        }
        data = append(data, item)
        lines+=1
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    index = 0
    fmt.Printf("CSV feeder fed with %d lines of data\n", lines)
    FeedChannel = make(chan map[string]string)
}