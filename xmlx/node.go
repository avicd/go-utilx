package xmlx

import (
	"github.com/avicd/go-utilx/logx"
	"strings"
)

type NodeType uint

const (
	ElementNode NodeType = iota
	AttributeNode
	TextNode
	CDataSectionNode
	ProcessingInstructionNode
	CommentNode
	DocumentNode
	DocumentTypeNode
	DirectiveNode
)

type Node struct {
	ParentNode,
	FirstChild,
	LastChild,
	PrevSibling,
	NextSibling *Node
	Type         NodeType
	Value        string
	Name         string
	NamespaceURI string
	Prefix       string
	Attrs        []*Node
	ChildNodes   []*Node
}

func (node *Node) AppendChild(child *Node) {
	if child == nil {
		return
	}
	child.ParentNode = node
	child.PrevSibling = node.LastChild
	node.ChildNodes = append(node.ChildNodes, child)
	if node.LastChild != nil {
		node.LastChild.NextSibling = child
	}
	if node.FirstChild == nil {
		node.FirstChild = child
	}
	node.LastChild = child
}

func (node *Node) AppendAll(nodes []*Node) {
	for _, p := range nodes {
		node.AppendChild(p)
	}
}

func cloneNodeList(list []*Node, parent *Node) []*Node {
	var newList []*Node
	for i, node := range list {
		node = node.CloneNode(true)
		node.ParentNode = parent
		if i > 0 {
			node.PrevSibling = newList[i-1]
			node.PrevSibling.NextSibling = node
		}
		newList = append(newList, node)
	}
	return newList
}

func (node *Node) CloneNode(deep bool) *Node {
	if node == nil {
		return nil
	}
	newNode := &Node{
		FirstChild:   node.FirstChild,
		LastChild:    node.LastChild,
		PrevSibling:  node.PrevSibling,
		NextSibling:  node.NextSibling,
		Type:         node.Type,
		Value:        node.Value,
		Name:         node.Name,
		NamespaceURI: node.NamespaceURI,
		Prefix:       node.Prefix,
		Attrs:        node.Attrs,
		ChildNodes:   node.ChildNodes,
	}
	newNode.Attrs = cloneNodeList(node.Attrs, newNode)
	if deep {
		if len(node.Attrs) > 0 {
			newNode.Attrs = cloneNodeList(node.Attrs, newNode)
		}
		if len(node.ChildNodes) > 0 {
			nodeList := cloneNodeList(node.ChildNodes, newNode)
			newNode.ChildNodes = nodeList
			newNode.FirstChild = nodeList[0]
			newNode.LastChild = nodeList[len(nodeList)-1]
		}
	}
	return newNode
}

func (node *Node) InsertBefore(newNode *Node, refNode *Node) {
	index := node.IndexOf(refNode)
	if index < 0 {
		logx.Error("Invalid ref node when calling *Node.InsertBefore")
		return
	}
	newNode.ParentNode = node
	var nodeList []*Node
	if index > 0 {
		nodeList = node.ChildNodes[:index-1]
	} else {
		node.FirstChild = newNode
	}
	nodeList = append(nodeList, newNode)
	nodeList = append(nodeList, node.ChildNodes[index-1:]...)
	node.ChildNodes = nodeList
	if refNode.PrevSibling != nil {
		refNode.PrevSibling.NextSibling = newNode
	}
	newNode.PrevSibling = refNode.PrevSibling
	newNode.NextSibling = refNode
	refNode.PrevSibling = newNode
}

func (node *Node) InsertAfter(newNode *Node, refNode *Node) {
	index := node.IndexOf(refNode)
	if index < 0 {
		logx.Error("Invalid ref node when calling *Node.InsertAfter")
		return
	}
	newNode.ParentNode = node
	var nodeList []*Node
	if refNode.NextSibling != nil {
		nodeList = append(node.ChildNodes[:index+1], newNode)
		nodeList = append(nodeList, node.ChildNodes[index+1:]...)
	} else {
		nodeList = append(node.ChildNodes, newNode)
	}
	node.ChildNodes = nodeList
	if refNode.NextSibling != nil {
		refNode.NextSibling.PrevSibling = newNode
	}
	newNode.NextSibling = node.NextSibling
	newNode.PrevSibling = refNode
	refNode.NextSibling = newNode
}

func (node *Node) AddSibling(next *Node) {
	if next == nil {
		return
	}
	node.ParentNode.InsertAfter(next, node)
}

func (node *Node) Remove() {
	if node.ParentNode != nil {
		node.ParentNode.RemoveChild(node)
	}
}

func (node *Node) RemoveChild(child *Node) {
	ni := node.IndexOf(child)
	if ni > -1 {
		if child.PrevSibling != nil {
			child.PrevSibling.NextSibling = child.NextSibling
		}
		if child.NextSibling != nil {
			child.NextSibling.PrevSibling = child.PrevSibling
		}
		if child == node.FirstChild {
			node.FirstChild = child.NextSibling
		}
		if child == node.LastChild {
			node.LastChild = child.PrevSibling
		}
		buf := node.ChildNodes[:ni]
		node.ChildNodes = append(buf, node.ChildNodes[ni+1:]...)
	}
}

func (node *Node) RemoveAll(list []*Node) {
	for _, nd := range list {
		node.RemoveChild(nd)
	}
}

func (node *Node) ClearContent() {
	node.ChildNodes = nil
	node.FirstChild = nil
	node.LastChild = nil
}

func (node *Node) IndexOf(ref *Node) int {
	for i, p := range node.ChildNodes {
		if p == ref {
			return i
		}
	}
	return -1
}

func (node *Node) Contains(ref *Node) bool {
	for _, p := range node.ChildNodes {
		if p == ref || p.Contains(ref) {
			return true
		}
	}
	return false
}

func (node *Node) HasChildren() bool {
	return len(node.ChildNodes) > 0
}

func (node *Node) GetRoot() *Node {
	var p *Node
	for p = node; p.ParentNode != nil; p = p.ParentNode {
	}
	return p
}

func (node *Node) Attr(name string) *Node {
	if node.Type == AttributeNode {
		if node.Name == name {
			return node
		}
	} else {
		for _, attr := range node.Attrs {
			if attr.Name == name {
				return attr
			}
		}
	}
	return nil
}

func (node *Node) AttrString(name string) string {
	if attr := node.Attr(name); attr != nil {
		return attr.Value
	}
	return ""
}

func (node *Node) TrimAttrString(name string) string {
	return strings.TrimSpace(node.AttrString(name))
}

func (node *Node) AttrBool(name string) bool {
	return strings.TrimSpace(node.TrimAttrString(name)) == "true"
}

func (node *Node) InnerText() string {
	switch node.Type {
	case ElementNode, DocumentNode:
		buffer := strings.Builder{}
		for _, p := range node.ChildNodes {
			buffer.WriteString(p.InnerText())
		}
		return buffer.String()
	case TextNode, CDataSectionNode, CommentNode:
		return node.Value
	}
	return ""
}

func (node *Node) NameWithPrefix() string {
	name := node.Name
	if node.Prefix != "" {
		name = node.Prefix + ":" + name
	}
	return name
}

func (node *Node) InnerXML() string {
	exporter := &Exporter{IncludeSelf: false}
	return exporter.ApplyOn(node)
}

func (node *Node) InnerHTML() string {
	return node.InnerXML()
}

func (node *Node) Export(exporter *Exporter) string {
	return exporter.ApplyOn(node)
}

func (node *Node) Find(selector string) []*Node {
	xpath, err := NewXpath(selector)
	if err != nil {
		logx.Error(err.Error())
		return nil
	}
	return xpath.SelectAll(node)
}

func (node *Node) FindOne(selector string) *Node {
	xpath, err := NewXpath(selector)
	if err != nil {
		logx.Error(err.Error())
		return nil
	}
	return xpath.SelectFirst(node)
}
