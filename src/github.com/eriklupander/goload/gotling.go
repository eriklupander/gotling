package main

import(
    "io/ioutil"
    "os"
	"time"
    "fmt"
    "sync"
"strings"
    "regexp"
   // "reflect"
    "reflect"
    "gopkg.in/yaml.v2"
)


var SimulationStart time.Time

func main() {

    // Start the web socket server, will not block exit until forced
    go StartWsServer()

    SimulationStart = time.Now()
    dir, _ := os.Getwd()
    dat, _ := ioutil.ReadFile(dir + "/samples/ltest00.yml")

    var t TestDef
    yaml.Unmarshal([]byte(dat), &t)

    actions := buildActionList(&t)
	spawnUsers(&t, actions)

    fmt.Printf("Done in %v\n", time.Since(SimulationStart))

}

func spawnUsers(t *TestDef, actions []interface{}) {
    resultsChannel := make(chan HttpReqResult, 1000) // buffer?
    go acceptResults(resultsChannel)
    wg := sync.WaitGroup{}
    for i := 0; i < t.Users; i++ {
        wg.Add(1)
        //go runActions(&t, resultsChannel, &wg)
        go launchActions(t, resultsChannel, &wg, actions)
        var waitDuration float32 = float32(t.Rampup) / float32(t.Users)
        time.Sleep( time.Duration( int(1000*waitDuration) )*time.Millisecond)
    }
    fmt.Println("Waiting at WaitGroup")
    wg.Wait()
}

func launchActions(t *TestDef, resultsChannel chan HttpReqResult, wg *sync.WaitGroup, actions []interface{}) {

    for i := 0; i < t.Iterations; i++ {

        var sessionMap = make(map[string]string)

        for _, action := range actions {

            if action != nil {
                // Just to test..
                action.(Action).Execute()
            }
            // TODO introduce an "execute()" interface function as implicit interface. Let the execution code be
            // encapsulated by the Action OO-style.
            actionType := fmt.Sprintf("%s", reflect.TypeOf(action))
            switch actionType {
            case "main.HttpReqAction":
                httpReqAction := action.(HttpReqAction)
                DoHttpRequest(httpReqAction, resultsChannel, sessionMap)

                break
            case "main.SleepAction":
                sleepAction := action.(SleepAction)
                time.Sleep(time.Duration(sleepAction.Duration) * time.Second)
                break
            default:
                break
            }
        }
    }
    wg.Done()
}

func buildActionList(t *TestDef) []interface{} {
    actions := make([]interface{}, len(t.Actions), len(t.Actions))
    for _, element := range t.Actions {
        for key, value := range element {

            switch key {
            case "sleep":
                action := value.(map[interface{}]interface{})
                actions = append(actions, SleepAction{action["duration"].(int)})
                break
            case "http":
                action := value.(map[interface{}]interface{})
                var responseHandler HttpResponseHandler
                if action["response"] != nil {
                    response := action["response"].(map[interface{}]interface{})
                    responseHandler.Jsonpath = response["jsonpath"].(string)
                    responseHandler.Variable = response["variable"].(string)
                    responseHandler.Index = response["index"].(string)
                }
                actions = append(actions, HttpReqAction{
                    action["method"].(string),
                    action["url"].(string),
                    getBody(action),
                    action["accept"].(string),
                    responseHandler,
                })

                break
            }
        }
    }
    return actions
}

func getBody(action map[interface{}]interface{}) string {
    var body string = ""
    if action["body"] != nil {
        body = action["body"].(string)
    }
    return body
}

//func runActions(t *TestDef, resultsChannel chan HttpReqResult, wg *sync.WaitGroup) {
//    for i := 0; i < t.Iterations; i++ {
//        // TODO separate parsing of the YAML from actual execution, e.g. don't parse out stuff each time. Looks bad.
//
//        for _, element := range t.Actions {
//
//            for key, value := range element {
//
//                switch key {
//                case "sleep":
//                    action := value.(map[interface{}]interface{})
//
//                    duration := action["duration"].(int)
//                    time.Sleep(time.Duration(duration) * time.Second)
//                    break
//                case "http":
//                    action := value.(map[interface{}]interface{})
//                    var responseHandler HttpResponseHandler
//                    if action["response"] != nil {
//                        response := action["response"].(map[interface{}]interface{})
//                        responseHandler.Jsonpath = response["jsonpath"].(string)
//                        responseHandler.Variable = response["variable"].(string)
//                        responseHandler.Index = response["index"].(string)
//                    }
//                    url := strings.TrimSpace(substParams(sessionMap, action["url"].(string)))
//                    var body string = ""
//                    if action["body"] != nil {
//                        body = action["body"].(string)
//                    }
//                    httpReq := HttpReqAction{action["method"].(string), url, substParams(sessionMap, body), action["accept"].(string), responseHandler}
//                    DoHttpRequest(httpReq, resultsChannel, sessionMap)
//                    break
//                }
//            }
//        }
//    }
//    wg.Done()
//}



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
                // Can be trouble. Uses too much CPU if low, limits throughput if too high
                time.Sleep(10*time.Microsecond)
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
            // This is troublesome. If too high, throughput is bad. Too low, CPU use goes up too much
            time.Sleep(10 * time.Microsecond)
        }
    }
}



// Move this to some utility file
var re = regexp.MustCompile("\\$\\{([a-zA-Z0-9]{0,})\\}")

func SubstParams(sessionMap map[string]string, textData string) string {
    if strings.ContainsAny(textData, "${") {
        res := re.FindAllStringSubmatch(textData, -1)
        for _, v := range res {
            textData = strings.Replace(textData, "${" + v[1] + "}", sessionMap[v[1]], 1)
        }
        return textData
    } else {
        return textData
    }
    return textData
}
