package evalx

import (
	"github.com/avicd/go-utilx/bufx"
	"github.com/avicd/go-utilx/refx"
	"go/ast"
)

type Context interface {
	ValueOf(ident string) (any, bool)
	MethodOf(ident string) (any, bool)
	CacheOf(text string) (ast.Expr, bool)
	Cache(text string, expr ast.Expr)
}

type Scope struct {
	binds   map[string]any
	backup  map[string]any
	callers map[string]any
	vars    []any
}

var top *Scope
var topCache bufx.Cache[string, ast.Expr]

func init() {
	top = NewScope()
	topCache = &bufx.LruCache[string, ast.Expr]{Size: 1000}
}

func SetCacheSize(size int) {
	topCache = &bufx.LruCache[string, ast.Expr]{Size: size}
}

func TopScope() *Scope {
	if top == nil {
		top = NewScope()
	}
	return top
}

func NewScope(vars ...any) *Scope {
	return &Scope{
		binds:   map[string]any{},
		backup:  map[string]any{},
		callers: map[string]any{},
		vars:    vars,
	}
}

func (it *Scope) ValueOf(ident string) (any, bool) {
	if val, ok := refx.PropOfId(it.binds, ident); ok {
		return val, true
	}
	for _, obj := range it.vars {
		if val, ok := refx.PropOfId(obj, ident); ok {
			return val, true
		}
	}
	if !it.IsTop() {
		return TopScope().ValueOf(ident)
	}
	return nil, false
}

func (it *Scope) MethodOf(ident string) (any, bool) {
	for _, obj := range it.vars {
		if method, ok := refx.MethodOfId(obj, ident); ok {
			return method, true
		}
	}
	if method, ok := it.callers[ident]; ok {
		return method, true
	} else if !it.IsTop() {
		return TopScope().MethodOf(ident)
	}
	return nil, false
}

func (it *Scope) CacheOf(text string) (ast.Expr, bool) {
	return topCache.Get(text)
}

func (it *Scope) Cache(text string, expr ast.Expr) {
	topCache.Put(text, expr)
}

func (it *Scope) IsTop() bool {
	return it == top
}

func (it *Scope) Bind(name string, value any) {
	if refx.IsFunc(value) {
		it.callers[name] = value
	} else {
		it.binds[name] = value
	}
}

func (it *Scope) UnBind(name string) {
	delete(it.binds, name)
}

func (it *Scope) Backup(name string) {
	if val, ok := it.binds[name]; ok {
		it.backup[name] = val
	}
}

func (it *Scope) Restore(name string) {
	if val, ok := it.backup[name]; ok {
		it.binds[name] = val
		delete(it.backup, name)
	}
}

func (it *Scope) Merge(val any) *Scope {
	var binds any
	var vars []any
	val = refx.ValueOf(val).Interface()
	switch tmp := val.(type) {
	case *Scope:
		if tmp != nil {
			binds = tmp.binds
			vars = tmp.vars
		}
	case Scope:
		binds = tmp.binds
		vars = tmp.vars
	default:
		binds = val
	}
	if binds != nil {
		if refx.IndirectType(binds) == refx.TypeOf(it.binds) {
			refx.Merge(&it.binds, binds)
		} else {
			it.Link(val)
		}
	}
	for _, item := range vars {
		it.Link(item)
	}
	return it
}

func (it *Scope) Link(items ...any) *Scope {
	if len(items) < 1 {
		return it
	}
	var dest []any
	dest = append(dest, items...)
	dest = append(dest, it.vars...)
	it.vars = dest
	return it
}

func (it *Scope) UnLink(obj any) *Scope {
	index := -1
	for i, p := range it.vars {
		if p == obj {
			index = i
			break
		}
	}
	if index > -1 {
		vars := it.vars[:index]
		it.vars = append(vars, it.vars[index+1:])
	}
	return it
}

func (it *Scope) Eval(text string) (any, error) {
	return Eval(text, it)
}
