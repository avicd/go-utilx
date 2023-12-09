package refx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testAsSt struct {
	Str string
}

func TestAs(t *testing.T) {
	var p *testAsSt
	p1 := testAsSt{Str: "1"}
	// Bool
	assert.Equal(t, false, AsBool(p))
	assert.Equal(t, true, AsBool(p1))
	assert.Equal(t, true, AsBool(1))
	assert.Equal(t, false, AsBool(0))

	// Integer
	assert.Equal(t, 10, AsInt(10.322))
	assert.Equal(t, 10, AsInt(int64(10)))
	assert.Equal(t, int8(23), AsInt8(23.21))
	assert.Equal(t, int16(12), AsInt16(12))
	assert.Equal(t, int32(15), AsInt32(15.12))
	assert.Equal(t, int64(32), AsInt64(32.1))

	// Uint
	assert.Equal(t, uint(10), AsUint(10.322))
	assert.Equal(t, uint8(10), AsUint8(int64(10)))
	assert.Equal(t, uint8(23), AsUint8(23.21))
	assert.Equal(t, uint16(12), AsUint16(12))
	assert.Equal(t, uint32(15), AsUint32(15.12))
	assert.Equal(t, uint64(32), AsUint64(32.1))

	// Float
	assert.Equal(t, float32(32.0), AsFloat32(32))
	assert.Equal(t, float32(16.0), AsFloat32(16))
	assert.Equal(t, 21.0, AsFloat64(21))
	assert.Equal(t, 12.0, AsFloat64(float32(12.0)))

	// String
	assert.Equal(t, "true", AsString(true))
	assert.Equal(t, "false", AsString(false))
	assert.Equal(t, "10.1", AsString(10.1))
	assert.Equal(t, "---", AsString("---"))
	assert.Equal(t, "0", AsString(0))
	assert.Equal(t, "", AsString(nil))
	assert.Equal(t, "", AsString(p))
}
