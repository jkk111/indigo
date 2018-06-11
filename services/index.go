package services

import (
  "runtime"
  "os/exec"
  "os"
  "fmt"
  "strings"
  "regexp"
  "github.com/jkk111/indigo/util"
  "github.com/jkk111/indigo/logger"
  "github.com/jkk111/indigo/database"
  "github.com/jkk111/indigo/git"
  Proxy "github.com/jkk111/indigo/proxy"
)

var GlobalRunningProcesses = make(map[string]*ActiveService)
var RunHistory = make(map[string]int64)

type ActiveService struct {
  id string
  cmd * exec.Cmd
}

func ParseCommitList(input string, ref string) (string, bool) {
  fmt.Println(input)
  input = strings.TrimSpace(input)


  if input == "" {
    return "", false
  }

  hashes := strings.Split(input, "\n")
  re := regexp.MustCompile(`\s+`)

  fmt.Println(hashes)

  for _, str := range hashes {
    parts := re.Split(str, -1)

    if parts[1] == ref {
      return parts[0], true
    }
  }

  return "", false
}

// Represents a specific clone of a repo!
func RepoId(repo string) {
  hash := git.BranchHash(repo, "master")

  if hash == "" {
    panic("Invalid Repo/Branch")
  }

  fmt.Println(hash)
}

func NewActiveService(service_id string, start []string, env []string) * ActiveService {
  cmd := exec.Command(start[0], start[1:]...)

  cmd_env := append(os.Environ(), env...)

  cmd.Env = cmd_env

  svc := &ActiveService{ 
    id: util.RandomId(),
    cmd: cmd,
  }

  GlobalRunningProcesses[svc.id] = svc

  stdout, err := cmd.StdoutPipe()

  if err != nil {
    panic(err)
  }

  stderr, err := cmd.StderrPipe()

  if err != nil {
    panic(err)
  }

  logger.CreateLogInstance(service_id, stdout, stderr)
  cmd.Start()

  return svc
}

func (this * ActiveService) Kill() {
  if runtime.GOOS == "windows" {
    this.cmd.Process.Kill()
  } else {
    this.cmd.Process.Signal(os.Interrupt)
  }
}

func Close() {
  for _, svc := range GlobalRunningProcesses {
    svc.Kill()
  }
}

func Load() {
  services := database.Services()
  proxy := Proxy.Instance()
  for _, service := range services {
    env := make([]string, 0)
    instance_no := RunHistory[service.Name]
    RunHistory[service.Name]++
    ln_port := util.GetSocket(service.Name, int(instance_no))
    port := fmt.Sprintf("PORT=%s", ln_port)
    env = append(env, port)

    instance_name := fmt.Sprintf("%s-%d", service.Name, instance_no)
    NewActiveService(service.Name, []string { service.Start }, env)

    fmt.Printf("Service: %+v\n", service, instance_name)

    proxy.AddRoute(service.Host, service.Path, ln_port, true)

    // services.NewActiveService("app", []string {"node", "-e", "console.log('go is great');" })
    // services.NewActiveService("app2", []string {"node", "-e", "console.log('node is great');" })
    // services.NewActiveService("app3", []string {"node", "-e", "console.log('python is great');" })
    // services.NewActiveService("app4", []string {"node", "-e", "console.log('java is not so great');" })
  }
}

func New(name string, repo string) {

}