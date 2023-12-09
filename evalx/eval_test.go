package evalx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testEval struct {
	A      int
	B      string
	GetAge func() int
}

func (p *testEval) GetName(add string) string {
	return "Allen" + add
}

func TestEvalCall(t *testing.T) {
	// Map
	v0 := map[string]any{
		"user": map[string]any{
			"getName": func() string {
				return "Allen"
			},
		},
	}
	v1 := map[string]any{
		"local": map[string]any{
			"users": []any{
				v0,
			},
		},
	}
	r1, err1 := NewScope(v0).Eval("user.getName()")
	assert.Equal(t, "Allen", r1)
	assert.Equal(t, nil, err1)

	r2, err2 := NewScope(v1).Eval("local.users[0]['user']['getName']()")
	assert.Equal(t, "Allen", r2)
	assert.Equal(t, nil, err2)

	r3, err3 := NewScope(v1).Eval("local.users[0].user['getName']()")
	assert.Equal(t, "Allen", r3)
	assert.Equal(t, nil, err3)

	// Struct
	p0 := &testEval{A: 25, B: "30", GetAge: func() int {
		return 21
	}}
	p0r1, p0err1 := NewScope(p0).Eval("GetAge()")
	assert.Equal(t, 21, p0r1)
	assert.Equal(t, nil, p0err1)

	p0r2, p0err2 := NewScope(p0).Eval("GetName(' is 21')")
	assert.Equal(t, "Allen is 21", p0r2)
	assert.Equal(t, nil, p0err2)

	p0r3, p0err3 := NewScope(p0).Eval("GetName(' is ' + A)")
	assert.Equal(t, "Allen is 25", p0r3)
	assert.Equal(t, nil, p0err3)

	p0r4, p0err4 := NewScope(p0).Eval("GetName0(' is ' + A)")
	assert.Equal(t, nil, p0r4)
	assert.NotEqual(t, nil, p0err4)

	// Bind
	getName := func(add string) string {
		return "Allen" + add
	}
	scope := NewScope()
	scope.Bind("getName", getName)
	b0r1, b0err1 := scope.Eval("getName(' is 21')")
	assert.Equal(t, "Allen is 21", b0r1)
	assert.Equal(t, nil, b0err1)

	bor2, b0err2 := scope.Eval("getName0(' is 21')")
	assert.Equal(t, nil, bor2)
	assert.NotEqual(t, nil, b0err2)
}

func TestScope_Cache(t *testing.T) {
	p0 := &testEval{A: 25, B: "30", GetAge: func() int {
		return 21
	}}
	NewScope(p0).Eval("GetAge()")
	_, ok := topCache.Get("GetAge()")
	assert.Equal(t, true, ok)
}
