# go-http-grace
Graceful shutdown for Go in HTTP based service.

# Installation
`go get github.com/andyxning/go-http-grace`

# Feature
* Graceful shutdown for Go with a HTTP service
* support shutdown timeout and it is configurable when initialize an instance of  `grace.GraceServer`
* **go-http-grace do not handle errors not caused by `Listener.Close()`. So, if
your `Serve` method returns errors for some other reasons, Nothing is guaranteed!** :)

# Usage
* Use `grace.HandleFunc` to register a route, instead of `http.HandleFunc`
* Use `grace.GraceServer` to initialize a server, instead of `http.Server`
* When declare a `grace.GraceServer`, register `grace.DefaultGraceServeMux` with `Handler` and initialize `ShutdownChan` and `ExitChan` elements.
    * Example
    ```go
    grace.GraceServer{
		Server: http.Server{
			Addr:           ADDRESS,
			Handler:        grace.DefaultGraceServeMux,
		},
        // Timeout: 2,
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
	ADDRESS string = "127.0.0.1:8081"
)

func sleep(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * time.Duration(5))
	fmt.Fprint(w, "hello world\n")
}

func check_health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK\n")
}

func main() {
	srv := &grace.GraceServer{
		Server: http.Server{
			Addr:           ADDRESS,
			Handler:        grace.DefaultGraceServeMux,
		},
        // Timeout: 2,
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
    * fourth one is used to run `curl http://localhost:8081/sleep`.

    **Gracefully Shutdown:**  
    After you run the `curl http://localhost:8081/sleep` command in the fourth terminal, then you can directly run the command ` kill PID` in the third terminal.

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
    If you want to check for the normal shutdown logic then you can wait for the `curl http://localhost:8081/sleep` to complete and then run the `kill PID`. you will find how normal shutdown happens. The normal output of the server is like this:
    > 2015/12/22 14:30:14 Receive shutdown signal terminated  
    > 2015/12/22 14:30:14 Shutdown gracefully. :)  
    > 2015/12/22 14:30:14 Exited. :)

# Workflow


# TODO
* Performance test
* Support for HTTPS
* Less copy-and-paste in `grace/server.go` from Go standard library

# License
This library is licensed under AGPL-3.0
