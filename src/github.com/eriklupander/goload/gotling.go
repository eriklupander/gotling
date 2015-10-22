package main

import(
    "reflect"
    "io/ioutil"
    "os"
    "gopkg.in/yaml.v2"
	"time"
    "fmt"
)





func FillStruct(data map[string]interface{}, result interface{}) {
    t := reflect.ValueOf(result).Elem()
    for k, v := range data {
        val := t.FieldByName(k)
        val.Set(reflect.ValueOf(v))
    }
}

/**
 * Loops indefinitely. The inner loop runs for exactly one second before submitting its
 * results to the WebSocket handler, then the aggregates are reset and restarted.
 */
func aggregatePerSecondHandler(perSecondChannel chan HttpReqResult) {

    for {

        var totalReq  int = 0
        var totalLatency int = 0
        until := time.Now().UnixNano() + 1000000000
        for time.Now().UnixNano() < until {
            select {
            case msg := <-perSecondChannel:
                totalReq++
                totalLatency += int(msg.Latency)
            default:

            }
        }
        // concurrently assemble the result and send it off to the websocket.
        go assembleAndSendResult(totalReq, totalLatency)
    }

}

func assembleAndSendResult(totalReq int, totalLatency int) {
    avgLatency := 0
    if totalReq > 0 {
        avgLatency = totalLatency / totalReq
    }
    fmt.Printf("\nBuilding stack frame for total latency: %d, latency: %d, totalReq: %d\n", totalLatency, avgLatency, totalReq)
    statFrame := StatFrame {
        time.Since(SimulationStart).Nanoseconds() / 100000000,
        avgLatency,
        totalReq,
    }
    BroadcastStatFrame(statFrame)
}

/**
 * Starts the per second aggregator and then forwards any HttpRequestResult messages to it through the channel.
 */
func acceptResults(resChannel chan HttpReqResult) {
    perSecondAggregatorChannel := make(chan HttpReqResult)
    go aggregatePerSecondHandler(perSecondAggregatorChannel)
	for {
		select {
            case msg := <-resChannel:
                perSecondAggregatorChannel <- msg
            default:
                time.Sleep(10 * time.Millisecond)
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

    // Start the web socket server, will block exit until forced
    StartWsServer()
}

func runActions(t TestDef, resultsChannel chan HttpReqResult) {
    for i := 0; i < t.Iterations; i++ {
        for _, element := range t.Actions {

            for key, value := range element {

                switch key {
                case "sleep":
                    a := value.(map[interface{}]interface{})

                    duration := a["duration"].(int)
                    time.Sleep(time.Duration(duration) * time.Second)
                    break
                case "http":
                    a := value.(map[interface{}]interface{})

                    httpReq := HttpReqAction{a["method"].(string), a["url"].(string), a["accept"].(string)}
                    DoHttpRequest(httpReq, resultsChannel)
                    break
                }
            }
        }
    }
}



