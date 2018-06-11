package git

import (
  "os/exec"
  "io"
  "io/ioutil"
  "strings"
  "regexp"
  "fmt"
  "path"
  "github.com/jkk111/indigo/util"
)

const ref_prefix_len = 11

type Branch struct {
  repo string 
  hash string
  ref string
}

func run(args ...string) string {
  cmd := exec.Command("git", args...)
  out, p_err := cmd.StdoutPipe()
  err, p_err2 := cmd.StderrPipe()

  if p_err != nil {
    panic(p_err)
  }

  if p_err2 != nil {
    panic(p_err2)
  }

  cmd.Start()
  reader := io.MultiReader(out, err)
  buf, read_err := ioutil.ReadAll(reader)
  cmd.Wait()

  if read_err != nil {
    panic(read_err)
  }

  output := strings.TrimSpace(string(buf))

  if output == "" {
    panic("No Output")
  }

  return output
}

func (this Branch) Ref() string {
  return  this.ref
}

func (this Branch) Branch() string {
  return this.ref[ref_prefix_len:]
}

func (this Branch) Hash() string {
  return this.hash
}

func (this Branch) Clone() {
  util.Rmdir(util.Path(path.Join("repos", this.hash)))
  run("clone", "-b", this.Branch(), this.repo, util.Path(path.Join("repos", this.hash)))
}

func LsRemote(repo string) map[string]Branch {
  remote := run("ls-remote", "--heads", repo)
  remotes := strings.Split(remote, "\n")
  branches := make(map[string]Branch, len(remotes))
  re := regexp.MustCompile(`\s+`)

  for _, branch := range remotes {
    parts := re.Split(branch, -1)
    b := Branch{ repo: repo, hash: parts[0], ref: parts[1] }
    branches[b.Branch()] = b
  }

  return branches
}

func Remotes(repo string) []string {
  str := run("ls-remote", "--heads", repo)
  branches := strings.Split(str, "\n")
  return branches
}

func BranchHash(repo string, branch string) string {
  branch_pattern := fmt.Sprintf("refs/heads/%s", branch)
  branches := Remotes(repo)

  re := regexp.MustCompile(`\s+`)
  for _, branch := range branches {
    parts := re.Split(branch, -1)

    if parts[1] == branch_pattern {
      return parts[0] 
    }
  }

  return ""
}

func Branches(repo string) []string {
  branches := Remotes(repo)
  re := regexp.MustCompile(`\s+`)

  branch_names := make([]string, 0)

  for _, branch := range branches {
    parts := re.Split(branch, -1)
    ref := parts[1]
    i := strings.LastIndex(ref, "/")
    branch_names = append(branch_names, ref[i + 1:])
  }

  return branch_names
}