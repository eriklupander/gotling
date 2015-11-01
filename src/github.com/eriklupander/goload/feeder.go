package main
import (
    "os"
    "bufio"
    "fmt"
    "log"
    "strings"
)

var data []map[string]string
var index = 0

var Outchan chan map[string]string

func Next() {
    if data != nil && len(data) > 0 {
        Outchan <- data[index]

        // Cycle, does this need to be synchronized?
        if index < len(data) {
            index += 1
        } else {
            index = 0
        }
    }

}

func Csv(filename string) {
    dir, _ := os.Getwd()
    file, _ := os.Open(dir + "/data/" + filename)

    scanner := bufio.NewScanner(file)
    var lines int = 0

    data = make([]map[string]string, 0, 0)

    scanner.Scan()

    headers := strings.Split(scanner.Text(), ",")
    for scanner.Scan() {
        line := strings.Split(scanner.Text(), ",")
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
    Outchan = make(chan map[string]string)
}