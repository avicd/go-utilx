package tokx

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestHasPair(t *testing.T) {
	pair := NewPair("#{", "}")
	assert.Equal(t, true, pair.Match("#{}"))
	assert.Equal(t, false, pair.Match("\\#{}"))
	assert.Equal(t, false, pair.Match("#{\\}"))
}

func TestPairToken_Map(t *testing.T) {
	pair := NewPair("#{", "}")
	sql := "INSERT INTO USER(name,age,height) VALUES(#{name},#{age},#{height})"
	var list []string
	val := pair.Map(sql, func(s string) string {
		list = append(list, s)
		return "?"
	})
	assert.Equal(t, val, "INSERT INTO USER(name,age,height) VALUES(?,?,?)")
	assert.Equal(t, true, reflect.DeepEqual(list, []string{"name", "age", "height"}))
}
