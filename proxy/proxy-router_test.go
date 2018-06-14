package proxy

import (
  "github.com/jkk111/indigo/sockets"
  "github.com/jkk111/indigo/util"
  "testing"
  "net/http"
  "io/ioutil"
  "net/url"
)

func Setup() {
  NewReverseProxy()
}

func TearDown() {
  ProxyInstance = nil
}

func TestInstance(t * testing.T) {
  Setup()
  inst := Instance()
  if inst == nil {
    t.Error("Expected Proxy Instance")
  }
  TearDown()
}

func TestAddRoute(t * testing.T) {
  Setup()
  inst := Instance()
  dummy_socket := util.GetSocket("test-listener", 0)
  dummy_socket2 := util.GetSocket("test-listener-2", 0)
  inst.AddRoute("*", "/TestPath", dummy_socket2, true)

  ln, err := sockets.Listen(dummy_socket)

  if err != nil {
    t.Error("Failed to create local http server")
  }

  ln2, err := sockets.Listen(dummy_socket2)

  if err != nil {
    t.Error("Failed to create local http server")
  }

  mux := http.NewServeMux()
  mux.HandleFunc("/", inst.Router)

  server := http.Server{ Handler: mux }
  go server.Serve(ln)

  mux2 := http.NewServeMux()
  mux2.HandleFunc("/", func(w http.ResponseWriter, r * http.Request) {
    w.Write([]byte("Hello World"))
  })

  server2 := http.Server{ Handler: mux2 }
  go server2.Serve(ln2)

  cli := http.Client{
    Transport: NewBetterRoundTripper(nil),
  }

  headers := make(http.Header, 0)
  headers.Add("X-Dest", "local")

  URL := &url.URL{
    Host: dummy_socket,
    Path: "/",
  }

  request := &http.Request{
    Method:"GET",
    URL: URL,
    Header:  headers,
  }

  resp, err := cli.Do(request)

  if err != nil {
    t.Error("Failed To Connect To Test Proxy")
  }

  buf, err := ioutil.ReadAll(resp.Body)

  if err != nil {
    t.Error("Failed To Read HTTP Response")
  }

  if string(buf) != "Invalid Route" {
    t.Error("Unexpected HTTP Response Expected", "Invalid Route", "Got:", string(buf))
  }

  URL = &url.URL{
    Host: dummy_socket,
    Path: "/TestPath",
  }

  request = &http.Request{
    Method:"GET",
    URL: URL,
    Header:  headers,
  }

  resp, err = cli.Do(request)

  if err != nil {
    t.Error("Failed To Connect To Test Proxy")
  }

  buf, err = ioutil.ReadAll(resp.Body)

  if err != nil {
    t.Error("Failed To Read HTTP Response")
  }

  if string(buf) != "Hello World" {
    t.Error("Unexpected HTTP Response Expected", "Hello World", "Got:", string(buf))
  }



  server2.Close()
  server.Close()

  TearDown()
}