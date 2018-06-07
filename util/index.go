package util

import (
  "os"
)

const FILE_MODE = 0700

func Mkdir(path string) {
  err := os.Mkdir(path, FILE_MODE)
  if err != nil {
    if !os.IsExist(err) {
      panic(err)
    }
  }
}