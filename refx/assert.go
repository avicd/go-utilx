package refx

import (
	"reflect"
)

func IsAnyOf(vl any, kinds ...reflect.Kind) bool {
	if len(kinds) > 0 {
		valKind := IndirectKind(vl)
		for _, kind := range kinds {
			switch kind {
			case reflect.Pointer:
				if IsPointer(vl) {
					return true
				}
			case reflect.Interface:
				if IsInterface(vl) {
					return true
				}
			default:
				if valKind == kind {
					return true
				}
			}
		}
	}
	return false
}

func IsInterface(vl any) bool {
	return KindOf(vl) == reflect.Interface
}

func IsPointer(vl any) bool {
	return KindOf(vl) == reflect.Pointer
}

func IsNil(vl any) bool {
	kind := KindOf(vl)
	switch kind {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return ValueOf(vl).IsNil()
	}
	return kind == reflect.Invalid
}

func IsBasic(vl any) bool {
	kind := IndirectKind(vl)
	switch kind {
	case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
		return false
	}
	return true
}

func IsError(vl any) bool {
	switch vl.(type) {
	case error:
		return true
	default:
		ref := TypeOf(func(err error) {})
		if TypeOf(vl) == ref.In(0) {
			return true
		}
	}
	return false
}

func IsList(vl any) bool {
	kind := IndirectKind(vl)
	switch kind {
	case reflect.Slice, reflect.Array:
		return true
	}
	return false
}

func IsInteger(vl any) bool {
	kind := IndirectKind(vl)
	return kind >= reflect.Int && kind <= reflect.Int64
}

func IsUInteger(vl any) bool {
	kind := IndirectKind(vl)
	return kind >= reflect.Uint && kind <= reflect.Uint64
}

func IsGeneralInt(vl any) bool {
	kind := IndirectKind(vl)
	return kind >= reflect.Int && kind <= reflect.Uint64
}

func IsFloat(vl any) bool {
	kind := IndirectKind(vl)
	return kind >= reflect.Float32 && kind <= reflect.Float64
}

func IsNumber(vl any) bool {
	kind := IndirectKind(vl)
	return kind >= reflect.Int && kind <= reflect.Uint64 ||
		kind >= reflect.Float32 && kind <= reflect.Float64
}

func IsInvalid(vl any) bool {
	return KindOf(vl) == reflect.Invalid
}

func IsBool(vl any) bool {
	return IndirectKind(vl) == reflect.Bool
}
func IsInt(vl any) bool {
	return IndirectKind(vl) == reflect.Int
}
func IsInt8(vl any) bool {
	return IndirectKind(vl) == reflect.Int8
}
func IsInt16(vl any) bool {
	return IndirectKind(vl) == reflect.Int16
}
func IsInt32(vl any) bool {
	return IndirectKind(vl) == reflect.Int32
}
func IsInt64(vl any) bool {
	return IndirectKind(vl) == reflect.Int64
}
func IsUint(vl any) bool {
	return IndirectKind(vl) == reflect.Uint
}
func IsUint8(vl any) bool {
	return IndirectKind(vl) == reflect.Uint8
}
func IsUint16(vl any) bool {
	return IndirectKind(vl) == reflect.Uint16
}
func IsUint32(vl any) bool {
	return IndirectKind(vl) == reflect.Uint32
}
func IsUint64(vl any) bool {
	return IndirectKind(vl) == reflect.Uint64
}
func IsUintptr(vl any) bool {
	return IndirectKind(vl) == reflect.Uintptr
}
func IsFloat32(vl any) bool {
	return IndirectKind(vl) == reflect.Float32
}
func IsFloat64(vl any) bool {
	return IndirectKind(vl) == reflect.Float64
}
func IsComplex64(vl any) bool {
	return IndirectKind(vl) == reflect.Complex64
}
func IsComplex128(vl any) bool {
	return IndirectKind(vl) == reflect.Complex128
}
func IsArray(vl any) bool {
	return IndirectKind(vl) == reflect.Array
}
func IsChan(vl any) bool {
	return IndirectKind(vl) == reflect.Chan
}
func IsFunc(vl any) bool {
	return IndirectKind(vl) == reflect.Func
}
func IsMap(vl any) bool {
	return IndirectKind(vl) == reflect.Map
}
func IsSlice(vl any) bool {
	return IndirectKind(vl) == reflect.Slice
}
func IsString(vl any) bool {
	return IndirectKind(vl) == reflect.String
}
func IsStruct(vl any) bool {
	return IndirectKind(vl) == reflect.Struct
}
func IsUnsafePointer(vl any) bool {
	return IndirectKind(vl) == reflect.UnsafePointer
}
