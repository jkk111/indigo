package proxy

import (
  "net"
  "time"
  "context"
  "strings"
  "github.com/Microsoft/go-winio"
)

func (this * BetterRoundTripper) socketRoundTrip(ctx context.Context, network string, addr string) (c net.Conn, err error) {
  addr = strings.Replace(addr, ":80", "", -1) // Sockets don't use ports, strip it
  timeout := time.Second * 30
  return = winio.DialPipe(addr, &timeout)
}