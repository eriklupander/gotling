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
//	"time"
	"fmt"
	"encoding/json"
//	"math/rand"
    "flag"
	"net/http"
    "log"
    "github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8182", "http service address")

var upgrader = websocket.Upgrader{} // use default options

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

func registerChannel(w http.ResponseWriter, r *http.Request) {
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

}

func StartWsServer() {
    fmt.Println("Starting WebSocket server")
    flag.Parse()
    log.SetFlags(0)

    http.HandleFunc("/start", registerChannel)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/" + r.URL.Path[1:])
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
	fmt.Println("Started WebSocket server")
}
