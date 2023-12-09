package goid

import (
	"github.com/avicd/go-utilx/logx"
	"reflect"
	"unsafe"
)

var gt reflect.Type

// return the pointer of runtime.g
func getgp() unsafe.Pointer

var goidOffset uintptr

func init() {
	gt = typeByString("runtime.g")
	if gt == nil {
		logx.Fatal("Failed to get type of runtime.g natively.")
	}
	if tf, ok := gt.FieldByName("goid"); ok {
		goidOffset = tf.Offset
	}
}

func Id() int64 {
	gp := getgp()
	if gp == nil {
		logx.Fatal("Failed to get gp from runtime natively.")
	}
	return *(*int64)(add(gp, goidOffset))
}

// Should be a built-in for unsafe.Pointer?
//
//go:nosplit
func add(p unsafe.Pointer, x uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}
