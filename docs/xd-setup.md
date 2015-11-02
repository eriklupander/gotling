# Setting up Spring XD as load target

To test a pure "write" scenario when developing a Gotling test, Spring XD can be utilized to provide HTTP and TCP endpoints in a convenient way.

### XD distributed mode setup
Running XD in distributed mode needs a number of subsystems to work properly, this is how to set them up:

#### ZooKeeper
To start the server (on Mac OS X or Linux. Windows, don't know): 
    bin/zkServer.sh start-foreground ../spring-xd-1.2.1.RELEASE/zookeeper/conf/zoo.cfg. 

#### Redis
Spring XD ships with the Redis source code and can be compiled or downloaded-
- Go to the redis/bin subdirectory under the Spring XD install path and run _install-redis_ script to build
- From the root path of the Redis install, start by executing _../bin/redis-server redis.conf_ 

#### HSQLDB
Change to the hsqldb path in the Spring XD install dir and start HSQLDB by running bin/hsqldb-server

#### SPRING XD
(1) Change the transport from Rabbit to Redis by uncommenting and modifying the xd.transport parameter as shown:

    #XD data transport (default is redis for distributed, local for single node)
    xd:
        transport: redis
        
(2) Configure the HSQLDB connection parameters, in this case the included default values will suffice - it is necessary to just uncomment the parameters.
     
     
    #Change the database host, port and name
    hsql:
      server:
        host: localhost
        port: 9101
        dbname: xdjob
    #Change database username and password
    spring:
      datasource:
        url: jdbc:hsqldb:hsql://${hsql.server.host:localhost}:${hsql.server.port:9101}/${hsql.server.dbname:xdjob}
        username: sa
        password:
        driverClassName: org.hsqldb.jdbc.JDBCDriver
        validationQuery: select 1 from INFORMATION_SCHEMA.SYSTEM_USERS
(3) Let's also enable the ability to shut down containers from the Admin UI. Follow the instructions in the configuration file, and uncomment the parameters.
    
    #---
    spring:
      profiles: container
    management:
      port: 0
(4) Setup the Redis connection parameters, using the default host and port.
    
    spring:
      redis:
       port: 6379
       host: localhost
(5) Finally, setup the ZooKeeper connection properties.
     
    # namespace is the path under the root where XD's top level nodes will be created
    # client connect string: host1:port1,host2:port2,...,hostN:portN
    zk:
      namespace: xd
      client:
         connect: localhost:2181


### Starting Spring XD
- Open a new shell, and start the Admin server by executing the spring-xd-1.2.1.RELEASE/xd/bin/xd-admin shell script.
- Start up a container in another shell by executing spring-xd-1.2.1.RELEASE/xd/bin/xd-container.
- Start a XD shell by: spring-xd-1.2.1.RELEASE/shell/bin> ./xd-shell


### Creating some streams and taps
TCP listener

    stream create --name tcptest --definition "tcp --port=8081 --outputType=text/plain | log" --deploy
HTTP listener

    stream create --name httptest --definition "http --port=10000 | log" --deploy
TCP tap/counter

    stream create --name tcptap --definition "tap:stream:tcptest > counter --name=tcpcount" --deploy    
HTTP tap/counter

    stream create --name httptap --definition "tap:stream:httptest > counter --name=httpcount" --deploy     
    
    
### Access Web GUI
http://localhost:9393/admin-ui

Remember that counters only becomes available in the GUI once some data have passed into them.


### Sample use-case: Fleet data

OpenStreetMap API's can be used to generate polyline or .gpx files for a road-following route:
Example, Gothenburg to Trollhättan

    http://router.project-osrm.org/viaroute?loc=57.703033,12.000526&loc=57.166582,12.529290&instructions=false&alt=false&output=gpx
    
See https://github.com/Project-OSRM/osrm-backend/wiki/Server-api

See also: http://mathematica.stackexchange.com/questions/11521/using-j-link-to-decode-google-maps-polyline-strings

The .gpx file is probably the easiest to transform into a feeder. To provide a good simulation, we'll probably need to generate a few hundred or so routes with random "start point" along the route.
I think the feeder format needs to be something like:

    route1,route2,route3,...,routeN
    57.10|12.30,57.90|12.70,58.30|11.40,...
    57.11|12.31,57.91|12.71,58.31|11.41,...
    57.12|12.32,57.92|12.72,58.32|11.42,...
    
I.e. each route goes into a column of its own. This won't quite work unless all routes are normalized to the same number of coords but that can probably be arranged.

So - first generate a number of routes using the OSM API (beware of their usage terms, no spamming!) and save as .gpx files. Then put those into a folder or something and build a matrix we can transform into the appropriate csv format.

Another option is to allow "namespaced" feeders. E.g. one user gets exclusive right to _one_ csv resource. Each CSV resource is namespaced by its file name and then bound to a user and its values accessed through $${varname} or something denoting them from "global" feeders.