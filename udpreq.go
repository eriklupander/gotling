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
	"net"
    "fmt"
    "time"
)

var udpconn *net.UDPConn

// Accepts a UdpAction and a one-way channel to write the results to.
func DoUdpRequest(udpAction UdpAction, resultsChannel chan HttpReqResult, sessionMap map[string]string) {

    address := SubstParams(sessionMap, udpAction.Address)
    payload := SubstParams(sessionMap, udpAction.Payload)

    if udpconn == nil {
        ServerAddr,err := net.ResolveUDPAddr("udp", address) //"127.0.0.1:10001")
        if err != nil {
            fmt.Println("Error ResolveUDPAddr remote: " + err.Error())
        }

        LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
        if err != nil {
            fmt.Println("Error ResolveUDPAddr local: " + err.Error())
        }

        udpconn, err = net.DialUDP("udp", LocalAddr, ServerAddr)
        if err != nil {
            fmt.Println("Error Dial: " + err.Error())
        }
    }
    //defer Conn.Close()
    start := time.Now()
    if udpconn != nil {
        _, err = fmt.Fprintf(udpconn, payload + "\r\n")
    }
    if err != nil {
        fmt.Printf("UDP request failed with error: %s\n", err)
        udpconn = nil
    }

    elapsed := time.Since(start)
    resultsChannel <- buildUdpResult(0, 200, elapsed.Nanoseconds(), udpAction.Title)

}

func buildUdpResult(contentLength int, status int, elapsed int64, title string) (HttpReqResult){
    httpReqResult := HttpReqResult {
        "UDP",
        elapsed,
        contentLength,
        status,
        title,
        time.Since(SimulationStart).Nanoseconds(),
    }
    return httpReqResult
}
