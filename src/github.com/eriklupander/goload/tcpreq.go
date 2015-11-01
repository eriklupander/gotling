package main

import (
	"github.com/eriklupander/goload/model"
    "net"
    "fmt"
)

var conn net.Conn
// Accepts a TcpAction and a one-way channel to write the results to.
func DoTcpRequest(tcpAction TcpAction, resultsChannel chan model.HttpReqResult, sessionMap map[string]string) {

    address := SubstParams(sessionMap, tcpAction.Address)
    payload := SubstParams(sessionMap, tcpAction.Payload)

    if conn == nil {
        conn, _ = net.Dial("tcp", address)

    }
    fmt.Fprintf(conn, payload + "\r\n")

}

