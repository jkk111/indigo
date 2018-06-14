package database

import (
  "fmt"
  "database/sql"
  "encoding/json"
  "github.com/jmoiron/sqlx"
  "github.com/mattn/go-sqlite3"
  "github.com/jkk111/indigo/assets"
  "github.com/jkk111/indigo/util"
)

var Instance * sqlx.DB

type Service struct {
  Id int `json:"id" db:"id"`
  Name string `json:"name" db:"name"`
  Desc string `json:"desc" db:"desc"`
  Host string `json:"host" db:"host"`
  Path string `json:"path" db:"path"`
  Repo string `json:"repo" db:"repo"`
  Branch string `json:"branch" db:"branch"`
  LatestHash string `json:"hash" db:"hash"`

  Start string `json:"start" db:"start"`

  StartArgsRaw string `json:"-" db:"args"`
  Args []string `json:"args" db:"-"`
  StartEnvRaw string `json:"-" db:"env"`
  Env []string `json:"env" db:"-"`

  Install string `json:"install" db:"install"`

  InstallArgsRaw string `json:"-" db:"installArgs"`
  InstallArgs []string `json:"installArgs" db:"-"`
  InstallEnvRaw string `json:"-" db:"installEnv"`
  InstallEnv []string `json:"installEnv" db:"-"`

  Enabled bool `json:"enabled" db:"enabled"`
}

func Exec(db * sqlx.DB, q string) (sql.Result, error) {
  return db.Exec(q)
}

func MustExec(db * sqlx.DB, q string) {
  _, err := Exec(db, q)
  if err != nil {
    panic(err)
  }
}

func Query(db * sqlx.DB, q string) (* sqlx.Rows, error) {
  return db.Queryx(q)
}

func MustQuery(db * sqlx.DB, q string) * sqlx.Rows {
  rows, err := Query(db, q)

  if err != nil {
    panic(err)
  }

  return rows
}

func AddService(service * Service) (sql.Result, error)  {
  sql := string(assets.MustAsset("resources/add_service.sql"))
  return Instance.NamedExec(sql, service)
}

func UpdateService(service * Service) (sql.Result, error) {
  sql := string(assets.MustAsset("resources/update_service.sql"))
  return Instance.NamedExec(sql, service)
}

func Must(err error) {
  if err != nil {
    panic(err)
  }
}

func Services() []*Service {
  services := make([]*Service, 0)
  rows := MustQuery(Instance, "SELECT * FROM services")

  for rows.Next() {
    svc := &Service{}

    err := rows.StructScan(
      svc,
    )

    if err != nil {
      panic(err)
    }

    Must(json.Unmarshal([]byte(svc.StartArgsRaw), &svc.Args))
    Must(json.Unmarshal([]byte(svc.StartEnvRaw), &svc.Env))
    Must(json.Unmarshal([]byte(svc.InstallArgsRaw), &svc.InstallArgs))
    Must(json.Unmarshal([]byte(svc.InstallEnvRaw), &svc.InstallEnv))

    services = append(services, svc)
  }

  rows.Close()

  return services
}

func Setup() {
  fmt.Println("Initializing Database")

  if Instance != nil {
    Instance.Close()
  }

  db_path := util.Path("store.db")
 
  _ = sqlite3.SQLiteDriver{}

  db, err := sqlx.Connect("sqlite3", db_path)
  table_setup_queries := string(assets.MustAsset("resources/setup.sql"))

  if err != nil {
    fmt.Println(err, db_path)
    panic(err)
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

func init() {
  Setup()
}