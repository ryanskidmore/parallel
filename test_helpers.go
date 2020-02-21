package parallel

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"testing"
)

// Assert fails the test if the condition is false.
func testAssert(t *testing.T, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		t.FailNow()
	}
}

// ok fails the test if an err is not nil.
func testNil(t *testing.T, val interface{}) {
	if val != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected value: %s\033[39m\n\n", filepath.Base(file), line, val)
		t.FailNow()
	}
}

func testNotNil(t *testing.T, val interface{}) {
	if val == nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: expected not nil, got nil\033[39m\n\n", filepath.Base(file), line)
		t.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func testEquals(t *testing.T, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		t.FailNow()
	}
}

type testStruct struct {
	Counter int
	Mutex   sync.Mutex
}
