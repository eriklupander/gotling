# gotling
Simple golang-based load test application using YAML documents as specification.

For a more full-blown explanation of what Gotling is about, see my blog post here: http://callistaenterprise.se/blogg/teknik/2015/11/22/gotling/

_Please note that this is my very first golang program and is probably full of anti-patterns and bad use of golang constructs._

## What it does
- Provides high-throughput load testing of HTTP services
    - Supports GET, POST, PUT and DELETE
    - Live metrics using HTML5 canvas from canvasjs
    - Request URLs and bodies can contain ${paramName} parameters
    - ${paramName} values can be extracted from HTTP response bodies and bound to a User context
- TCP sockets
    - Can send line-terminated text strings
    - Possible to use ${varname} variables in payload
    - Does not currently support body parsing or variable extraction from data incoming on the TCP socket.

## Building

Of course you need the Go SDK, I think I've used version 1.5 or so. These instructions are based on OS X 10.11, but should apply to Windows and Linux too.

#### 1. Clone the source from github
    git clone git@github.com:eriklupander/gotling.git
    
#### 2. Open a command shell 
Go into the root project directory, for example ~/projects/gotling.

#### 3. Set GOPATH

    export GOPATH=~/projects/gotling
    
#### 4. Use go get to fethc deps

    cd src/github.com/eriklupander/gotling
    go get
    
May take a little while, after go get finishes, you should see 

    src/github.com/NodePrime/
    src/github.com/eriklupander/
    src/github.com/gorilla/
    src/github.com/tobyhede/
    src/gopkg.in/yaml.v2/
    
#### 5. Run

To run the standard bundled "demo" simulation, use go run like this:

    cd ~/projects/gotling
    go run src/github.com/eriklupander/gotling/*.go samples/demosimulation.yml
    


## Usage
Define your test setup in a .yml file

    ---
    iterations: 10          # Number of iterations per user
    users: 10               # Number of users
    rampup: 20              # seconds
    actions:                # List of actions. Currently supports http, sleep
      - http:
          method: GET                                  # GET, POST, PUT, DELETE
          url: http://localhost:8183/courses           # URL. Can include ${paramName} parameters
          accept: json                                 # Only 'json' is currently supported
          response:                                    # Defines response handling
            jsonpath: $[*].id+                         # JsonPath expression to capture data from JSON response
            variable: courseId                         # Parameter name for captured value
            index: first                               # If > 1 results found - which result to store in 'variable': 
                                                       # first, random, last
      - sleep:
          duration: 3                                  # Sleep duration in seconds. Will block current user
      - http:
          method: GET
          url: http://localhost:8183/courses/${courseId}
          accept: json
          response:
            jsonpath: $.author+
            variable: author
            index: first
      - sleep:
            duration: 3
      - tcp:
            address: 127.0.0.1:8081                     # TCP socket connection
            payload: |TYPE|1|${UID}|${email}            # Sample data using parameter substitution

### Feeders and user context
Gotling currently support CSV feeds of data. First line needs to be comma-separated headers with the following lines containing data, e.g:

    id,name,size
    1,Bear,Large
    2,Cat,Small
    3,Deer,Medium
    
These values can be accessed through ${varname} matching the column header.

#### The UID
Each "user" gets a unique "UID" assigned to it, typically an integer from 10000 + random(no of users). Perhaps I can tweak this to either use UUID's or configurable intervals. Anyway, the UID can be used using ${UID} and can be useful for grouping data etc.



## Realtime dashboard
Access at http://localhost:8182

Click "connect" to connect to the currently executing test.

![Gotling dashboard](gotling-dashboard.png)

## Uses the following libraries
- gopkg.in/yaml.v2
- NodePrime/jsonpath - https://github.com/NodePrime/jsonpath/blob/master/README.md
- gorilla/websocket
- highcharts

## License
Licensed under the MIT license.

See LICENSE
