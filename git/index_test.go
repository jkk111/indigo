package git

import (
  "fmt"
  "testing"
  "github.com/jkk111/indigo/util"
)

const SELF_REPO = "https://github.com/jkk111/indigo"
const INVALID_REPO = "https://github.com/jkk111/indigo-invalid/"

func TestLsRemote(t * testing.T) {
  branches := LsRemote(SELF_REPO)

  if branches == nil {
    t.Error("Expected Map got Nil")
  }

  if branches["master"] == nil {
    t.Error("Expected Map To Contain 'master' Branch")
  }
}

func TestInvalidLsRemote(t * testing.T) {
  branches := LsRemote(INVALID_REPO)

  if branches != nil {
    t.Error("Expected Map To Be Nil, Got", branches)
  }
}

func TestClone(t * testing.T) {
  util.BASE_NAME = ".indigo-test"
  // Setup + Cleanup stray Test
  util.Rmdir(util.DataDir())
  util.Mkdir(util.DataDir())
  util.Hide(util.DataDir())

  branches := LsRemote(SELF_REPO)
  success := branches["master"].Clone()

  if !success {
    t.Error("Clone Failed")
  }

  // Cleanup
  util.Rmdir(util.DataDir())
}

// Note: This test isn't great since it will fail if the format of branch is ever changed
// Will Require updating as changes happen.
func TestBranchString(t * testing.T) {
  branches := LsRemote(SELF_REPO)
  branch := branches["master"]
  expect := fmt.Sprintf(`{"repo":"%s","hash":"%s","ref":"%s"}`, branch.Repo, branch.Hash, branch.Ref)

  if branch.String() != expect {
    t.Error("Expected", expect, "Got", branch.String())
  }
}

func TestBranchBranch(t * testing.T) {
  branches := LsRemote(SELF_REPO)
  branch := branches["master"]

  if branch.Branch() != "master" {
    t.Error("Expected Branch 'master'", "Got", branch.Branch())
  }
}