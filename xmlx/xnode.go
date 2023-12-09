package xmlx

import (
	"strings"
)

type XNodeHandler func(node *Node) bool

type XNode struct {
	*Node
	handleNode XNodeHandler
}

func NewXpathNode(node *Node) *XNode {
	return &XNode{Node: node}
}

func chooseNode(node *Node, sel string) bool {
	if strings.HasSuffix(sel, "()") {
		switch sel {
		case "node()":
			return true
		case "text()":
			return node.Type == TextNode
		}
	}

	if node.Type != ElementNode && node.Type != DocumentNode {
		return false
	}
	if sel == "*" {
		return true
	}
	return sel == node.Name
}

func (it *XNode) setCheckin(checker XNodeHandler) {
	it.handleNode = checker
}

func (it *XNode) Ancestor(sel string) {
	for p := it.ParentNode; p != nil; p = p.ParentNode {
		if !chooseNode(p, sel) {
			continue
		}
		if !it.handleNode(p) {
			break
		}
	}
}

func (it *XNode) AncestorOrSelf(sel string) {
	if !it.chooseSelf(sel) {
		return
	}
	it.Ancestor(sel)
}

func (it *XNode) Attribute(sel string) {
	for _, p := range it.Attrs {
		if sel == p.Name && !it.handleNode(p) {
			break
		}
	}
}

func (it *XNode) Child(sel string) {
	for _, p := range it.ChildNodes {
		if chooseNode(p, sel) && !it.handleNode(p) {
			break
		}
	}
}

func nodeTreeLoop(node *Node, sel string, handleNode XNodeHandler) {
	for _, p := range node.ChildNodes {
		if chooseNode(p, sel) && !handleNode(p) {
			break
		}
		if len(p.ChildNodes) > 0 {
			nodeTreeLoop(p, sel, handleNode)
		}
	}
}

func (it *XNode) Descendant(sel string) {
	nodeTreeLoop(it.Node, sel, it.handleNode)
}

func (it *XNode) DescendantOrSelf(sel string) {
	if !it.chooseSelf(sel) {
		return
	}
	it.Descendant(sel)
}

func (it *XNode) Following(sel string) {
	for p := it.NextSibling; p != nil; p = p.NextSibling {
		xnd := NewXpathNode(p)
		next := true
		xnd.setCheckin(func(node *Node) bool {
			next = it.handleNode(node)
			return next
		})
		xnd.DescendantOrSelf(sel)
		if !next {
			return
		}
	}
	if it.ParentNode != nil {
		xnd := NewXpathNode(it.ParentNode)
		xnd.setCheckin(func(node *Node) bool {
			return it.handleNode(node)
		})
		xnd.Following(sel)
	}
}

func (it *XNode) Namespace(sel string) {
	if it.NamespaceURI == sel {
		it.handleNode(it.Node)
	}
}

func (it *XNode) Parent() {
	if chooseNode(it.Node, "*") {
		it.handleNode(it.ParentNode)
	}
}

func (it *XNode) Preceding(sel string) {
	if !it.chooseSelf(sel) {
		return
	}
	it.PrecedingSibling(sel)
}

func (it *XNode) PrecedingSibling(sel string) {
	for p := it.PrevSibling; p != nil; p = p.PrevSibling {
		if !it.handleNode(p) {
			break
		}
	}
}

func (it *XNode) Self(sel string) {
	it.chooseSelf(sel)
}

func (it *XNode) Root() {
	root := it.GetRoot()
	if chooseNode(root, "*") {
		it.handleNode(root)
	}
}

func (it *XNode) chooseSelf(sel string) bool {
	if chooseNode(it.Node, sel) && !it.handleNode(it.Node) {
		return false
	}
	return true
}
