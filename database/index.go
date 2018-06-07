package database

import (
  "fmt"
  "path"
  "os/user"
  "database/sql"
  "github.com/mattn/go-sqlite3"
  "github.com/jkk111/indigo/assets"
)

var Instance * sql.DB

type Service struct {
  Id int `json:"id"`
  Name string `json:"name"`
  Desc string `json:"desc"`
  Host string `json:"host"`
  Path string `json:"path"`
  Repo string `json:"repo"`
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

func Services() []*Service {
  services := make([]*Service, 0)
  rows := MustQuery(Instance, "SELECT * FROM services")

  for rows.Next() {
    svc := &Service{}
    rows.Scan(&svc.Id, &svc.Name, &svc.Desc, &svc.Host, &svc.Path, &svc.Repo, &svc.Enabled)
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