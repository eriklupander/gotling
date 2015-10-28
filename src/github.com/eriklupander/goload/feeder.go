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

func Next() (map[string]string) {
    if data != nil && len(data) > 0 {

        // Cycle
        if index < len(data) {
            index += 1
        } else {
            index = 0
        }

        return data[index]
    }
    return nil
}

func Csv(filename string) {
    dir, _ := os.Getwd()
    file, _ := os.Open(dir + "/data/" + filename)

    scanner := bufio.NewScanner(file)
    var lines int = 0

    data = make([]map[string]string, 0, 10)

    scanner.Scan()

    headers := strings.Split(scanner.Text(), ",")
   // fmt.Printf("Headers: %v\n", headers)
    for scanner.Scan() {
        line := strings.Split(scanner.Text(), ",")
      //  fmt.Printf("Line: %v\n", line)
        item := make(map[string]string)
        for n := 0; n < len(headers); n++ {

            item[headers[n]] = line[n]
            data = append(data, item)
        }

        lines+=1
    }
    fmt.Printf("Data: %v\n", data)
    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    index = 0
    fmt.Printf("CSV feeder fed with %d lines of data\n", lines)
}