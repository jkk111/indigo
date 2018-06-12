package proxy

import (
  "net"
  "net/http"
  "net/http/httputil"
  "fmt"
  "time"
  "strings"
  URL "net/url"
  "bytes"
  "github.com/jkk111/indigo/assets"
  "reflect"
)

var transport * BetterRoundTripper = newBetterRoundTripper(nil)
var ProxyInstance * ReverseProxy

type BetterRoundTripper struct {
  transport http.RoundTripper
}

func newBetterRoundTripper(tr http.RoundTripper) * BetterRoundTripper {
  if tr == nil {
    tr = http.DefaultTransport
  }

  return &BetterRoundTripper{ tr }
}

type BufferedCloser struct {
  *bytes.Buffer
}

func (this * BufferedCloser) Close() (err error) {
  return
}

func NewBufferedCloser(buf []byte) * BufferedCloser {
  b := bytes.NewBuffer(buf)
  return &BufferedCloser{b}
}

func conn_failed() * http.Response {
  r := &http.Response{}
  r.StatusCode = 502
  r.Body = NewBufferedCloser(assets.MustAsset("resources/ServiceUnavailable.html"))
  return r
}

func conn_refused() * http.Response {
  r := &http.Response{}
  r.StatusCode = 502
  r.Body = NewBufferedCloser([]byte("Connection Failed, Server Refused Connection"))
  return r
}

func handle_proxy_error(req * http.Request, err error) * http.Response {
  fmt.Println(err)
  switch t := err.(type) {
    case *net.OpError:
      if t.Op == "dial" {
        return conn_failed()
      } else if t.Op == "read" {
        return conn_refused()
      }

    default:
      fmt.Println("Unknown Error:", reflect.TypeOf(err))
      return conn_failed()
  }

  return nil
}

// We store some essential paths in the headers to handle passing,
// Easy in v 1-3 because node allows dynamic mutation of the request
// Work-around for golang
func (this * BetterRoundTripper) SocketRoundTrip(req * http.Request) (res * http.Response, err error) {
  req.URL.Scheme = "http"

  tp := &http.Transport{
    Proxy: http.ProxyFromEnvironment,
    DialContext: this.socketRoundTrip,
    MaxIdleConns:          100,
    IdleConnTimeout:       90 * time.Second,
    TLSHandshakeTimeout:   10 * time.Second,
    ExpectContinueTimeout: 1 * time.Second,
  }

  return tp.RoundTrip(req)
}

func (this * BetterRoundTripper) RoundTrip(req * http.Request) (*http.Response, error) {
  var res * http.Response
  var err error

  dest := req.Header.Get("X-Dest")

  if dest == "" {
    return conn_refused(), nil
  }

  // Special Case, Windows specific named pipe.
  if dest == "local" {
    res, err = this.SocketRoundTrip(req)
  } else if dest == "unix" { // Unix Socket
    // res, err := UnixRoundTrip(req)
  } else {
    res, err = this.transport.RoundTrip(req) 
  }

  if err != nil {
    return handle_proxy_error(req, err), nil
  } 

  return res, err
}

type ProxyRuleSet struct {
  Rules map[string]*HttpProxyRules
}

func NewProxyRuleSet() * ProxyRuleSet {
  rules := make(map[string]*HttpProxyRules)
  return &ProxyRuleSet{ Rules: rules }
}

func (this * ProxyRuleSet) domain(domain string) {
  if this.Rules[domain] == nil {
    this.Rules[domain] = NewHttpProxyRules()
  }
}

func (this * ProxyRuleSet) Match(domain string, url string) (route * RuleMatch) {
  this.domain(domain)
  this.domain("*")
  wildcard_match := this.Rules["*"].Match(url)
  match := this.Rules[domain].Match(url)

  if wildcard_match == nil {
    return match
  } else if match == nil {
    return wildcard_match
  }

  if wildcard_match.Strength == -1 &&  match.Strength == -1 {
    return nil
  } 

  if wildcard_match.Strength > match.Strength {
    return wildcard_match
  } else {
    return match
  }
}

func (this * ProxyRuleSet) Add(domain string, url string, route string, local bool) {
  this.domain(domain)
  this.Rules[domain].Add(url, route, local)
}

type ReverseProxy struct {
  proxy_rules * ProxyRuleSet
}

func NewReverseProxy() * ReverseProxy {
  rp := &ReverseProxy {
    proxy_rules: NewProxyRuleSet(),
  }

  ProxyInstance = rp

  return rp
}

func (this * ReverseProxy) AddRoute(domain string, url string, route string, local bool) {
  fmt.Println("Add Route", this)
  this.proxy_rules.Add(domain, url, route, local)
}

func (this * ReverseProxy) RemoveRoute(domain string, url string) {

}

func (this * ReverseProxy) Router(w http.ResponseWriter, r * http.Request) {
  path := r.URL.Path
  host := r.URL.Hostname()

  if host == "" {
    host = r.Header.Get("Host")
  }

  if host == "" {
    host = "localhost"
  }

  fmt.Println("Initial Host", host, path, r.Header["Host"])

  if host == "" {
    host = r.Host
  }

  match := this.proxy_rules.Match(host, path)

  if match != nil {
    fixed_path := strings.Replace(path, match.Prefix, "", -1)

    if fixed_path == "" {
      fixed_path = "/"
    } else if fixed_path[0] != '/' {
      fixed_path = "/" + fixed_path
    }

    m := match.Match

    var url * URL.URL

    if match.Local {
      url = &URL.URL{
        Scheme: "http",
        Path: fixed_path,
        Host: m,
      }
    } else {
      url, err := url.Parse(m)

      if err != nil {
        w.Write([]byte("Broke"))
        return
      }

      url.Path = fixed_path      
    }

    proxy := httputil.NewSingleHostReverseProxy(url)
    proxy.Transport = transport
    r.URL.Path = fixed_path

    r.Header.Set("X-Original-Host", host)
    r.Header.Set("X-Original-Path", path)
    r.Header.Set("X-Proxied-Path", fixed_path)
    r.Header.Set("X-Proxied-URL", m)
    var dest string

    if match.Local {
      dest = "local"
    } else {
      dest = "remote"
    }

    r.Header.Set("X-Dest", dest)

    proxy.ServeHTTP(w, r)
  
  } else {
    fmt.Println("404", host, path)
    w.Write([]byte("Invalid Route"))
  }
}

func SetProxyTransport(tr http.RoundTripper) {
  transport.transport = tr
}

func Instance() * ReverseProxy {
  return ProxyInstance
}