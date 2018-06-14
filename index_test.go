package main

import (
  "os"
  "testing"
)

func TestMain(t * testing.T) {
  go main()
  <- ready
  close <- os.Interrupt
  srv.Shutdown(nil)
}