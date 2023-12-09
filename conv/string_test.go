package conv

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCamelCase(t *testing.T) {
	var str string
	str = CamelCase("camel_case")
	assert.Equal(t, "camelCase", str)
	str = CamelCase("camel-case")
	assert.Equal(t, "camelCase", str)
	str = CamelCase("camel:case")
	assert.Equal(t, "camelCase", str)
	str = BigCamelCase("camel__case")
	assert.Equal(t, "CamelCase", str)
}

func TestUnderLineCase(t *testing.T) {
	var str string
	str = UnderLineCase("camelCase")
	assert.Equal(t, "camel_case", str)
}
