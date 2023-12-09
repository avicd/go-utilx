package refx

import (
	"unsafe"
)

var (
	TBool          bool
	TByte          byte
	TBytes         []byte
	TInt8          int8
	TInt16         int16
	TInt32         int32
	TInt64         int64
	TUint8         uint8
	TUint16        uint16
	TUint32        uint32
	TUint64        uint64
	TUintptr       uintptr
	TFloat32       float32
	TFloat64       float64
	TComplex64     complex64
	TComplex128    complex128
	TArray         [1]any
	TChan          chan any
	TFunc          func()
	TSlice         []any
	TString        []any
	TStruct        struct{}
	TUnsafePointer unsafe.Pointer
	TAny           any
	TMapStrAny     map[string]any
)
