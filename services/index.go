package services

import (
  "runtime"
  "os/exec"
  "sync"
  "os"
  "io"
  "fmt"
  "regexp"
  "reflect"
  path "path/filepath"
  "github.com/jkk111/indigo/util"
  "github.com/jkk111/indigo/logger"
  "github.com/jkk111/indigo/database"
  "github.com/jkk111/indigo/git"
  Proxy "github.com/jkk111/indigo/proxy"
)

var GlobalRunningProcesses = make(map[string]*ActiveService)
var RunHistory = make(map[string]int64)
var RunningMutex = sync.Mutex{}
var RunCountMutex = sync.Mutex{}

type ActiveService struct {
  id string
  cmd * exec.Cmd
}

func must_pipe(p io.ReadCloser, err error) io.ReadCloser {
  if err != nil {
    panic(err)
  }

  return p
}

func InstallService(service * database.Service, hash string) {
  dir := util.Path(path.Join("repos", hash))

  cmd := exec.Command(service.Install, service.InstallArgs...)

  cmd_env := append(os.Environ(), service.InstallEnv...)
  cmd.Env = cmd_env
  cmd.Dir = dir

  stdout := must_pipe(cmd.StdoutPipe())
  stderr := must_pipe(cmd.StderrPipe())

  logger.CreateLogInstance(service.Name + "-install", stdout, stderr)
  err := cmd.Run()

  if err != nil {
    fmt.Println(err)
    panic(err)
  }
}

func NewActiveService(service_id string, commit string, start []string, env []string) * ActiveService {
  fmt.Println("Running", start[0], start[1:])
  cmd := exec.Command(start[0], start[1:]...)
  re := regexp.MustCompile(`\${datadir}`)
  datadir := util.DataDir()
  for i, v := range env {
    env[i] = re.ReplaceAllString(v, datadir)
  }

  cmd_env := append(os.Environ(), env...)

  cmd.Env = cmd_env
  cmd.Dir = util.Path(path.Join("repos", commit))

  svc := &ActiveService{ 
    id: util.RandomId(),
    cmd: cmd,
  }

  stdout, err := cmd.StdoutPipe()

  if err != nil {
    panic(err)
  }

  stderr, err := cmd.StderrPipe()

  if err != nil {
    panic(err)
  }

  logger.CreateLogInstance(service_id, stdout, stderr)
  RunningMutex.Lock()
  GlobalRunningProcesses[svc.id] = svc
  RunningMutex.Unlock()
  var start_err error
  start_err = cmd.Start()
  go func() {
    if start_err == nil {
      err := cmd.Wait()
      if err != nil {
        fmt.Println(reflect.TypeOf(err))
        fmt.Println(err)
      } else {
        RunningMutex.Lock()
        delete(GlobalRunningProcesses, svc.id)
        RunningMutex.Unlock()
      }
    }
  }()
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
    fmt.Printf("%+v\n", service)
    env := service.Env
    instance_no := RunHistory[service.Name]
    RunCountMutex.Lock()
    RunHistory[service.Name]++
    RunCountMutex.Unlock()
    ln_port := util.GetSocket(service.Name, int(instance_no))
    port := fmt.Sprintf("socket=%s", ln_port)
    env = append(env, port)

    params := make([]string, len(service.Args) + 1)
    params[0] = service.Start
    copy(params[1:], service.Args)

    branch := service.Branch
    branches := git.LsRemote(service.Repo)

    if branches == nil {
      continue
    }

    var b * git.Branch

    if branches[branch] == nil {
      b = branches["master"]
    } else {
      b = branches[branch]
    }

    if !util.Exists(util.Path(path.Join("repos", b.Hash))) {
      b.Clone()
      InstallService(service, b.Hash)
    }

    service.LatestHash = b.Hash

    NewActiveService(service.Name, service.LatestHash, params, env)
    proxy.AddRoute(service.Host, service.Path, ln_port, true)
  }
}

func New(service * database.Service) {
  git.LsRemote(service.Repo)
}