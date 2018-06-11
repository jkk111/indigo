package proxy

import (
  "net"
  "context"
  "github.com/jkk111/indigo/sockets"
)

func (this * BetterRoundTripper) socketRoundTrip(ctx context.Context, network string, addr string) (c net.Conn, err error) {
  return sockets.Dial(addr)
}