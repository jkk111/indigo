// +build !windows

package proxy

import (
  "net"
  "context"
  "strings"
  "github.com/jkk111/indigo/sockets"
)

func (this * BetterRoundTripper) socketRoundTrip(ctx context.Context, network string, addr string) (c net.Conn, err error) {
  sockets.Dial(addr)
}