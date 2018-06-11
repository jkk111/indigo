package sockets

import (
  "net"
  "time"
  "strings"
  "github.com/Microsoft/go-winio"
)

const TIMEOUT = time.Second * 30

func Dial(path string) (net.Conn, error) {
  path = strings.Replace(path, ":80", "", -1)
  timeout := TIMEOUT
  return winio.DialPipe(path, &timeout)
}

func Listen(path string) (net.Listener, error) {
  path = strings.Replace(path, ":80", "", -1)
  return winio.ListenPipe(path, nil) 
}