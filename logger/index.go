package logger

import (
  "io"
  "os"
  "fmt"
  "time"
  "path"
  "github.com/jkk111/indigo/util"
)

var LOG_PATH string

func init() {
  LOG_PATH = util.Path("logs")
  util.Mkdir(LOG_PATH)
}

type LogInstance struct {
  id string // Application ID
  started int64
  fout * os.File
  ferr * os.File
}

var handles = make(map[string]*LogInstance)

func CreateLogInstance(name string, out io.ReadCloser, err io.ReadCloser) * LogInstance {
  start := time.Now().UnixNano()

  fout_name := fmt.Sprintf("%s-%d-out.log", name, start)
  ferr_name := fmt.Sprintf("%s-%d-err.log", name, start)

  app_log_path := path.Join(LOG_PATH, name)

  util.Mkdir(app_log_path)

  fout, e := os.Create(path.Join(app_log_path, fout_name))
  ferr, e2 := os.Create(path.Join(app_log_path, ferr_name))

  if e != nil {
    panic(e)
  }

  if e2 != nil {
    panic(e2)
  }

  fmt.Fprintf(fout, "STDOUT Log for %s Started at %d\n", name, start)
  fmt.Fprintf(ferr, "STDERR Log for %s Started at %d\n", name, start)

  go func() {
    defer fout.Close()
    defer ferr.Close()
    // Dump to file until process exits
    fmt.Println(io.Copy(fout, out))
    fmt.Println(io.Copy(ferr, err))
  }()

  inst := &LogInstance{ 
    name,
    start,
    fout,
    ferr,
  }

  handles[name] = inst

  return inst
}

