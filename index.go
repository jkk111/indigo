package main

import "net"
import "net/http"
import "fmt"
import "os"
import "os/signal"
import Proxy "github.com/jkk111/indigo/proxy"
import "github.com/jkk111/indigo/admin"
import "github.com/jkk111/indigo/database"
import "github.com/jkk111/indigo/services"
import "github.com/jkk111/indigo/sockets"

var proxy = Proxy.NewReverseProxy()
var srv http.Server
var close chan os.Signal
var ready chan bool = make(chan bool, 1)

func cleanup(c chan os.Signal) {
  for sig := range c {
    database.Instance.Close()
    if sig == os.Interrupt {
      srv.Shutdown(nil)
    }
  }
}

func init() {
  close = make(chan os.Signal, 1)
  signal.Notify(close, os.Interrupt)
  go cleanup(close)
}

func StartServer() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", proxy.Router)
  mux.Handle("/admin/", admin.Router)

  port := os.Getenv("PORT")
  socket := os.Getenv("SOCKET")

  var ln net.Listener
  var err error

  if port == "" {
    port = ":80"
  } else if port[0] != ':' {
    port = fmt.Sprintf(":%s", port)
  }

  if socket != "" {
    ln, err = sockets.Listen(socket)
  } else {
    ln, err = net.Listen("tcp", port)
  }

  if err != nil {
    panic(err)
  }

  srv = http.Server{ Handler: mux }
  ready <- true

  if err := srv.Serve(ln); err != nil {
    fmt.Println(err)
  }
}

func main() {
  test_mode := os.Getenv("TEST")

  if test_mode == "true" {
    os.Exit(0)
  }

  services.Load()
  StartServer()
}
