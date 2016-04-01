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
	"time"
	"fmt"
)

/**
 * Loops indefinitely. The inner loop runs for exactly one second before submitting its
 * results to the WebSocket handler, then the aggregates are reset and restarted.
 */
func aggregatePerSecondHandler(perSecondChannel chan *HttpReqResult) {

	for {

		var totalReq  int = 0
		var totalLatency int = 0
		until := time.Now().UnixNano() + 1000000000
		for time.Now().UnixNano() < until {
			select {
			case msg := <-perSecondChannel:
				totalReq++
				totalLatency += int(msg.Latency/1000) // measure in microseconds
			default:
			// Can be trouble. Uses too much CPU if low, limits throughput if too high
				time.Sleep(100*time.Microsecond)
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
		time.Since(SimulationStart).Nanoseconds() / 1000000000, // seconds
		avgLatency,                                             // microseconds
		totalReq,
	}
	fmt.Printf("Time: %d Avg latency: %d μs (%d ms) req/s: %d\n", statFrame.Time, statFrame.Latency, statFrame.Latency / 1000, statFrame.Reqs)
	BroadcastStatFrame(statFrame)
}

/**
 * Starts the per second aggregator and then forwards any HttpRequestResult messages to it through the channel.
 */
func acceptResults(resChannel chan HttpReqResult) {
	perSecondAggregatorChannel := make(chan *HttpReqResult, 5)
	go aggregatePerSecondHandler(perSecondAggregatorChannel)
	for {
		select {
		case msg := <-resChannel:
			perSecondAggregatorChannel <- &msg
			writeResult(&msg) // sync write result to file for later processing.
            break
		case <-	time.After(100 * time.Microsecond):
            break
//		default:
//			// This is troublesome. If too high, throughput is bad. Too low, CPU use goes up too much
//			// Using a sync channel kills performance
//			time.Sleep(100 * time.Microsecond)
		}
	}
}
