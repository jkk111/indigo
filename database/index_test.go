package database

import (
  "testing"
  "github.com/jkk111/indigo/util"
)

func SetupTest() {
  if Instance != nil {
    Instance.Close()
  }
  util.BASE_NAME = ".indigo-test"
  util.Rmdir(util.DataDir())
  util.Mkdir(util.DataDir())
  Setup()
}

func TearDown() {
  Instance.Close()
  util.Rmdir(util.DataDir())
}

func TestServices(t * testing.T) {
  SetupTest()
  services := Services()

  if len(services) != 3 {
    t.Error("Expected 3 Services GOT", len(services))
  }

  TearDown()
}

func TestAddService(t * testing.T) {
  SetupTest()
  services := Services()

  if len(services) != 3 {
    t.Error("Expected 3 Services GOT", len(services))
  }

  next := &Service{ 
    Name: "Test Name", 
    Path: "Test Path", 
    StartArgsRaw: "[]",
    StartEnvRaw: "[]",
    InstallArgsRaw: "[]",
    InstallEnvRaw: "[]",
    Enabled: true,
  }

  result, err := AddService(next)

  if err != nil {
    t.Error("Error Adding Service")
  }

  insert_id, err := result.LastInsertId()

  if err != nil {
    t.Error("Failed To Retrieve Insert Id")
  }

  if insert_id != 4 {
    t.Error("Invalid Insert Id")
  }

  TearDown()
}

func TestUpdateService(t * testing.T) {
  SetupTest()
  services := Services()

  if len(services) != 3 {
    t.Error("Expected 3 Services GOT", len(services))
  }

  next := services[0]

  next.Name = "Updated Test"

  result, err := UpdateService(next)

  if err != nil {
    t.Error("Error Updating Service")
  }

  affected, err := result.RowsAffected()

  if err != nil {
    t.Error("Failed To Retrieve Rows Affected")
  }

  if affected != 1 {
    t.Error("Invalid Number Of Rows Affected")
  }

  TearDown()
}