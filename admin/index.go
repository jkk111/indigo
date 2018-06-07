package admin

import (
  "encoding/json"
  "net/http"
  "fmt"
  "path"
  // "strings"
  "github.com/jkk111/indigo/assets"
  "github.com/jkk111/indigo/database"
)

var router * http.ServeMux
var Router * http.ServeMux

func serve_static_asset(w http.ResponseWriter, r * http.Request) {
  p := r.URL.EscapedPath()
  urlpath := p
  fmt.Println(urlpath)

  data, err := assets.Asset(fmt.Sprintf("resources/admin/%s", urlpath))

  if err == nil {
    w.Write(data)
    return
  } else {
    fmt.Printf("resources/admin/%s\n", urlpath)
  }

  data, err2 := assets.Asset(path.Join("resources/admin", urlpath, "index.html"))
  fmt.Println(path.Join("resources/admin", urlpath, "index.html"), fmt.Sprintf("resources/admin/%s", urlpath), data, err, err2)
  if err2 == nil {
    w.Write(data)
    return
  }

  w.Write([]byte("404"))
}

func services(w http.ResponseWriter, r * http.Request) {
  marshalled, err := json.Marshal(database.Services())
  if err != nil {
    panic(err)
  }
  w.Write(marshalled)
}

func verify(w http.ResponseWriter, r * http.Request) {
  r.URL.Path = r.URL.Path[6:]
  router.ServeHTTP(w, r)
}

func init() {
  Router = http.NewServeMux()
  Router.HandleFunc("/", verify)
  router = http.NewServeMux()
  router.HandleFunc("/", serve_static_asset)
  router.HandleFunc("/services", services)
}