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
	fmt.Printf("Time: %d Avg latency: %d μs req/s: %d\n", statFrame.Time, statFrame.Latency, statFrame.Reqs)
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
