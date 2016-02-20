package main

import (
	"net"
    "fmt"
    "time"
)

var conn net.Conn
// Accepts a TcpAction and a one-way channel to write the results to.
func DoTcpRequest(tcpAction TcpAction, resultsChannel chan HttpReqResult, sessionMap map[string]string) {

    address := SubstParams(sessionMap, tcpAction.Address)
    payload := SubstParams(sessionMap, tcpAction.Payload)

    if conn == nil {

        conn, err = net.Dial("tcp", address)
        if err != nil {
            fmt.Printf("TCP socket closed, error: %s\n", err)
            conn = nil
            return
        }
       // conn.SetDeadline(time.Now().Add(100 * time.Millisecond))
    }

    start := time.Now()

    _, err = fmt.Fprintf(conn, payload + "\r\n")
    if err != nil {
        fmt.Printf("TCP request failed with error: %s\n", err)
        conn = nil
    }

    elapsed := time.Since(start)
    resultsChannel <- buildTcpResult(0, 200, elapsed.Nanoseconds(), tcpAction.Title)

}

func buildTcpResult(contentLength int, status int, elapsed int64, title string) (HttpReqResult){
    httpReqResult := HttpReqResult {
        "TCP",
        elapsed,
        contentLength,
        status,
        title,
        time.Since(SimulationStart).Nanoseconds(),
    }
    return httpReqResult
}
