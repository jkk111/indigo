package admin

import (
  "bytes"
  "encoding/json"
  "testing"
  "net/http"
  "net/url"
  "github.com/jkk111/indigo/sockets"
  "github.com/jkk111/indigo/proxy"
  "github.com/jkk111/indigo/util"
  "github.com/jkk111/indigo/database"
)

const REPO = "https://github.com/jkk111/indigo"

func Setup() {
  util.BASE_NAME = ".indigo-test"
  database.Instance.Close()
  database.Setup()
}

func TearDown() {
  database.Instance.Close()
  util.Rmdir(util.DataDir())
}

type CloseBuffer struct {
  *bytes.Buffer
}

func (this * CloseBuffer) Close() error {
  return nil
}

func TestAdminAPI(t * testing.T) {
  Setup()
  socket := util.GetSocket("indigo-test", 2)
  ln, err := sockets.Listen(socket)

  if err != nil {
    t.Error("Failed to create test server")
  }

  mux := http.NewServeMux()
  mux.Handle("/", Router)
  srv := http.Server{ Handler: mux }

  go srv.Serve(ln)

  headers := make(http.Header, 0)
  headers.Add("X-Dest", "local")

  URL := &url.URL{
    Host: socket,
    Path: "/admin/",
  }

  req := &http.Request{
    Method: "GET",
    URL: URL,
    Header: headers,
  }

  if err != nil {
    t.Error("Failed To Connect To Test Server")
  }

  cli := http.Client{
    Transport: proxy.NewBetterRoundTripper(nil),
  }

  resp, err := cli.Do(req)

  if err != nil {
    t.Error("Failed To Connect To Server")
  }

  if resp.StatusCode != 200 {
    t.Error("Invalid Response Code")
  }

  URL.Path = "/admin/404-test"

  resp, err = cli.Do(req)

  if err != nil {
    t.Error("Failed To Connect To Server")
  }

  if resp.StatusCode != 404 {
    t.Error("Invalid Response Code", resp.StatusCode)
  }

  URL.Path = "/admin/services"

  resp, err = cli.Do(req)

  if err != nil {
    t.Error("Failed To Connect To Server")
  }

  if resp.StatusCode != 200 {
    t.Error("Invalid Response Code", resp.StatusCode)
  }

  URL.Path = "/admin/branches"
  URL.RawQuery = "repo=" + REPO

  resp, err = cli.Do(req)

  if err != nil {
    t.Error("Failed To Connect To Server")
  }

  if resp.StatusCode != 200 {
    t.Error("Invalid Response Code", resp.StatusCode)
  }

  service := &database.Service{
    Id: 4,
    Name: "test-service",
    Desc: "Test Description",
    Enabled: true,
    Start: "Test",
    Args: []string { "start" },
    Env: []string { "TEST=${datadir}/test" },
    Host: "*",
    Path: "/static",
    Install: "Test",
    InstallArgs: []string { "install" },
    InstallEnv: []string {},
    Repo: REPO,
  }

  URL.Path = "/admin/add_service"

  buf, err := json.Marshal(service)

  if err != nil {
    t.Error("Failed To Serialize Service")
  }

  req.Method = "POST"
  req.Body = &CloseBuffer{bytes.NewBuffer(buf)}

  resp, err = cli.Do(req)

  if err != nil {
    t.Error("Failed To Connect To Server")
  }

  if resp.StatusCode != 200 {
    t.Error("Invalid Response Code", resp.StatusCode)
  }

  service.Name = "Test-Service-Updated"

  URL.Path = "/admin/update_service"

  buf, err = json.Marshal(service)

  if err != nil {
    t.Error("Failed To Serialize Service")
  }

  req.Method = "POST"
  req.Body = &CloseBuffer{bytes.NewBuffer(buf)}

  resp, err = cli.Do(req)

  if err != nil {
    t.Error("Failed To Connect To Server")
  }

  if resp.StatusCode != 200 {
    t.Error("Invalid Response Code", resp.StatusCode)
  }
}