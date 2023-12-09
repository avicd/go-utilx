package refx

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type testOp struct {
	A string
	B int
	C map[string]any
	D []string
	E *testOp
}

func TestValueOf(t *testing.T) {
	v1 := reflect.ValueOf(0)
	assert.Equal(t, true, v1 == ValueOf(v1))
	assert.Equal(t, true, v1 == ValueOf(&v1))
}

func TestKindOf(t *testing.T) {
	var p *testOp
	var v1 = &p
	var v2 = ValueOf(p)
	var v3 = TypeOf(p)

	// Wrap Test
	assert.Equal(t, reflect.Pointer, KindOf(p))
	assert.Equal(t, reflect.Pointer, KindOf(v1))
	assert.Equal(t, reflect.Pointer, KindOf(v2))
	assert.Equal(t, reflect.Pointer, KindOf(v3))

	// Basic
	assert.Equal(t, reflect.Bool, KindOf(TBool))
	assert.Equal(t, reflect.Uint8, KindOf(TByte))
	assert.Equal(t, reflect.Int8, KindOf(TInt8))
	assert.Equal(t, reflect.Int16, KindOf(TInt16))
	assert.Equal(t, reflect.Int32, KindOf(TInt32))
	assert.Equal(t, reflect.Int64, KindOf(TInt64))
	assert.Equal(t, reflect.Uint8, KindOf(TUint8))
	assert.Equal(t, reflect.Uint16, KindOf(TUint16))
	assert.Equal(t, reflect.Uint32, KindOf(TUint32))
	assert.Equal(t, reflect.Uint64, KindOf(TUint64))
	assert.Equal(t, reflect.Uintptr, KindOf(TUintptr))
	assert.Equal(t, reflect.Float32, KindOf(TFloat32))
	assert.Equal(t, reflect.Float64, KindOf(TFloat64))
	assert.Equal(t, reflect.Complex64, KindOf(TComplex64))
	assert.Equal(t, reflect.Complex128, KindOf(TComplex128))
	assert.Equal(t, reflect.Array, KindOf(TArray))
	assert.Equal(t, reflect.Chan, KindOf(TChan))
	assert.Equal(t, reflect.Func, KindOf(TFunc))
	assert.Equal(t, reflect.Slice, KindOf(TSlice))
	assert.Equal(t, reflect.Struct, KindOf(TStruct))
	assert.Equal(t, reflect.UnsafePointer, KindOf(TUnsafePointer))
	assert.Equal(t, reflect.Invalid, KindOf(TAny))
}

func TestIndirectKind(t *testing.T) {
	var v1 *testOp
	var v2 *[]string
	var v3 *[1]string
	assert.Equal(t, reflect.Struct, IndirectKind(v1))
	assert.Equal(t, reflect.Slice, IndirectKind(v2))
	assert.Equal(t, reflect.Array, IndirectKind(v3))

	v4 := &v1
	v5 := &v2
	v6 := &v3
	assert.Equal(t, reflect.Struct, IndirectKind(v4))
	assert.Equal(t, reflect.Slice, IndirectKind(v5))
	assert.Equal(t, reflect.Array, IndirectKind(v6))

	v1 = &testOp{E: &testOp{A: "A"}}
	v2 = &[]string{"A", "B", "C"}
	v3 = &[1]string{"A"}

	assert.Equal(t, reflect.Struct, IndirectKind(v1))
	assert.Equal(t, reflect.Struct, IndirectKind(v1.E))
	assert.Equal(t, reflect.Slice, IndirectKind(v2))
	assert.Equal(t, reflect.Array, IndirectKind(v3))
}

func TestIndirect(t *testing.T) {
	p0 := "A"
	p1 := &p0
	p2 := &p1
	p3 := &p2
	assert.Equal(t, "A", Indirect(p1).Interface())
	assert.Equal(t, "A", Indirect(p2).Interface())
	assert.Equal(t, "A", Indirect(p3).Interface())

	v1 := &testOp{A: "A"}
	v2 := &v1
	v3 := &v2
	assert.Equal(t, Indirect(v1).UnsafeAddr(), Indirect(v2).UnsafeAddr())
	assert.Equal(t, Indirect(v2).UnsafeAddr(), Indirect(v3).UnsafeAddr())
}

func TestIndirectType(t *testing.T) {
	var p0 testOp
	p1 := &p0
	p2 := &p1
	p3 := &p2

	assert.Equal(t, reflect.TypeOf(p0), IndirectType(p1))
	assert.Equal(t, reflect.TypeOf(p0), IndirectType(p2))
	assert.Equal(t, reflect.TypeOf(p0), IndirectType(p3))
}

func TestNewOf(t *testing.T) {
	p0 := 10
	v1 := NewOf(p0)
	v1.Set(ValueOf(p0))
	assert.Equal(t, true, v1.Interface() == 10)
	assert.Equal(t, reflect.Int, v1.Kind())

	var list []string
	v2 := NewOf(list)
	v2.Set(ValueOf([]string{"A1", "AN"}))
	if list2, ok := v2.Interface().([]string); ok {
		assert.Equal(t, "A1", list2[0])
		assert.Equal(t, "AN", list2[1])
	} else {
		panic(errors.New("new of Slice failed"))
	}

	var s0 *testOp
	v3 := NewOf(s0)
	v3.Set(ValueOf(&testOp{A: "A1", B: 22}))
	if s1, ok := v3.Interface().(*testOp); ok {
		assert.Equal(t, "A1", s1.A)
		assert.Equal(t, 22, s1.B)
	} else {
		panic(errors.New("new of Struct failed"))
	}
}

func TestAssign(t *testing.T) {
	// String
	var str string

	Assign(&str, 10)
	assert.Equal(t, "10", str)

	Assign(&str, "A1")
	assert.Equal(t, "A1", str)

	Assign(&str, nil)
	assert.Equal(t, "", str)

	// Int
	var num int
	Assign(&num, 10)
	assert.Equal(t, 10, num)

	Assign(&num, "10")
	assert.Equal(t, 10, num)

	Assign(&num, "203.312")
	assert.Equal(t, 203, num)

	// Float
	var decimal float64
	Assign(&decimal, 10)
	assert.Equal(t, 10.0, decimal)

	Assign(&decimal, "10")
	assert.Equal(t, 10.0, decimal)

	Assign(&decimal, "203.312")
	assert.Equal(t, 203.312, decimal)

	// Struct
	var p0 *testOp

	Assign(&p0, &testOp{A: "A", B: 203})
	assert.Equal(t, "A", p0.A)
	assert.Equal(t, 203, p0.B)

	// Auto Addr
	p0 = nil
	Assign(&p0, testOp{A: "A", B: 203})
	assert.Equal(t, "A", p0.A)
	assert.Equal(t, 203, p0.B)

	// Map

	var map0 map[string]any
	Assign(&map0, map[string]any{"A": "A"})
	assert.Equal(t, "A", map0["A"])

}

func TestClone(t *testing.T) {
	// basic
	var str string
	s1 := "A<-B"
	Clone(&str, s1)
	assert.Equal(t, s1, str)

	var num int64
	s2 := int64(200)
	Clone(&num, s2)
	assert.Equal(t, s2, num)

	// Struct
	var p0 *testOp
	s3 := &testOp{A: "AAAA", B: 10, C: map[string]any{"C1": 11}, D: []string{"L0"}}
	p1 := &testOp{A: "AAAA"}
	p2 := s3
	Clone(&p0, s3)
	assert.Equal(t, false, p0 == s3)
	assert.Equal(t, true, p2 == s3)
	assert.Equal(t, true, reflect.DeepEqual(p0, p2))
	assert.Equal(t, true, reflect.DeepEqual(p0, p2))
	assert.Equal(t, false, reflect.DeepEqual(p1, p2))

	// Slice
	var s4 []string
	s4 = append(s4, "A1")
	s4 = append(s4, "A1")
	s4 = append(s4, "A1")

	var list []string
	Clone(&list, s4)
	assert.Equal(t, true, reflect.DeepEqual(s4, list))

	// Map
	s5 := map[string]any{
		"Names": []string{"Allen", "FlashMan"},
		"Age":   78,
	}
	var map0 map[string]any
	Clone(&map0, s5)
	assert.Equal(t, true, reflect.DeepEqual(s5, map0))

}

func TestMerge(t *testing.T) {
	// Cover
	p0 := &testOp{A: "", B: 10, C: map[string]any{"C1": 11}}
	p1 := &testOp{A: "A", B: 20, C: map[string]any{"C1": 21}}
	Merge(p0, p1)
	refP0 := &testOp{A: "A", B: 20, C: map[string]any{"C1": 21}}
	assert.Equal(t, true, reflect.DeepEqual(p0, refP0))

	// Ignore Zero Value When Covering
	p0 = &testOp{A: "A", B: 20, C: map[string]any{"C1": 21}}
	p1 = &testOp{A: "", B: 0, C: map[string]any{"C1": ""}}
	Merge(p0, p1)
	refP0 = &testOp{A: "A", B: 20, C: map[string]any{"C1": 21}}
	assert.Equal(t, true, reflect.DeepEqual(p0, refP0))

	// Create new When Dest is nil
	p0 = nil
	p1 = &testOp{C: map[string]any{"C1": "21", "C2": "C2"}, D: []string{"D3", "D4"}}

	Merge(&p0, p1)
	assert.Equal(t, true, reflect.DeepEqual(p0, p1))

	// Map
	p0 = &testOp{C: map[string]any{"C1": 21}, D: []string{"D1", "D2"}}
	p1 = &testOp{C: map[string]any{"C1": "", "C2": "C2"}, D: []string{"D3", "D4"}}
	Merge(p0, p1)
	refP0 = &testOp{C: map[string]any{"C1": 21, "C2": "C2"}, D: []string{"D1", "D2", "D3", "D4"}}
	assert.Equal(t, true, reflect.DeepEqual(p0, refP0))

	// Slice
	list0 := []string{"D1", "D2"}
	list1 := []string{"D3", "D4"}
	refList := []string{"D1", "D2", "D3", "D4"}
	Merge(&list0, list1)
	assert.Equal(t, true, reflect.DeepEqual(list0, refList))

	// Slice->Array
	arr0 := [3]string{"D1", "D2"}
	list1 = []string{"D3", "D4"}

	refArr := [3]string{"D3", "D4", ""}
	Merge(&arr0, list1)
	assert.Equal(t, true, reflect.DeepEqual(arr0, refArr))

	// Basic
	var str string
	var num float64

	Merge(&str, "STR0")
	Merge(&num, 201.232)
	assert.Equal(t, "STR0", str)
	assert.Equal(t, 201.232, num)
}

func TestFieldOf(t *testing.T) {
	list0 := []*testOp{
		{A: "A-Field", B: 10, C: map[string]any{"C1": 11}},
		{C: map[string]any{"C1": "", "C2": "C2"}, D: []string{"D3", "D4"}},
	}
	v0, exist := FieldOf(list0, 0, "A")
	assert.Equal(t, true, exist)
	assert.Equal(t, "A-Field", v0.Interface())

	v0, exist = FieldOf(list0, 0, "C", "C1")
	assert.Equal(t, true, exist)
	assert.Equal(t, 11, v0.Interface())
}

func TestTypeOfField(t *testing.T) {
	list0 := []*testOp{
		{A: "A-Field", B: 10, C: map[string]any{"C1": 11}},
		{C: map[string]any{"C1": "", "C2": "C2"}, D: []string{"D3", "D4"}},
	}
	v0, exist := TypeOfField(list0, 0, "B")
	assert.Equal(t, true, exist)
	assert.Equal(t, reflect.Int, v0.Kind())

	v0, exist = TypeOfField(list0, 1, "C", "C2")
	assert.Equal(t, true, exist)
	assert.Equal(t, reflect.String, v0.Kind())

	var p *testOp
	v0, exist = TypeOfField(p, "A")
	assert.Equal(t, true, exist)
	assert.Equal(t, reflect.String, v0.Kind())

	v0, exist = TypeOfField(p, "E", "E", "E")
	assert.Equal(t, true, exist)
	assert.Equal(t, reflect.Pointer, v0.Kind())
}

func TestPropOf(t *testing.T) {
	list0 := []*testOp{
		{A: "A-Field", B: 10, C: map[string]any{"C1": 11}},
		{C: map[string]any{"C1": "", "C2": "C2"}, D: []string{"D3", "D4"}},
	}

	v0, exist := PropOf(list0, 0, "A")
	assert.Equal(t, true, exist)
	assert.Equal(t, "A-Field", v0)

	v0, exist = PropOf(list0, 0, "C", "C1")
	assert.Equal(t, true, exist)
	assert.Equal(t, 11, v0)

	map0 := map[string]any{
		"10": "Value A",
	}
	v0, exist = PropOf(map0, "10")
	assert.Equal(t, true, exist)
	assert.Equal(t, "Value A", v0)

	v0, exist = PropOf(map0, "21")
	assert.Equal(t, false, exist)
	assert.Equal(t, nil, v0)
}

func TestSet(t *testing.T) {
	n0 := 1
	n1 := 20
	Set(&n0, n1)
	assert.Equal(t, 20, n0)

	p0 := &testOp{A: "A-Field", B: 10, C: map[string]any{"C1": 11}}
	Set(p0, "A-Field-Changed", "A")
	Set(p0, "MAP:D", "C", "D")
	assert.Equal(t, "A-Field-Changed", p0.A)
	assert.Equal(t, "MAP:D", p0.C["D"])

	p0 = nil
	Set(&p0, "MAP:D", "C", "D")
	assert.Equal(t, "MAP:D", p0.C["D"])

	var list0 [4]string
	ok := Set(&list0, "A", 2)
	assert.Equal(t, true, ok)
}
