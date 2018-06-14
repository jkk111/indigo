package proxy

import (
  "testing"
  "github.com/jkk111/indigo/util"
)

func TestSocketRoundTrip(t * testing.T) {
  addr := util.GetSocket("Testing", 0)
  brt := newBetterRoundTripper(nil)
  brt.socketRoundTrip(nil, "", addr)
}