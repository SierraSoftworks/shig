package core_test

import (
	"testing"
)

type testOutput struct {
	t *testing.T
}

func (o *testOutput) Printf(format string, a ...interface{}) {
	o.t.Logf(format, a...)
}

func (o *testOutput) Println(a ...interface{}) {
	o.t.Log(a...)
}
