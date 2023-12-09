package tokx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDoubleQuota(t *testing.T) {
	text := "'testA' + \"testB\""
	assert.Equal(t, "\"testA\" + \"testB\"", DoubleQuota(text))
}
