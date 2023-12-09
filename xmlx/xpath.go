package xmlx

import (
	"fmt"
	"github.com/avicd/go-utilx/conv"
	"github.com/avicd/go-utilx/evalx"
	"github.com/avicd/go-utilx/refx"
	"io"
	"strings"
)

type Xpath struct {
	stacks []*xStack
	expr   string
}

func NewXpath(text string) (*Xpath, error) {
	stacks, err := parseXpath(text)
	if err != nil {
		return nil, err
	}
	return &Xpath{stacks: stacks, expr: text}, nil
}

func (it *Xpath) SelectFirst(node *Node) *Node {
	for _, stack := range it.stacks {
		list := stack.eval(node, SelectFirst)
		if len(list) > 0 {
			return list[0]
		}
	}
	return nil
}

func (it *Xpath) SelectAll(node *Node) []*Node {
	var list []*Node
	for _, stack := range it.stacks {
		rslt := stack.eval(node, SelectAll)
		if len(rslt) > 0 {
			list = append(list, rslt...)
		}
	}
	return list
}

func (it *Xpath) SelectLast(node *Node) *Node {
	list := it.SelectAll(node)
	if len(list) > 0 {
		return list[len(list)-1]
	}
	return nil
}

type xStack struct {
	call   string
	sel    string
	axis   string
	choose string
	root   bool
	next   *xStack
}

func (it *xStack) pushNext() *xStack {
	next := &xStack{}
	if it != nil {
		it.next = next
	}
	return next
}

func (it *xStack) evalStack(input *Node, handler XNodeHandler, pool *indexPool) {
	xnd := NewXpathNode(input)
	var list []*Node
	xnd.setCheckin(func(node *Node) bool {
		pool.cacheIndex(node)
		list = append(list, node)
		if it.choose == "" {
			if it.next != nil {
				it.next.evalStack(node, handler, pool)
			} else {
				return handler(node)
			}
		}
		return true
	})
	ctx := newEvalCtx(xnd)
	ctx.bind("indexPool", pool)
	evalx.Eval(it.axis, ctx)
	if it.choose != "" && len(list) > 0 {
		for _, p := range list {
			pctx := newEvalCtx(NewXpathNode(p))
			val, _ := evalx.Eval(it.choose, pctx)
			var pass *Node
			if refx.IsGeneralInt(val) {
				if refx.AsInt64(val) == int64(pool.GetIndex(p)) {
					pass = p
				}
			} else if refx.AsBool(val) {
				pass = p
			}
			if pass != nil {
				if it.next != nil {
					it.next.evalStack(pass, handler, pool)
				} else {
					handler(pass)
				}
			}
		}
	}
}

func (it *xStack) eval(input *Node, stype SelectType) []*Node {
	var list []*Node
	pool := &indexPool{cache: map[*Node]map[*Node]int{}}
	var handler XNodeHandler
	rc := map[*Node]bool{}
	handler = func(node *Node) bool {
		if _, ok := rc[node]; !ok {
			rc[node] = true
			list = append(list, node)
			if stype == SelectFirst {
				return false
			}
		}
		return true
	}
	it.evalStack(input, handler, pool)
	return list
}

func parseXpath(text string) ([]*xStack, error) {
	r := strings.NewReader(strings.TrimSpace(text))
	var pre byte
	var stacks []*xStack
	var root *xStack
	var current *xStack
	buf := &strings.Builder{}
	index := -1
	var sl byte
	var cl bool
	split := func() {
		if current != nil && buf.Len() > 0 {
			current.sel = buf.String()
			buf.Reset()
		}
	}
	collect := func() {
		split()
		stacks = append(stacks, root)
		root = nil
		current = nil
		index = -1
	}
	for {
		ch, err := r.ReadByte()
		if err == io.EOF {
			collect()
			break
		}
		index++
		if (sl == 0) || ch == sl {
			switch ch {
			case '/':
				if pre == '/' {
					current.call = "DescendantOrSelf"
				} else {
					split()
					current = current.pushNext()
					current.call = "Child"
					current.root = index == 0
				}
			case '[':
				current.sel = buf.String()
				buf.Reset()
				cl = true
			case ']':
				current.choose = buf.String()
				buf.Reset()
				cl = false
			case '\'', '"':
				if sl == 0 {
					sl = ch
				} else {
					sl = 0
				}
				buf.WriteByte('"')
			case ':':
				if pre == ':' {
					current.call = conv.BigCamelCase(buf.String())
					buf.Reset()
				}
			case '.':
				if pre == '.' {
					current.call = "Parent"
				}
			case '@':
				if !cl {
					if current == nil {
						current = current.pushNext()
					}
					current.call = "Attribute"
				} else {
					buf.WriteString(xpathAttr)
				}
			case '=':
				if !(pre == '!' || pre == '>' || pre == '<') {
					buf.WriteString("==")
				} else {
					buf.WriteByte('=')
				}
			case '|':
				collect()
			default:
				if current == nil {
					current = current.pushNext()
					current.call = "Child"
				}
				buf.WriteByte(ch)
			}
		} else {
			buf.WriteByte(ch)
		}

		if current != nil && root == nil {
			root = current
		}
		pre = ch
	}
	for hi, hd := range stacks {
		if hd.root {
			stacks[hi] = &xStack{next: hd, call: "self.Root"}
		} else {
			hd.call = "self." + hd.call
		}
		for p := stacks[hi]; p != nil; p = p.next {
			var axis string
			if p.sel != "" {
				axis = fmt.Sprintf("%s(\"%s\")", p.call, p.sel)
			} else {
				axis = p.call + "()"
			}
			p.axis = axis
		}
	}
	return stacks, nil
}
