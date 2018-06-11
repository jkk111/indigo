package util

import (
  "os"
  "fmt"
  "flag"
  "path"
  "os/user"
  "github.com/satori/go.uuid"
)

const FILE_MODE = 0700

func Mkdir(path string) {
  err := os.MkdirAll(path, FILE_MODE)
  if err != nil {
    if !os.IsExist(err) {
      panic(err)
    }
  }
}

func Rmdir(path string) {
  fmt.Println("Removing", path)
  os.RemoveAll(path) 
}

func RandomId() string {
  return uuid.Must(uuid.NewV4()).String()
}

func DataDir() string {
  current, err := user.Current()

  if err != nil {
    panic(err)
  }

  data_path := path.Join(current.HomeDir, ".indigo")
  return data_path
}

func Path(res string) string {
  return path.Join(DataDir(), res)
}

func init() {
  resetPtr := flag.Bool("reset", false, "Clear local data.")
  flag.Parse()
  reset := *resetPtr

  if reset {
    fmt.Println("Cleaing All Stored Data")
    os.RemoveAll(DataDir())
  }

  Mkdir(DataDir())
  Hide(DataDir())
}