// +build !windows

package sockets

import (
  "net"
  "time"
  "strings"
)

const TIMEOUT = time.Second * 30

func Dial(path string) (net.Conn, error) {
  path = strings.Replace(path, ":80", "", -1)
  timeout := TIMEOUT
  return net.Dial("unix", addr) 
}

func Listen(path string) (net.Listener, error) {
  path = strings.Replace(path, ":80", "", -1)
  return net.Listen("unix", path) 
}