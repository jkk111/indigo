package util

import (
  "os"
  "github.com/satori/go.uuid"
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

func RandomId() string {
  return uuid.Must(uuid.NewV4()).String()
}