package services

import (
  "testing"
  "github.com/jkk111/indigo/util"
  "github.com/jkk111/indigo/git"
  "github.com/jkk111/indigo/database"
)

const REPO = "https://github.com/jkk111/indigo"

func Setup() {
  util.BASE_NAME = ".indigo-test"
}

func TestInstallServcice(t * testing.T) {
  Setup()
  branches := git.LsRemote(REPO)
  branch := branches["master"]
  branch.Clone()

  service := &database.Service{
    Start: "indigo",
    Args: []string { },
    Env: []string { "TEST=true" },
    Install: "go",
    InstallArgs: []string { "Install" },
  }

  InstallService(service, branch.Hash)
  TearDown()
}

func TearDown() {
  util.Rmdir(util.DataDir())
}
