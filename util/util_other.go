// +build !windows

package util

import (
  "path"
  "os"
  "fmt"
)

func Hide(p string) {
  base := path.Base(p)
  if base[0] != '.' {
    dir := path.Dir(p)
    err := os.Rename(p, path.Join(dir, "." + base))

    if err != nil {
      panic(err)
    }
  }
}

func GetSocket(name string, instance int) string {
  return fmt.Sprintf("/var/run/%s-%d.socket", name, instance)
}