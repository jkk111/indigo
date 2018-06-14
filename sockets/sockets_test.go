package sockets

import (
  "testing"
  "github.com/jkk111/indigo/util"
)

func TestListen(t * testing.T) {
  socket := util.GetSocket("test", 0)
  ln, err := Listen(socket)

  if err != nil {
    t.Error(err)
  }

  ln.Close()
}

func TestDial(t * testing.T) {
  socket := util.GetSocket("test", 1)
  ln, err := Listen(socket)
  running := true

  go func() {
    for {
      conn, err := ln.Accept()
      if running && err != nil {
        t.Error("Failed To handle Connection")
      } else if running {
        conn.Close()
      }
    }
  }()

  if err != nil {
    t.Error("Failed To Setup Listener")
  }

  conn, err := Dial(socket)

  if err != nil {
    t.Error("Failed To Dial")
  }
  running = false
  conn.Close()
  ln.Close()
}