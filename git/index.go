package git

import (
  "encoding/json"
  "os/exec"
  "io"
  "io/ioutil"
  "strings"
  "regexp"
  "path"
  "github.com/jkk111/indigo/util"
)

const ref_prefix_len = 11

type Branch struct {
  Repo string `json:"repo"`
  Hash string `json:"hash"`
  Ref string `json:"ref"`
}

func (this * Branch) String() string {
  data, err := json.Marshal(&this)

  if err != nil {
    panic(err)
  }

  return string(data)
}

func run(args ...string) (string, error) {
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
  run_err := cmd.Wait()

  if read_err != nil {
    panic(read_err)
  }

  output := strings.TrimSpace(string(buf))

  if output == "" {
    panic("No Output")
  }

  return output, run_err
}

func (this * Branch) Branch() string {
  return this.Ref[ref_prefix_len:]
}

func (this * Branch) Clone() bool {
  util.Rmdir(util.Path(path.Join("repos", this.Hash)))
  _, err := run("clone", "-b", this.Branch(), this.Repo, util.Path(path.Join("repos", this.Hash)))

  if err != nil {
    return false
  }
  return true
}

func LsRemote(repo string) map[string]*Branch {
  if repo == "" {
    return nil
  }

  remote, err := run("ls-remote", "--heads", repo)

  if err != nil {
    return nil
  }

  remotes := strings.Split(remote, "\n")
  branches := make(map[string]*Branch, len(remotes))
  re := regexp.MustCompile(`\s+`)

  for _, branch := range remotes {
    parts := re.Split(branch, -1)

    if len(parts) != 2 {
      return nil
    }

    b := &Branch{ Repo: repo, Hash: parts[0], Ref: parts[1] }
    branches[b.Branch()] = b
  }

  return branches
}
