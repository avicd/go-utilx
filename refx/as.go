package refx

import (
	"fmt"
	"reflect"
)

func AsBool(val any) bool {
	var ret bool
	Assign(&ret, val)
	return ret
}

func AsInt(val any) int {
	var ret int
	Assign(&ret, val)
	return ret
}

func AsInt8(val any) int8 {
	var ret int8
	Assign(&ret, val)
	return ret
}

func AsInt16(val any) int16 {
	var ret int16
	Assign(&ret, val)
	return ret
}

func AsInt32(val any) int32 {
	var ret int32
	Assign(&ret, val)
	return ret
}

func AsInt64(val any) int64 {
	var ret int64
	Assign(&ret, val)
	return ret
}

func AsUint(val any) uint {
	var ret uint
	Assign(&ret, val)
	return ret
}

func AsUint8(val any) uint8 {
	var ret uint8
	Assign(&ret, val)
	return ret
}

func AsUint16(val any) uint16 {
	var ret uint16
	Assign(&ret, val)
	return ret
}

func AsUint32(val any) uint32 {
	var ret uint32
	Assign(&ret, val)
	return ret
}

func AsUint64(val any) uint64 {
	var ret uint64
	Assign(&ret, val)
	return ret
}

func AsFloat32(val any) float32 {
	var ret float32
	Assign(&ret, val)
	return ret
}

func AsFloat64(val any) float64 {
	var ret float64
	Assign(&ret, val)
	return ret
}

func AsString(val any) string {
	var ret string
	Assign(&ret, val)
	return ret
}

func AsListOf[T any](values ...T) []T {
	return values
}

func AsList(args ...any) []any {
	var dest []any
	for _, arg := range args {
		vl := ValueOf(arg)
		if IsList(vl) {
			for i := 0; i < vl.Len(); i++ {
				dest = append(dest, vl.Index(i).Interface())
			}
		} else {
			dest = append(dest, vl.Interface())
		}
	}
	return dest
}

func AsOf(ref any, val any) reflect.Value {
	if IsInterface(ref) {
		return ValueOf(val)
	}
	dest := NewOf(ref)
	Assign(dest.Addr(), val)
	return dest
}

func AsError(val any) error {
	if val == nil {
		return nil
	}
	switch tmp := val.(type) {
	case error:
		return tmp
	}
	return fmt.Errorf("%v", val)
}
