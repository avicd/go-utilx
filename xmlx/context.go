package xmlx

import (
	"github.com/avicd/go-utilx/refx"
	"go/ast"
	"strings"
)

type SelectType uint

const (
	SelectAll SelectType = iota
	SelectFirst
)

type xpathContext struct {
	binds map[string]any
	xnd   *XNode
}

func newEvalCtx(xnd *XNode) *xpathContext {
	ctx := &xpathContext{binds: map[string]any{}, xnd: xnd}
	ctx.bind("self", xnd)
	return ctx
}

func (ctx *xpathContext) bind(name string, value any) {
	ctx.binds[name] = value
}

func (ctx *xpathContext) ValueOf(ident string) (any, bool) {
	if val, ok := ctx.binds[ident]; ok {
		return val, ok
	}
	if strings.HasPrefix(ident, xpathAttr) {
		ident = strings.TrimPrefix(ident, xpathAttr)
		return ctx.xnd.AttrString(ident), true
	} else {
		var target *Node
		xnd := NewXpathNode(ctx.xnd.Node)
		xnd.setCheckin(func(node *Node) bool {
			target = node
			return false
		})
		xnd.Child(ident)
		if target != nil {
			return target.InnerText(), true
		}
	}
	return nil, false
}

func (ctx *xpathContext) MethodOf(ident string) (any, bool) {
	return refx.MethodOfId(ctx.xnd, ident)
}

func (ctx *xpathContext) CacheOf(text string) (ast.Expr, bool) {
	return nil, false
}

func (ctx *xpathContext) Cache(text string, expr ast.Expr) {

}

// parent->child->index
type indexPool struct {
	cache map[*Node]map[*Node]int
}

func (c *indexPool) cacheIndex(node *Node) {
	pool := c.cache[node.ParentNode]
	if pool == nil {
		pool = map[*Node]int{}
		c.cache[node.ParentNode] = pool
	}
	if _, ok := pool[node]; !ok {
		pool[node] = len(pool) + 1
	}
}

func (c *indexPool) GetIndex(node *Node) int {
	pool := c.cache[node.ParentNode]
	if pool == nil {
		pool = map[*Node]int{}
		c.cache[node.ParentNode] = pool
	}
	if val, ok := pool[node]; ok {
		return val
	}
	return -1
}
