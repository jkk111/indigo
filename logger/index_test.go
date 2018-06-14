package logger

import (
  "testing"
  "github.com/jkk111/indigo/util"
)

type TestCloser struct {
  data []byte
  offset int
  closed bool
}

func NewTestCloser(buf []byte) * TestCloser {
  return &TestCloser{ data: buf }
}

func (this * TestCloser) Read(buf []byte) (read int, e error) {
  read = copy(buf, this.data[this.offset:])
  this.offset += read
  return
}

func (this * TestCloser) Close() error {
  this.closed = true
  return nil
}

func Setup() {
  util.BASE_NAME = "indigo-test"
}

func TestCreateLogInstance(t * testing.T) {
  buf_out := NewTestCloser([]byte("Hello World"))
  buf_err := NewTestCloser([]byte("Hello World"))
  CreateLogInstance("test", buf_out, buf_err)
}