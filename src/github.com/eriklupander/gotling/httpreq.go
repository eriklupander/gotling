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
    "io/ioutil"
    "log"
    "math/rand"
    "net/http"
    "strings"
    "time"
    //"fmt"
    "gopkg.in/xmlpath.v2"
    "github.com/NodePrime/jsonpath"
    "bytes"
    "crypto/tls"

)

// Accepts a Httpaction and a one-way channel to write the results to.
func DoHttpRequest(httpAction HttpAction, resultsChannel chan HttpReqResult, sessionMap map[string]string) {
    req := buildHttpRequest(httpAction, sessionMap)

    start := time.Now()
    var DefaultTransport http.RoundTripper = &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    resp, err := DefaultTransport.RoundTrip(req)

    if err != nil {
        log.Printf("HTTP request failed: %s", err)
    } else {
        elapsed := time.Since(start)
        responseBody, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            //log.Fatal(err)
            log.Printf("Reading HTTP response failed: %s\n", err)
            httpReqResult := buildHttpResult(0, resp.StatusCode, elapsed.Nanoseconds(), httpAction.Title)

            resultsChannel <- httpReqResult
        } else {
            defer resp.Body.Close()


            if httpAction.StoreCookie != "" {
                for _, cookie := range resp.Cookies() {

                    if cookie.Name == httpAction.StoreCookie {
                        sessionMap["____" + cookie.Name] = cookie.Value
                    }
                }
            }

            // if action specifies response action, parse using regexp/jsonpath
            processResult(httpAction, sessionMap, responseBody)

            httpReqResult := buildHttpResult(len(responseBody), resp.StatusCode, elapsed.Nanoseconds(), httpAction.Title)

            resultsChannel <- httpReqResult
        }
    }
}

func buildHttpResult(contentLength int, status int, elapsed int64, title string) HttpReqResult {
    httpReqResult := HttpReqResult{
        "HTTP",
        elapsed,
        contentLength,
        status,
        title,
        time.Since(SimulationStart).Nanoseconds(),
    }
    return httpReqResult
}

func buildHttpRequest(httpAction HttpAction, sessionMap map[string]string) *http.Request {
    var req *http.Request
    var err error
    if httpAction.Body != "" {
        reader := strings.NewReader(SubstParams(sessionMap, httpAction.Body))
        req, err = http.NewRequest(httpAction.Method, SubstParams(sessionMap, httpAction.Url), reader)
    } else if httpAction.Template != "" {
        reader := strings.NewReader(SubstParams(sessionMap, httpAction.Template))
        req, err = http.NewRequest(httpAction.Method, SubstParams(sessionMap, httpAction.Url), reader)
    } else {
        req, err = http.NewRequest(httpAction.Method, SubstParams(sessionMap, httpAction.Url), nil)
    }
    if err != nil {
        log.Fatal(err)
    }

    // Add headers
    req.Header.Add("Accept", httpAction.Accept)
    if (httpAction.ContentType != "") {
        req.Header.Add("Content-Type", httpAction.ContentType)
    }

    // Add cookies stored by subsequent requests in the sessionMap having the kludgy ____ prefix
    for key, value := range sessionMap {
        if strings.HasPrefix(key, "____") {

            cookie := http.Cookie{
                Name: key[4:len(key)],
                Value: value,
            }

            req.AddCookie(&cookie)
        }
    }

    return req
}

/**
 * If the httpAction specifies a Jsonpath in the Response, try to extract value(s)
 * from the responseBody.
 *
 * TODO extract both Jsonpath handling and Xmlpath handling into separate functions, and write tests for them.
 *
 * Uses github.com/NodePrime/jsonpath
 */
func processResult(httpAction HttpAction, sessionMap map[string]string, responseBody []byte) {
    if httpAction.ResponseHandler.Jsonpath != "" {
        paths, err := jsonpath.ParsePaths(httpAction.ResponseHandler.Jsonpath)
        if err != nil {
            panic(err)
        }
        eval, err := jsonpath.EvalPathsInBytes(responseBody, paths)
        if err != nil {
            panic(err)
        }

        // TODO optimization: Don't reinitialize each time, reuse this somehow.
        resultsArray := make([]string, 0, 10)
        for {
            if result, ok := eval.Next(); ok {

                value := strings.TrimSpace(result.Pretty(false))
                resultsArray = append(resultsArray, trimChar(value, '"'))
            } else {
                break
            }
        }
        if eval.Error != nil {
            panic(eval.Error)
        }

        passResultIntoSessionMap(resultsArray, httpAction, sessionMap)
    }


    if httpAction.ResponseHandler.Xmlpath != "" {
        path := xmlpath.MustCompile(httpAction.ResponseHandler.Xmlpath)
        r := bytes.NewReader(responseBody)
        root, err := xmlpath.Parse(r)

        if err != nil {
            log.Fatal(err)
        }

        iterator := path.Iter(root)
        hasNext := iterator.Next()
        if hasNext {
            resultsArray := make([]string, 0, 10)
            for {
                if hasNext {
                    node := iterator.Node()
                    resultsArray = append(resultsArray, node.String())
                    hasNext = iterator.Next()
                } else {
                    break
                }
            }
            passResultIntoSessionMap(resultsArray, httpAction, sessionMap)
        }
    }

    // log.Println(string(responseBody))
}

/**
 * Trims leading and trailing byte r from string s
 */
func trimChar(s string, r byte) string {
    sz := len(s)

    if sz > 0 && s[sz-1] == r {
        s = s[:sz-1]
    }
    sz = len(s)
    if sz > 0 && s[0] == r {
        s = s[1:sz]
    }
    return s
}



func passResultIntoSessionMap(resultsArray []string, httpAction HttpAction, sessionMap map[string]string) {
    resultCount := len(resultsArray)

    if resultCount > 0 {
        switch httpAction.ResponseHandler.Index {
        case FIRST:
            sessionMap[httpAction.ResponseHandler.Variable] = resultsArray[0]
            break
        case LAST:
            sessionMap[httpAction.ResponseHandler.Variable] = resultsArray[resultCount-1]
            break
        case RANDOM:
            if resultCount > 1 {
                sessionMap[httpAction.ResponseHandler.Variable] = resultsArray[rand.Intn(resultCount-1)]
            } else {
                sessionMap[httpAction.ResponseHandler.Variable] = resultsArray[0]
            }
            break
        }

    } else {
        // TODO how to handle requested, but missing result?
    }
}
