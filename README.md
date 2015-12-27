# go-http-grace

[![Join the chat at https://gitter.im/andyxning/go-http-grace](https://badges.gitter.im/andyxning/go-http-grace.svg)](https://gitter.im/andyxning/go-http-grace?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Build Status](https://travis-ci.org/andyxning/go-http-grace.svg)](https://travis-ci.org/andyxning/go-http-grace)  
Graceful shutdown for Go in HTTP based service.

# Installation
`go get github.com/andyxning/go-http-grace`

# Feature
* Graceful shutdown for Go with a HTTP service
* support shutdown timeout and it is configurable when initialize an instance of  `grace.Server`
* **go-http-grace do not handle errors not caused by `Listener.Close()`. So, if
your `http.Server.Serve` method returns errors for some other reasons.
 Nothing is guaranteed!** :)

# Usage
* Use `grace.HandleFunc` to register a route, instead of `http.HandleFunc`
* Use `grace.Server` to initialize a server, instead of `http.Server`
* When declare a `grace.Server`, register `grace.DefaultServeMux` with `Handler` and initialize `ShutdownChan` and `ExitChan` elements.
    * Example
    ```go
    grace.Server{
		Server: http.Server{
			Addr:           Address,
			Handler:        grace.DefaultServeMux,
		},
        // Timeout: 10,
		ShutdownChan: make(chan os.Signal, 1),
		ExitChan:     make(chan bool, 1),
	}
    ```

# Example
* First, you should install `go-http-grace`
* Then, you can copy the content of the code below. Assume you save it with `demo.go`

```go
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/andyxning/go-http-grace/grace"
)

const (
	Address string = "127.0.0.1:8081"
)

func sleep(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * time.Duration(10))
	w.Write([]byte("hello world\n"))
}

func check_health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK\n")
}

func main() {
	srv := &grace.Server{
		Server: http.Server{
			Addr:           Address,
			Handler:        grace.DefaultServeMux,
			MaxHeaderBytes: 1 << 30,
		},
		// Timeout:      10,
		ShutdownChan: make(chan os.Signal, 1),
		ExitChan:     make(chan bool, 1),
	}

	grace.HandleFunc("/", sleep)
	grace.HandleFunc("/healthcheck", check_health)

	srv.ListenAndServe()
}

```
* open your browser and input `127.0.0.1:8081/healthcheck` to check if the server is running. If you see a string `OK` in
the output, then it looks that everything is good for now. Congratulations. :)
* use `ps -ef|grep demo` to check for the running server's **pid**.
* you should open two more terminals
    * first one is used to start the server. Already done in step One. :)
    * second one is to run `ps -ef|grep demo`. Already done in step Two. :)
    * third one is used to run `kill PID`. **PID** is the pid of the running server. Get it from the second terminal. :)
    * fourth one is used to run `curl http://127.0.0.1:8081/sleep`.

    **Gracefully Shutdown:**  
    After you run the `curl http://127.0.0.1:8081/sleep` command in the fourth terminal, then you can directly run the command ` kill PID` in the third terminal.

    You can now get the output from the first terminal like this(**Note: the time of second line is about 4 seconds later than the first lien's**):
    > 2015/12/22 14:02:53 Receive shutdown signal terminated
    > 2015/12/22 14:02:57 Shutdown gracefully. :)  
    > 2015/12/22 14:02:57 Exited. :)

    and, get the output from the fourth terminal like this:
    > hello world

    This means that we shutdown the server gracefully.

    **Timeout Shutdown:**  
    If you want to check for the timeout logic then you can set the sleep time in `sleep` function to be `2x` of the timeout time of shutdown. Then, run the same test just like the above. You will find how timeout works. The timeout output of the server is like this(**Note: the time of second line is 5 seconds later than the first line's**):
    > 2015/12/22 14:32:20 Receive shutdown signal terminated  
    > 2015/12/22 14:32:25 Shutdown timeout in 5s  
    > 2015/12/22 14:32:25 Shutdown!!!. There are still 1 HTTP connections  
    > 2015/12/22 14:32:25 Exited. :)

    **Normal Shutdown:**  
    If you want to check for the normal shutdown logic then you can wait for the `curl http://127.0.0.1:8081/sleep` to complete and then run the `kill PID`. you will find how normal shutdown happens. The normal output of the server is like this:
    > 2015/12/22 14:30:14 Receive shutdown signal terminated  
    > 2015/12/22 14:30:14 Shutdown gracefully. :)  
    > 2015/12/22 14:30:14 Exited. :)

# Performance Test
* Tools: `ab`.
* Test command: `ab -n 10000 -c 128 http://127.0.0.1:8081/healthcheck`
* Environment:
    * Darwin Kernel Version 14.5.0
    * 8 GB 1867 MHz DDR3
    * 2.7 GHz Intel Core i5
    * OS X Yosemite 10.10.5
* Test request: "/healthcheck"
* `net/http` server source code

```go
package main

import (
	"fmt"
	"net/http"
	"time"
)

const (
	Address string = "127.0.0.1:8081"
)

func sleep(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * time.Duration(10))
	w.Write([]byte("hello world\n"))
}

func check_health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK\n")
}

func main() {
	srv := &http.Server{
		Addr:           Address,
		MaxHeaderBytes: 1 << 30,
	}

	http.HandleFunc("/", sleep)
	http.HandleFunc("/healthcheck", check_health)

	srv.ListenAndServe()
}
```
* Conclusion: **The performance of `go-http-grace` is almost the same as `net/http`**
## go-http-grace

```
Server Software:
Server Hostname:        127.0.0.1
Server Port:            8081

Document Path:          /healthcheck
Document Length:        3 bytes

Concurrency Level:      128
Time taken for tests:   1.707 seconds
Complete requests:      10000
Failed requests:        0
Total transferred:      1190000 bytes
HTML transferred:       30000 bytes
Requests per second:    5858.88 [#/sec] (mean)
Time per request:       21.847 [ms] (mean)
Time per request:       0.171 [ms] (mean, across all concurrent requests)
Transfer rate:          680.87 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        2   12  51.1      7     562
Processing:     4   10  35.8      7     562
Waiting:        3   10  35.8      7     562
Total:          7   22  62.6     14     573

Percentage of the requests served within a certain time (ms)
  50%     14
  66%     15
  75%     16
  80%     16
  90%     17
  95%     18
  98%     19
  99%    570
 100%    573 (longest request)
```

## net/http

```
Server Software:
Server Hostname:        127.0.0.1
Server Port:            8081

Document Path:          /healthcheck
Document Length:        3 bytes

Concurrency Level:      128
Time taken for tests:   1.725 seconds
Complete requests:      10000
Failed requests:        0
Total transferred:      1190000 bytes
HTML transferred:       30000 bytes
Requests per second:    5796.20 [#/sec] (mean)
Time per request:       22.083 [ms] (mean)
Time per request:       0.173 [ms] (mean, across all concurrent requests)
Transfer rate:          673.58 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        2    9  37.1      7     589
Processing:     2   12  53.9      7     589
Waiting:        1   12  54.0      7     589
Total:          5   22  65.5     14     599

Percentage of the requests served within a certain time (ms)
  50%     14
  66%     15
  75%     15
  80%     16
  90%     17
  95%     18
  98%     21
  99%    595
 100%    599 (longest request)
```

# Workflow


# TODO
* Support for HTTPS
* Less copy-and-paste in `grace/server.go` from Go standard library

# License
This library is licensed under AGPL-3.0
