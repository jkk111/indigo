package proxy

import "regexp"
import "strings"
import "fmt"

type ProxyRule struct {
  Route string // Socket // Host
  Local bool
}

type HttpProxyRules struct {
  Rules map[string]*ProxyRule
  Default string
}

func NewHttpProxyRules() * HttpProxyRules {
  rules := make(map[string]*ProxyRule)
  return &HttpProxyRules{ Rules: rules }
}

func (this * HttpProxyRules) Add(path string, route string, local bool) {
  if this.Rules[path] == nil {
    this.Rules[path] = &ProxyRule{route, local}
  } else {
    fmt.Println(path, route)
    panic("Overwriting Existing Rule")
  }
}

type RuleMatch struct {
  Match string
  Prefix string
  Strength int
  Local bool
}

func (this * HttpProxyRules) Match(url string) * RuleMatch {
  if url[0] != '/' {
    url = "/" + url
  }

  path := url
  var bestMatchStrength int
  var bestMatch * ProxyRule
  var bestPrefix string
  matched := false

  for prefix, rule := range this.Rules {
    var trailing_slash bool
    var pathPrefixRe * regexp.Regexp

    if prefix[len(prefix) - 1] == '/' {
      pathPrefixRe = regexp.MustCompile(strings.Replace(prefix, "/", "\\/", -1))
      trailing_slash = true
    } else {
      pathPrefixRe = regexp.MustCompile("(" + strings.Replace(prefix, "/", "\\/", -1) + `)(?:\W|$)`)
      trailing_slash = false
    }

    testPrefixMatchIndex := pathPrefixRe.FindStringIndex(path)
    testPrefixMatch := pathPrefixRe.FindStringSubmatch(path)

    if testPrefixMatchIndex != nil && testPrefixMatchIndex[0] == 0 {
      var url_prefix string
      if trailing_slash {
        url_prefix = testPrefixMatch[0]
      } else {
        url_prefix = testPrefixMatch[1]
      }

      matchStrength := computeMatchStrength(path, testPrefixMatch[0])

      if matchStrength > bestMatchStrength {
        if url_prefix != "" {
          bestPrefix = url_prefix
        }

        bestMatch = rule
        bestMatchStrength = matchStrength
        matched = true
      }
    }
  }

  if !matched {
    bestMatchStrength = -1
    return nil
  }

  return &RuleMatch { bestMatch.Route, bestPrefix, bestMatchStrength, bestMatch.Local }
}

func computeMatchStrength(url string, match string) int {
  i := 0

  l := len(url)
  lm := len(match)

  if lm < l {
    l = lm
  }

  for ;i < l; i++ {
    if url[i] != match[i] {
      return i
    }
  }

  return i
}