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


func acceptResults(resChannel chan HttpReqResult) {
	for {
		select {
		case msg := <-resChannel:
			fmt.Println("received message", msg)
			statFrame := StatFrame {
                time.Since(SimulationStart).Nanoseconds(),
                msg.Latency,
            }
			BroadcastStatFrame(statFrame)
		default:
			//fmt.Println("no message received")
            time.Sleep(100 * time.Millisecond)
		}

	}
}

var SimulationStart time.Time

func main() {

    SimulationStart = time.Now()
    dir, _ := os.Getwd()
    dat, _ := ioutil.ReadFile(dir + "/samples/ltest00.yml")

    var t TestDef
    yaml.Unmarshal([]byte(dat), &t)

	resultsChannel := make(chan HttpReqResult, 1000) // buffer?
	go acceptResults(resultsChannel)

    for i := 0; i < t.Users; i++ {
        go runActions(t, resultsChannel)
    }
//    for _, element := range t.Actions {
//        //fmt.Printf("At index: %d with value: %v\n", index, element)
//        for key, value := range element {
//            //fmt.Println("Key:", key, "Value:", value)
//            switch key {
//            case "sleep":
//                a := value.(map[interface {}]interface {})
//                duration := a["duration"].(int)
//                time.Sleep(time.Duration(duration) * time.Second)
//                break
//            case "http":
//                a := value.(map[interface {}]interface{})
//
//
//				httpReq := HttpReqAction{a["method"].(string), a["url"].(string), a["accept"].(string)}
//				go DoHttpRequest(httpReq, resultsChannel)
//                break
//            }
//        }
//    }
    // Start the web socket server, will block exit until forced
    StartWsServer()
    //time.Sleep(5 * time.Second)
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

func runActions(t TestDef, resultsChannel chan HttpReqResult) {
    for i := 0; i < t.Iterations; i++ {
        for _, element := range t.Actions {
            //fmt.Printf("At index: %d with value: %v\n", index, element)
            for key, value := range element {
                //fmt.Println("Key:", key, "Value:", value)
                switch key {
                case "sleep":
                    a := value.(map[interface{}]interface{})
                    duration := a["duration"].(int)
                    time.Sleep(time.Duration(duration) * time.Second)
                    break
                case "http":
                    a := value.(map[interface{}]interface{})


                    httpReq := HttpReqAction{a["method"].(string), a["url"].(string), a["accept"].(string)}
                    go DoHttpRequest(httpReq, resultsChannel)
                    break
                }
            }
        }
    }
}



