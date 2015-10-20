package main

import(
    "fmt"
    "reflect"
    "io/ioutil"
    "os"
    "gopkg.in/yaml.v2"
	"time"
)





func FillStruct(data map[string]interface{}, result interface{}) {
    t := reflect.ValueOf(result).Elem()
    for k, v := range data {
        val := t.FieldByName(k)
        val.Set(reflect.ValueOf(v))
    }
}


func acceptResults(resChannel chan StatFrame) {
	for {
		select {
		case msg := <-resChannel:
			fmt.Println("received message", msg)
			//
			BroadcastStatFrame(msg)
		default:
			fmt.Println("no message received")
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	simulationStart := time.Now()
    dir, _ := os.Getwd()
    dat, _ := ioutil.ReadFile(dir + "\\samples\\ltest00.yml")

    var t TestDef
    yaml.Unmarshal([]byte(dat), &t)

	resultsChannel := make(chan StatFrame, 1000) // no buffer
	go acceptResults(resultsChannel)

    for _, element := range t.Actions {
        //fmt.Printf("At index: %d with value: %v\n", index, element)
        for key, value := range element {
            //fmt.Println("Key:", key, "Value:", value)
            switch key {
            case "sleep":
                a := value.(map[interface {}]interface {})
                duration := a["duration"].(int)
                time.Sleep(time.Duration(duration) * time.Second)
                break
            case "http":
                a := value.(map[interface {}]interface{})


				httpReq := HttpReq{a["method"].(string), a["url"].(string), a["accept"].(string)}
				go DoHttpRequest(httpReq, resultsChannel, simulationStart)
                break
            }
        }
    }
	time.Sleep(5 * time.Second)
//
//    message := make(chan string) // no buffer
//    count := 3
//
//    go func() {
//        for i := 1; i <= count; i++ {
//            fmt.Println("send message")
//            message <- fmt.Sprintf("message %d", i)
//        }
//    }()
//
//    time.Sleep(time.Second * 3)
//
//    for i := 1; i <= count; i++ {
//        fmt.Println(<-message)
//    }
}



