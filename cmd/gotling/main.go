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
	"fmt"
	"github.com/eriklupander/gotling/internal/pkg/action"
	"github.com/eriklupander/gotling/internal/pkg/feeder"
	"github.com/eriklupander/gotling/internal/pkg/result"
	"github.com/eriklupander/gotling/internal/pkg/runtime"
	ws "github.com/eriklupander/gotling/internal/pkg/server"
	"github.com/eriklupander/gotling/internal/pkg/testdef"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
	//"github.com/davecheney/profile"
	"math/rand"
	"strconv"
)

func main() {

	spec := parseSpecFile()

	// defer profile.Start(profile.CPUProfile).Stop()

	// Start the web socket server, will not block exit until forced
	go ws.StartWsServer()

	runtime.SimulationStart = time.Now()
	dir, _ := os.Getwd()
	dat, _ := ioutil.ReadFile(dir + "/" + spec)

	var t testdef.TestDef
	err := yaml.Unmarshal([]byte(dat), &t)
	fail(err)

	if !testdef.ValidateTestDefinition(&t) {
		return
	}

	actions, isValid := action.BuildActionList(&t)
	if !isValid {
		return
	}

	if t.Feeder.Type == "csv" {
		feeder.Csv(t.Feeder.Filename, ",")
	} else if t.Feeder.Type != "" {
		log.Fatal("Unsupported feeder type: " + t.Feeder.Type)
	}

	result.OpenResultsFile(dir + "/results/log/latest.log")
	spawnUsers(&t, actions)

	fmt.Printf("Done in %v\n", time.Since(runtime.SimulationStart))
	fmt.Println("Building reports, please wait...")
	result.CloseResultsFile()
	//buildReport()
}

func parseSpecFile() string {
	if len(os.Args) == 1 {
		fmt.Println("No command line arguments, exiting...")
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

func spawnUsers(t *testdef.TestDef, actions []action.Action) {
	resultsChannel := make(chan result.HttpReqResult, 10000) // buffer?
	go result.AcceptResults(resultsChannel)
	wg := sync.WaitGroup{}
	for i := 0; i < t.Users; i++ {
		wg.Add(1)
		UID := strconv.Itoa(rand.Intn(t.Users+1) + 10000)
		go launchActions(t, resultsChannel, &wg, actions, UID)
		var waitDuration float32 = float32(t.Rampup) / float32(t.Users)
		time.Sleep(time.Duration(int(1000*waitDuration)) * time.Millisecond)
	}
	fmt.Println("All users started, waiting at WaitGroup")
	wg.Wait()
}

func launchActions(t *testdef.TestDef, resultsChannel chan result.HttpReqResult, wg *sync.WaitGroup, actions []action.Action, UID string) {
	var sessionMap = make(map[string]string)

	for i := 0; i < t.Iterations; i++ {

		// Make sure the sessionMap is cleared before each iteration - except for the UID which stays
		cleanSessionMapAndResetUID(UID, sessionMap)

		// If we have feeder data, pop an item and push its key-value pairs into the sessionMap
		feedSession(t, sessionMap)

		// Iterate over the actions. Note the use of the command-pattern like Execute method on the Action interface
		for _, action := range actions {
			if action != nil {
				action.Execute(resultsChannel, sessionMap)
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

func feedSession(t *testdef.TestDef, sessionMap map[string]string) {
	if t.Feeder.Type != "" {
		go feeder.NextFromFeeder()       // Do async
		feedData := <-feeder.FeedChannel // Will block here until feeder delivers value over the FeedChannel
		for item := range feedData {
			sessionMap[item] = feedData[item]
		}
	}
}

func fail(err error) {
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		os.Exit(1)
	}
}
