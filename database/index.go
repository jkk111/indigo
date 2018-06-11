package database

import (
  "fmt"
  "path"
  "flag"
  "os"
  "os/user"
  "database/sql"
  "encoding/json"
  "github.com/mattn/go-sqlite3"
  "github.com/jkk111/indigo/assets"
  "github.com/jkk111/indigo/util"
)

var Instance * sql.DB

type Service struct {
  Id int `json:"id"`
  Name string `json:"name"`
  Desc string `json:"desc"`
  Host string `json:"host"`
  Path string `json:"path"`
  Repo string `json:"repo"`

  Start string `json:"string"`
  Args []string `json:"string"`
  Env []string `json:"string"`

  Install string `json:"string"`
  InstallArgs []string `json:"installArgs"`
  InstallEnv []string `json:"installEnv"`

  Enabled bool `json:"enabled"`
}

func Exec(db * sql.DB, q string) (sql.Result, error) {
  return db.Exec(q)
}

func MustExec(db * sql.DB, q string) {
  _, err := Exec(db, q)
  if err != nil {
    panic(err)
  }
}

func Query(db * sql.DB, q string) (* sql.Rows, error) {
  return db.Query(q)
}

func MustQuery(db * sql.DB, q string) * sql.Rows {
  rows, err := Query(db, q)

  if err != nil {
    panic(err)
  }

  return rows
}


/*
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  desc TEXT NOT NULL DEFAULT "echo \"No Description Specified\"",
  host TEXT NOT NULL, 
  path TEXT NOT NULL, 
  repo TEXT NOT NULL DEFAULT "echo \"No Repo Specified\"",
  start TEXT NOT NULL DEFAULT "echo \"No Start Command Specified\"",
  install TEXT NOT NULL DEFAULT "echo \"No Install Command Specified\"",
  enabled BOOL NOT NULL DEFAULT 1,
*/

func Services() []*Service {
  services := make([]*Service, 0)
  rows := MustQuery(Instance, "SELECT * FROM services")

  for rows.Next() {
    svc := &Service{}
    var startArgs string
    var startEnv string 
    var installArgs string
    var installEnv string 
    rows.Scan(
      &svc.Id, 
      &svc.Name, 
      &svc.Desc, 
      &svc.Host, 
      &svc.Path, 
      &svc.Repo, 
      &svc.Start, 
      &startArgs, 
      &startEnv,
      &svc.Install, 
      &installArgs,
      &installEnv,
      &svc.Enabled,
    )

    var StartArgs []string
    var StartEnv []string

    var InstallArgs []string
    var InstallEnv []string

    json.Unmarshal([]byte(startArgs), &StartArgs)
    json.Unmarshal([]byte(startEnv), &StartEnv)

    json.Unmarshal([]byte(installArgs), &InstallArgs)
    json.Unmarshal([]byte(installEnv), &InstallEnv)


    svc.Args = StartArgs
    svc.Env = StartEnv

    svc.InstallArgs = InstallArgs
    svc.InstallEnv = InstallEnv

    services = append(services, svc)
  }

  rows.Close()

  return services
}

func init() {
  fmt.Println("Initializing Database")
  current, err := user.Current()

  if err != nil {
    panic(err)
  }

  data_path := path.Join(current.HomeDir, ".indigo")

  resetPtr := flag.Bool("reset", false, "Clear local data.")
  flag.Parse()

  reset := *resetPtr

  if reset {
    os.RemoveAll(data_path)
  }

  util.Mkdir(data_path)
  util.Hide(data_path)
  db_path := path.Join(data_path, "store.db")
  _ = sqlite3.SQLiteDriver{}
  db, err := sql.Open("sqlite3", db_path)
  table_setup_queries := string(assets.MustAsset("resources/setup.sql"))

  if err != nil {
    fmt.Println(err)
  } else {
    MustExec(db, table_setup_queries)
    Exec(db, "INSERT INTO services(name, host, path) VALUES('static', '*', '/')")
    Exec(db, "INSERT INTO services(name, host, path) VALUES('conversion', '*', '/conv')")
    Exec(db, "INSERT INTO services(name, host, path, enabled) VALUES('old_service', '*', '/old', 0)")
    db = db
  }

  Instance = db
  fmt.Println("Database Ready")
}