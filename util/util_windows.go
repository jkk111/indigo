package util

import (
  "syscall"
  "fmt"
)

func Hide(path string) {
  nameptr, err := syscall.UTF16PtrFromString(path)

  if err != nil {
    panic(err)
  }

  err = syscall.SetFileAttributes(nameptr, syscall.FILE_ATTRIBUTE_HIDDEN)

  if err != nil {
    panic(err)
  }
}

// Windows doesn't support sockets, we can fake it using named pipes
// Mainly for development, as ideally, we'll be running on linux host
// file://./pipe/<svc_name>-<instance_no>.sock
func GetSocket(name string, instance int) string {
  return fmt.Sprintf(`//./pipe/%s-%d.sock`, name, instance)
}