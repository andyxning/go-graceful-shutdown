package grace

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	// ShutdownTimeout wait timeout in seconds
	ShutdownTimeout int = 3
)

// GraceServer is used to replace `http.Server`
type GraceServer struct {
	http.Server
	// We need the internal listener to shutdown http service gracefully.
	graceListener net.Listener
	Timeout       int
	ShutdownChan  chan os.Signal
	ExitChan      chan bool
}

// ListenAndServe listens on the TCP network address srv.Addr and handles HTTP
// requests.
//
// Note: This method is copied from `net/http/server.go`. I know this is ugly.
//       However, i do not know other ways to achieve the same goal.
//       If you have a good way to do the same job. RP is welcome. :)
func (srv *GraceServer) ListenAndServe() (err error) {
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// determine which timeout values to use: the instance one or the global one
	var timeout int
	if srv.Timeout == 0 {
		timeout = ShutdownTimeout
	} else {
		// A user defined timeout has been specified and it is not zero.
		timeout = srv.Timeout
	}

	srv.graceListener = tcpKeepAliveListener{ln.(*net.TCPListener)}

	// register graceful shutdown signal `SIGTERM`
	signal.Notify(srv.ShutdownChan, syscall.SIGTERM)

	// run shutdown signal monitor
	go func() {
		select {
		case sig := <-srv.ShutdownChan:
			log.Println("Receive shutdown signal", sig)
			// this will cause `srv.Server.Serve` to return with an error.
			srv.graceListener.Close()
			go func() {
				for {
					if defaultGraceHTTPBarrier.GetCounter() == 0 {
						defaultGraceHTTPBarrier.Barrier <- true
						break
					} else {
						time.Sleep(time.Millisecond * time.Duration(10))
					}
				}
			}()
			select {
			case <-time.After(time.Second * time.Duration(timeout)):
				log.Printf("Shutdown timeout in %ds\n", timeout)
				log.Printf("Shutdown!!!. There are still %d HTTP connections\n", defaultGraceHTTPBarrier.GetCounter())
			case <-defaultGraceHTTPBarrier.Barrier:
				log.Print("Shutdown gracefully. :)")
			}
			// we can exit now. :)
			srv.ExitChan <- true
		}
	}()

	// run http server
	err = srv.Server.Serve(srv.graceListener)
	// we only process the error for the reason of close the listening socket.
	// e.g., just like we invoke the `Close` method on listener.
	// all other errors causing `Serve` to return will be returned to the caller
	// directly. And, in such a situation the grace shutdown is not guaranteed!
	if v, ok := err.(*net.OpError); ok {
		if v.Err.Error() != "use of closed network connection" {
			return err
		} else {
			err = nil
		}
	}

	<-srv.ExitChan
	log.Println("Exited. :)")

	return
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
//
// Note: This method is copied from `net/http/server.go`. I know this is ugly.
//       However, i do not know other ways to achieve the same goal.
//       If you have a good way to do the same job. RP is welcome. :)
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
