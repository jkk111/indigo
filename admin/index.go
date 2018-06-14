package admin

import (
  "encoding/json"
  "net/http"
  "fmt"
  "path"
  // "strings"
  "github.com/jkk111/indigo/git"
  "github.com/jkk111/indigo/assets"
  "github.com/jkk111/indigo/database"
)

var router * http.ServeMux
var Router * http.ServeMux

func serve_static_asset(w http.ResponseWriter, r * http.Request) {
  p := r.URL.EscapedPath()
  urlpath := p

  data, err := assets.Asset(fmt.Sprintf("resources/admin/%s", urlpath))

  if err == nil {
    w.Write(data)
    return
  }

  data, err2 := assets.Asset(path.Join("resources/admin", urlpath, "index.html"))
  if err2 == nil {
    w.Write(data)
    return
  }

  w.WriteHeader(404)
  w.Write([]byte("404"))
}

func services(w http.ResponseWriter, r * http.Request) {
  marshalled, err := json.Marshal(database.Services())
  if err != nil {
    panic(err)
  }
  w.Write(marshalled)
}

func Branches(w http.ResponseWriter, r * http.Request) {
  qs := r.URL.Query()
  repo_qs := qs["repo"]
  if repo_qs == nil || len(repo_qs) == 0 {
    w.Write([]byte("Must Specify Repo"))
  }

  repo := repo_qs[0]

  defer func() {
    if r := recover(); r != nil {
      w.Write([]byte("[]"))
    }
  }()

  branches := git.LsRemote(repo)

  buf := must_unmarshal_raw(branches)
  w.Write(buf)
}

func must_unmarshal_raw(iface interface{}) []byte {
  buf, err := json.Marshal(iface)

  if err != nil {
    panic(err)
  }

  return buf
}

func must_unmarshal(iface interface{}) string {
  return string(must_unmarshal_raw(iface))
}

func set_empty(ptr * []string) {
  if *ptr == nil {
    *ptr = make([]string, 0)
  }
}

func add_service(w http.ResponseWriter, r * http.Request) {
  if r.Method == "POST" {
    decoder := json.NewDecoder(r.Body)
    svc := &database.Service{}
    err := decoder.Decode(svc)

    if err != nil {
      panic(err)
    }

    set_empty(&svc.Args)
    set_empty(&svc.Env)
    set_empty(&svc.InstallArgs)
    set_empty(&svc.InstallEnv)

    args := must_unmarshal(svc.Args)
    env := must_unmarshal(svc.Env)

    installArgs := must_unmarshal(svc.InstallArgs)
    installEnv := must_unmarshal(svc.InstallEnv)

    svc.StartArgsRaw = args
    svc.StartEnvRaw = env
    svc.InstallArgsRaw = installArgs
    svc.InstallEnvRaw = installEnv

    database.AddService(svc)
  } else {
    w.Write([]byte("Invalid Method"))   
  }
}

func update_service(w http.ResponseWriter, r * http.Request) {
  if r.Method == "POST" {
    decoder := json.NewDecoder(r.Body)
    svc := &database.Service{}
    err := decoder.Decode(svc)

    if err != nil {
      panic(err)
    }

    set_empty(&svc.Args)
    set_empty(&svc.Env)
    set_empty(&svc.InstallArgs)
    set_empty(&svc.InstallEnv)

    args := must_unmarshal(svc.Args)
    env := must_unmarshal(svc.Env)

    installArgs := must_unmarshal(svc.InstallArgs)
    installEnv := must_unmarshal(svc.InstallEnv)

    svc.StartArgsRaw = args
    svc.StartEnvRaw = env
    svc.InstallArgsRaw = installArgs
    svc.InstallEnvRaw = installEnv

    database.UpdateService(svc)
  } else {
    w.Write([]byte("Invalid Method"))   
  }
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
  router.HandleFunc("/add_service", add_service)
  router.HandleFunc("/update_service", update_service)
  router.HandleFunc("/branches", Branches)
}