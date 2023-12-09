package evalx

import "github.com/avicd/go-utilx/datax"

type Target = uint

const (
	PROPERTY Target = iota
	METHOD
)

type Stack struct {
	Ctx    Context
	Target datax.LinkedList[Target]
	Error  error
	X      any
	Y      any
}

func StackOf(accessor Context) *Stack {
	return &Stack{Ctx: accessor}
}

func (it *Stack) PushTarget(target Target) {
	it.Target.Push(target)
}

func (it *Stack) PopTarget() Target {
	if val, ok := it.Target.Pop(); ok {
		return val
	}
	return PROPERTY
}
