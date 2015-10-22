package main
import (
//	"time"
	"fmt"
	"encoding/json"
//	"math/rand"
//	"golang.org/x/net/websocket"
    "flag"
	"net/http"
    "log"
    "github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8182", "http service address")

var upgrader = websocket.Upgrader{} // use default options

//func startDataSimulation(ws *websocket.Conn) {

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
//}
//
//func registerChannel(ws *websocket.Conn) {
//	wsChannels = append(wsChannels, ws)
//	fmt.Printf("Added Web Socket channel to registry, size is now %d connections", len(wsChannels))
//}

func Remove(item int) {
    connectionRegistry = append(connectionRegistry[:item], connectionRegistry[item+1:]...)
}

func BroadcastStatFrame(statFrame StatFrame) {
    for index, wsConn := range connectionRegistry {
		serializedFrame, _  := json.Marshal(statFrame)
		err := wsConn.WriteMessage(1, serializedFrame)
		if err != nil {
			// Detected disconnected channel. Need to clean up.

            fmt.Printf("Could not write to channel: %v", err)
            wsConn.Close()
            Remove(index)
		}
	}

}

var connectionRegistry = make([]*websocket.Conn, 0, 10)

func echo(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/start" {
        http.Error(w, "Not found", 404)
        return
    }
    if r.Method != "GET" {
        http.Error(w, "Method not allowed", 405)
        return
    }
    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Print("upgrade:", err)
        return
    }
    connectionRegistry = append(connectionRegistry, c)
    /**
    defer c.Close()
    for {
        mt, message, err := c.ReadMessage()
        if err != nil {
            log.Println("read:", err)
            break
        }
        log.Printf("recv: %s", message)
        err = c.WriteMessage(mt, message)
        if err != nil {
            log.Println("write:", err)
            break
        }
    }
    */
}

func StartWsServer() {


    fmt.Println("Starting WebSocket server")
    flag.Parse()
    log.SetFlags(0)

	//http.Handle("/start", websocket.Handler(registerChannel))
    http.HandleFunc("/start", echo)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/" + r.URL.Path[1:])
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
	fmt.Println("Started WebSocket server")
}
