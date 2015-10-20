package main
import (
//	"time"
	"fmt"
	"encoding/json"
//	"math/rand"
	"golang.org/x/net/websocket"
	"net/http"
)

var wsChannels []*websocket.Conn

func startDataSimulation(ws *websocket.Conn) {

//	for i := 0; i < 1000; i++ {
//		var statFrame StatFrame
//		statFrame.Time = int64(i*1000)
//		statFrame.ReqS = (rand.Intn(20)+150)
//		serializedFrame, err  := json.Marshal(statFrame)
//		if err != nil {
//			panic(err)
//		}
//		ws.Write(serializedFrame)
//		fmt.Printf("sent StatFrame: %s\n", serializedFrame)
//		time.Sleep(1000 * time.Millisecond)
//	}
}

func registerChannel(ws *websocket.Conn) {
	wsChannels = append(wsChannels, ws)
	fmt.Printf("Added Web Socket channel to registry, size is now %d connections", len(wsChannels))
}

func BroadcastStatFrame(statFrame StatFrame) {
	for _, ws := range wsChannels {
		serializedFrame, _  := json.Marshal(statFrame)
		ws.Write(serializedFrame)
	}

}

func StartWsServer() {
	fmt.Println("Starting WebSocket server")
	http.Handle("/start", websocket.Handler(registerChannel))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/" + r.URL.Path[1:])
	})
	err := http.ListenAndServe(":8182", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
	fmt.Println("Started WebSocket server")
}
