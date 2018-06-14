package util

import (
  "testing"
  "os/user"
  "path/filepath"
)

const TEST_DIR = ".indigo-test-2"

func TestRandomId(t * testing.T) {
  id := RandomId()

  if len(id) != 36 {
    t.Error("Expected Id Of Length 36 Got:", len(id))
  }
}

func Setup() {
  BASE_NAME = TEST_DIR
}

func TearDown() {
  Rmdir(DataDir())
}

func TestMkdir(t * testing.T) {
  Setup()
  dir_path := Path("repos")
  Mkdir(dir_path)

  if !Exists(dir_path) {
    t.Error("Mkdir Failed")
  } else {
    Rmdir(dir_path)
  }
}

func TestRmdir(t * testing.T) {
  Setup()
  dir_path := Path("repos")
  Mkdir(dir_path)

  if !Exists(dir_path) {
    t.Error("Mkdir Failed")
  }

  Rmdir(dir_path)

  if Exists(dir_path) {
    t.Error("Rmdir Failed")
  }
}

func TestExists(t * testing.T) {
  Setup()
  dir_path := Path("repos")

  if Exists(dir_path) {
    Rmdir(dir_path)
    if Exists(dir_path) {
      t.Error("TestExists Setup Failed")
    }
  }

  Mkdir(dir_path)
  if !Exists(dir_path) {
    t.Error("Failed To Create Directory")
  }
}

func TestDataDir(t * testing.T) {
  Setup()

  current, err := user.Current()

  if err != nil {
    t.Error("Failed To Get User")
  }

  base_dir := current.HomeDir
  base_dir = filepath.Join(base_dir, TEST_DIR)

  if base_dir != DataDir() {
    t.Error("Expected", base_dir,
            "Got", DataDir(),
    )
  }

  TearDown()
}

func TestPath(t * testing.T) {
  Setup()

  current, err := user.Current()

  if err != nil {
    t.Error("Failed To Get User")
  }
  
  base_dir := current.HomeDir
  base_dir = filepath.Join(base_dir, TEST_DIR)
  test_dir := filepath.Join(base_dir, "test")

  if test_dir != Path("test") {
    t.Error("Expected", test_dir,
            "Got", Path("test"),
    )
  }

  TearDown()
}
