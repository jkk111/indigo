package services

import (
  "fmt"
  "testing"
  "github.com/jkk111/indigo/util"
  "github.com/jkk111/indigo/git"
  "github.com/jkk111/indigo/database"
  "github.com/jkk111/indigo/proxy"
)

const REPO = "https://github.com/jkk111/indigo"

func Setup() {
  util.BASE_NAME = ".indigo-test"
  if database.Instance != nil {
    database.Instance.Close()
  }

  fmt.Println("Exists", util.Exists(util.DataDir()))

  util.Rmdir(util.DataDir())
  util.Mkdir(util.DataDir())
  database.Setup()
  proxy.NewReverseProxy()
}

func TestInstallService(t * testing.T) {
  Setup()
  branches := git.LsRemote(REPO)
  branch := branches["master"]
  branch.Clone()

  service := &database.Service{
    Start: "indigo",
    Args: []string { },
    Env: []string { "TEST=true" },
    Install: "go",
    InstallArgs: []string { "build" },
  }

  InstallService(service, branch.Hash)
  TearDown()
}

func TestLoad(t * testing.T) {
  Setup()

  branches := git.LsRemote(REPO)
  branch := branches["master"]

  service := &database.Service{
    Name: "test-service",
    Desc: "test service",
    Host: "*", 
    Path: "/test-path", 
    Repo: REPO,
    Branch: "master",
    LatestHash: branch.Hash,
    Start: "indigo",
    Args: []string { },
    Env: []string { "TEST=true" },
    StartArgsRaw: "[]",
    StartEnvRaw: `[ "TEST=true" ]`,
    Install: "go",
    InstallArgs: []string { "build" },
    InstallArgsRaw: `[ "build" ]`,
    InstallEnvRaw: `[]`,
  }

  database.AddService(service)

  Load()
  TearDown()
}

func TearDown() {
  util.Rmdir(util.DataDir())
  if database.Instance != nil {
    database.Instance.Close()
  }
  Close()
}
