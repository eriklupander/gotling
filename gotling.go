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

import(
    "io/ioutil"
    "os"
	"time"
    "fmt"
    "sync"
    //"github.com/davecheney/profile"
    "math/rand"
    "strconv"
    "gopkg.in/yaml.v2"
    "log"

)


var SimulationStart time.Time


func main() {

    spec := parseSpecFile()

   // defer profile.Start(profile.CPUProfile).Stop()

    // Start the web socket server, will not block exit until forced
    go StartWsServer()

    SimulationStart = time.Now()
    dir, _ := os.Getwd()
    dat, _ := ioutil.ReadFile(dir + "/" + spec)

    var t TestDef
    yaml.Unmarshal([]byte(dat), &t)

    if !ValidateTestDefinition(&t) {
        return
    }

    actions, isValid := buildActionList(&t)
    if !isValid {
        return
    }

    if t.Feeder.Type == "csv" {
        Csv(t.Feeder.Filename, ",")
    } else if t.Feeder.Type != "" {
        log.Fatal("Unsupported feeder type: " + t.Feeder.Type)
    }

    OpenResultsFile(dir + "/results/log/latest.log" )
	spawnUsers(&t, actions)

    fmt.Printf("Done in %v\n", time.Since(SimulationStart))
    fmt.Println("Building reports, please wait...")
    CloseResultsFile()
    //buildReport()
}

func parseSpecFile() (string) {
    if len(os.Args) == 1 {
        fmt.Errorf("No command line arguments, exiting...\n")
        panic("Cannot start simulation, no YAML simulaton specification supplied as command-line argument")
    }
    var s, sep string
    for i := 1; i < len(os.Args); i++ {
        s += sep + os.Args[i]
        sep = " "
    }
    if s == "" {
        panic(fmt.Sprintf("Specified simulation file '%s' is not a .yml file", s))
    }
    return s
}

func spawnUsers(t *TestDef, actions []Action) {
    resultsChannel := make(chan HttpReqResult, 10000) // buffer?
    go acceptResults(resultsChannel)
    wg := sync.WaitGroup{}
    for i := 0; i < t.Users; i++ {
        wg.Add(1)
        UID := strconv.Itoa(rand.Intn(t.Users+1) + 10000)
        go launchActions(t, resultsChannel, &wg, actions, UID)
        var waitDuration float32 = float32(t.Rampup) / float32(t.Users)
        time.Sleep( time.Duration( int(1000*waitDuration) )*time.Millisecond)
    }
    fmt.Println("All users started, waiting at WaitGroup")
    wg.Wait()
}

func launchActions(t *TestDef, resultsChannel chan HttpReqResult, wg *sync.WaitGroup, actions []Action, UID string) {
    var sessionMap = make(map[string]string)

    for i := 0; i < t.Iterations; i++ {

        // Make sure the sessionMap is cleared before each iteration - except for the UID which stays
        cleanSessionMapAndResetUID(UID, sessionMap)

        // If we have feeder data, pop an item and push its key-value pairs into the sessionMap
        feedSession(t, sessionMap)

        // Iterate over the actions. Note the use of the command-pattern like Execute method on the Action interface
        for _, action := range actions {
			if action != nil {
				action.(Action).Execute(resultsChannel, sessionMap)
			}
        }
    }
    wg.Done()
}

func cleanSessionMapAndResetUID(UID string, sessionMap map[string]string) {
    // Optimization? Delete all entries rather than reallocate map from scratch for each new iteration.
    for k := range sessionMap {
        delete(sessionMap, k)
    }
    sessionMap["UID"] = UID
}

func feedSession(t *TestDef, sessionMap map[string]string) {
    if t.Feeder.Type != "" {
        go NextFromFeeder() // Do async
        feedData := <- FeedChannel // Will block here until feeder delivers value over the FeedChannel
        for item := range feedData {
            sessionMap[item] = feedData[item]
        }
    }
}




