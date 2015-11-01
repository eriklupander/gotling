package main

import(
    "io/ioutil"
    "os"
	"time"
    "fmt"
    "sync"
    "gopkg.in/yaml.v2"

    "github.com/davecheney/profile"
	"github.com/eriklupander/goload/model"
)


var SimulationStart time.Time


func main() {

    defer profile.Start(profile.CPUProfile).Stop()

    // Start the web socket server, will not block exit until forced
    go StartWsServer()

    SimulationStart = time.Now()
    dir, _ := os.Getwd()
    dat, _ := ioutil.ReadFile(dir + "/samples/ltest00.yml")

    var t model.TestDef
    yaml.Unmarshal([]byte(dat), &t)

    if !ValidateTestDefinition(&t) {
        return
    }

    if t.Feeder.Type == "csv" {
         Csv(t.Feeder.Filename)
    }
    actions := buildActionList(&t)



    OpenResultsFile(dir + "/results/log/" + string(SimulationStart.UnixNano()) + ".log" )
	spawnUsers(&t, actions)

    fmt.Printf("Done in %v\n", time.Since(SimulationStart))
    fmt.Println("Building reports, please wait...")
    CloseResultsFile()
    //buildReport()
}

func spawnUsers(t *model.TestDef, actions []interface{}) {
    resultsChannel := make(chan model.HttpReqResult, 10000) // buffer?
    go acceptResults(resultsChannel)
    wg := sync.WaitGroup{}
    for i := 0; i < t.Users; i++ {
        wg.Add(1)
        go launchActions(t, resultsChannel, &wg, actions)
        var waitDuration float32 = float32(t.Rampup) / float32(t.Users)
        time.Sleep( time.Duration( int(1000*waitDuration) )*time.Millisecond)
    }
    fmt.Println("Waiting at WaitGroup")
    wg.Wait()
}

func launchActions(t *model.TestDef, resultsChannel chan model.HttpReqResult, wg *sync.WaitGroup, actions []interface{}) {
    var sessionMap = make(map[string]string)
    for i := 0; i < t.Iterations; i++ {

        // Optimization? Delete all entries rather than reallocate map from scratch for each new iteration.
        for k := range sessionMap {
            delete(sessionMap, k)
        }

        // If we have feeder data, pop an item and push its key-value pairs into the sessionMap
        if t.Feeder.Type != "" {
            go Next()  // FEL, får samma data två gånger, en per goroutine...
            feedData := <- Outchan
            for item := range feedData {
                sessionMap[item] = feedData[item]
            }
        }

        for _, action := range actions {

			if action != nil {
				action.(model.Action).Execute(resultsChannel, sessionMap)
			}
        }
    }
    wg.Done()
}






