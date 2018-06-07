package main

import "net/http"
import "fmt"
import "os"
import "os/signal"
import "io/ioutil"
import Proxy "github.com/jkk111/indigo/proxy"
import "os/user"
import "path"
import "github.com/jkk111/indigo/admin"
import "github.com/jkk111/indigo/database"
import "github.com/jkk111/indigo/util"

var proxy = Proxy.NewReverseProxy()
var srv http.Server

const FILE_MODE = 0700 // Owner Accessible Only
  
func read_file(path string) string {
  f, err := os.Open(path)

  if err != nil {
    panic(err)
  }

  defer f.Close()
  data, err := ioutil.ReadAll(f)

  if err != nil {
    panic(err)
  }

  return string(data)
}

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
  current, err := user.Current()

  if err != nil {
    panic(err)
  }

  data_path := path.Join(current.HomeDir, ".indigo")
  repo_path := path.Join(data_path, "repos")

  util.Mkdir(data_path)
  util.Mkdir(repo_path)
  util.Hide(data_path)

  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)
  go cleanup(c)
}

func StartServer() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", proxy.Router)
  mux.Handle("/admin/", admin.Router)

  srv = http.Server{Addr: ":80", Handler: mux}
  if err := srv.ListenAndServe(); err != nil {
    fmt.Println(err)
  }
}

func main() {
  services := database.Services()

  for _, service := range services {
    proxy.AddRoute(service.Host, service.Path, util.GetSocket(service.Name, 0), true)
  }

  fmt.Printf("Added %d Routes\n", len(services))

  StartServer()
}