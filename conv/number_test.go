package conv

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseNum(t *testing.T) {
	var info *numInfo
	info = parseNum("fsdfjj0x32332.3232weqw")
	assert.Equal(t, "32332", info.integer)
	assert.Equal(t, "3232", info.decimal)
	assert.Equal(t, 16, info.base)
	assert.Equal(t, false, info.negative)

	info = parseNum("-1.223")
	assert.Equal(t, "1", info.integer)
	assert.Equal(t, "223", info.decimal)
	assert.Equal(t, 10, info.base)
	assert.Equal(t, true, info.negative)

	info = parseNum("2.0323")
	assert.Equal(t, "2", info.integer)
	assert.Equal(t, "0323", info.decimal)
	assert.Equal(t, 10, info.base)
	assert.Equal(t, false, info.negative)

	info = parseNum("0x2.0323fgdf")
	assert.Equal(t, "2", info.integer)
	assert.Equal(t, "0323f", info.decimal)
	assert.Equal(t, 16, info.base)
	assert.Equal(t, false, info.negative)

	info = parseNum("0o78721")
	assert.Equal(t, "7", info.integer)
	assert.Equal(t, "", info.decimal)
	assert.Equal(t, 8, info.base)
	assert.Equal(t, false, info.negative)

	info = parseNum("-0b10010.11001")
	assert.Equal(t, "10010", info.integer)
	assert.Equal(t, "11001", info.decimal)
	assert.Equal(t, 2, info.base)
	assert.Equal(t, true, info.negative)
}

func TestParseInt(t *testing.T) {
	assert.Equal(t, int64(255), ParseInt("0xff"))
	assert.Equal(t, int64(57), ParseInt("0o71"))
	assert.Equal(t, int64(2), ParseInt("0b10"))
	assert.Equal(t, int64(4454), ParseInt("4454"))
	assert.Equal(t, int64(18), ParseInt("0x12.212"))
}

func TestParseFloat(t *testing.T) {
	assert.Equal(t, 1.9375, ParseFloat("0x1.f"))
	assert.Equal(t, 1.875, ParseFloat("0o1.7"))
	assert.Equal(t, 1.8125, ParseFloat("0b1.1101"))
	assert.Equal(t, 2.1121, ParseFloat("2.1121"))
	assert.Equal(t, 2.0, ParseFloat("2"))
}
