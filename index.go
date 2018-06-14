package main

import "net/http"
import "fmt"
import "os"
import "os/signal"
import Proxy "github.com/jkk111/indigo/proxy"
import "github.com/jkk111/indigo/admin"
import "github.com/jkk111/indigo/database"
import "github.com/jkk111/indigo/services"

var proxy = Proxy.NewReverseProxy()
var srv http.Server

func cleanup(c chan os.Signal) {
  for sig := range c {
    fmt.Println("Received Exit Signal", sig)
    database.Instance.Close()
    if sig == os.Interrupt {
      srv.Shutdown(nil)
    }
  }
}

func init() {
  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)
  go cleanup(c)
}

func StartServer() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", proxy.Router)
  mux.Handle("/admin/", admin.Router)

  port := os.Getenv("PORT")

  if port == "" {
    port = ":80"
  } else if port[0] != ':' {
    port = fmt.Sprintf(":%s", port)
  }

  srv = http.Server{Addr: port, Handler: mux}
  if err := srv.ListenAndServe(); err != nil {
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
