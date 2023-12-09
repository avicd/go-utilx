package refx

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type testInterface interface {
	TestIs()
}

type testAssert struct {
}

func (ta *testAssert) TestIs() {

}

func TestIs(t *testing.T) {
	p0 := &testAssert{}
	mp := map[string]any{}
	var p1 any = p0

	// IsAnyOf
	assert.Equal(t, false, IsAnyOf(p0, reflect.Interface))
	assert.Equal(t, true, IsAnyOf(p0, reflect.Interface, reflect.Pointer))
	assert.Equal(t, true, IsAnyOf(p0, reflect.Map, reflect.Struct))

	// IsInterface
	assert.Equal(t, true, IsInterface(reflect.TypeOf(mp).Elem()))

	// IsPointer
	assert.Equal(t, true, IsPointer(p0))
	assert.Equal(t, true, IsPointer(p1))

	// IsBasic
	assert.Equal(t, false, IsBasic(p1))
	assert.Equal(t, true, IsBasic(0))
	assert.Equal(t, true, IsBasic("0"))

	var p11 *testOp
	// IsNil
	assert.Equal(t, true, IsNil(nil))
	assert.Equal(t, true, IsNil(p11))
	assert.Equal(t, false, IsNil(p0))

	// IsList
	list := make([]int, 10)
	assert.Equal(t, true, IsList([]string{}))
	assert.Equal(t, false, IsList(false))
	assert.Equal(t, false, IsList(102))
	assert.Equal(t, true, IsList(list))

	//Integer
	assert.Equal(t, true, IsInteger(10))
	assert.Equal(t, true, IsUInteger(uint(1)))
	assert.Equal(t, false, IsUInteger(uintptr(1)))
	assert.Equal(t, false, IsGeneralInt(uintptr(1)))

	// Float
	assert.Equal(t, false, IsFloat(uintptr(1)))
	assert.Equal(t, false, IsFloat(int64(1)))
	assert.Equal(t, false, IsFloat(uint64(1)))
	assert.Equal(t, true, IsFloat(10.2))

	// IsNumber
	assert.Equal(t, true, IsNumber(10.2))
	assert.Equal(t, true, IsNumber(uint64(1)))
	assert.Equal(t, true, IsNumber(200))

	// reflect Kind
	assert.Equal(t, true, IsInvalid(nil))
	assert.Equal(t, true, IsInvalid(TypeOf(nil)))
	assert.Equal(t, false, IsInvalid(0))

	assert.Equal(t, true, IsBool(true))
	assert.Equal(t, true, IsBool(false))
	assert.Equal(t, false, IsBool(nil))
	assert.Equal(t, false, IsBool(0))

	assert.Equal(t, true, IsInt(200))
	assert.Equal(t, false, IsInt("2"))

	assert.Equal(t, true, IsInt8(int8(10)))
	assert.Equal(t, true, IsInt8(int8(10)))

	assert.Equal(t, true, IsInt16(int16(10)))
	assert.Equal(t, true, IsInt16(int16(10)))

	assert.Equal(t, false, IsInt32(int16(10)))
	assert.Equal(t, true, IsInt32(int32(10)))

	assert.Equal(t, false, IsInt64(int16(10)))
	assert.Equal(t, true, IsInt64(int64(21)))

	assert.Equal(t, false, IsUint(uintptr(10)))
	assert.Equal(t, true, IsUint(uint(21)))

	assert.Equal(t, false, IsUint8(uint32(10)))
	assert.Equal(t, true, IsUint8(uint8(21)))

	assert.Equal(t, false, IsUint16(uint32(10)))
	assert.Equal(t, true, IsUint16(uint16(21)))

	assert.Equal(t, false, IsUint32(uint64(10)))
	assert.Equal(t, true, IsUint32(uint32(21)))

	assert.Equal(t, false, IsUint64(uint(10)))
	assert.Equal(t, true, IsUint64(uint64(21)))

	assert.Equal(t, false, IsUintptr(uint(10)))
	assert.Equal(t, true, IsUintptr(uintptr(21)))

	assert.Equal(t, false, IsUintptr(uint(10)))
	assert.Equal(t, true, IsUintptr(uintptr(21)))

	assert.Equal(t, false, IsFloat32(2.54))
	assert.Equal(t, true, IsFloat32(float32(21)))

	assert.Equal(t, true, IsFloat64(2.54))
	assert.Equal(t, false, IsFloat64(float32(2.231)))

	assert.Equal(t, false, IsComplex64(2.54))
	assert.Equal(t, true, IsComplex64(complex64(1+1i)))

	assert.Equal(t, false, IsComplex128(2.54))
	assert.Equal(t, true, IsComplex128(2.231+2i))

	var buf [10]string
	assert.Equal(t, true, IsArray(buf))
	assert.Equal(t, false, IsArray(2.231+2i))

	channel := make(chan int, 1)
	assert.Equal(t, true, IsChan(channel))
	assert.Equal(t, false, IsChan(2.231+2i))

	var fn func()
	assert.Equal(t, true, IsFunc(fn))
	assert.Equal(t, false, IsFunc(2.231))

	assert.Equal(t, true, IsMap(map[string]any{}))
	assert.Equal(t, false, IsMap(2.231))

	assert.Equal(t, true, IsSlice(make([]int, 10)))
	assert.Equal(t, false, IsSlice(2.231))

	assert.Equal(t, true, IsString(""))
	assert.Equal(t, true, IsString("9999"))
	assert.Equal(t, false, IsString(nil))
	assert.Equal(t, false, IsString(2.231))

	assert.Equal(t, true, IsStruct(&testAssert{}))
	assert.Equal(t, true, IsStruct(testAssert{}))
	assert.Equal(t, false, IsStruct(map[int]any{}))
	assert.Equal(t, false, IsStruct(nil))
	assert.Equal(t, false, IsStruct(2.231))

	assert.Equal(t, true, IsUnsafePointer(reflect.ValueOf(&testAssert{}).UnsafePointer()))
	assert.Equal(t, false, IsUnsafePointer(2.231))

}

type testErr struct {
}

func (terr testErr) Error() string {
	return ""
}

func TestIsError(t *testing.T) {
	var p func() (int, error)
	tp := TypeOf(p)
	println(IsInterface(tp.Out(1)))
	fmt.Printf("%v", NewOf(tp.Out(1)).Interface())
	assert.Equal(t, true, IsError(tp.Out(1)))
	assert.Equal(t, true, IsError(testErr{}))
	assert.Equal(t, true, IsError(errors.New("")))
}
