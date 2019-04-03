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