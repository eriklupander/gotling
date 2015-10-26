package main
import (
	"time"
	"github.com/eriklupander/goload/model"
)

/**
 * Loops indefinitely. The inner loop runs for exactly one second before submitting its
 * results to the WebSocket handler, then the aggregates are reset and restarted.
 */
func aggregatePerSecondHandler(perSecondChannel chan model.HttpReqResult) {

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
		time.Since(SimulationStart).Nanoseconds() / 100000000,
		avgLatency,
		totalReq,
	}
	BroadcastStatFrame(statFrame)
}

/**
 * Starts the per second aggregator and then forwards any HttpRequestResult messages to it through the channel.
 */
func acceptResults(resChannel chan model.HttpReqResult) {
	perSecondAggregatorChannel := make(chan model.HttpReqResult, 5)
	go aggregatePerSecondHandler(perSecondAggregatorChannel)
	for {
		select {
		case msg := <-resChannel:
			perSecondAggregatorChannel <- msg
		default:
		// This is troublesome. If too high, throughput is bad. Too low, CPU use goes up too much
			time.Sleep(100 * time.Microsecond)
		}
	}
}
