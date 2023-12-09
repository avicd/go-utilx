package xmlx

import (
	"bufio"
	"encoding/xml"
	"github.com/avicd/go-utilx/logx"
	"golang.org/x/net/html/charset"
	"io"
	"strings"
	"unicode"
)

const cdataOpen = "<![CDATA["
const xmlnsPrefix = "xmlns"
const xpathAttr = "__Attr__"

type xmlStack struct {
	*Node
	prefix map[string]string
	ns     map[string]string
	next   *xmlStack
	prev   *xmlStack
}

func (stack *xmlStack) pushNext(node *Node, attrs []xml.Attr) *xmlStack {
	next := &xmlStack{
		Node:   node,
		prefix: map[string]string{},
		ns:     map[string]string{},
	}
	node.Attrs = []*Node{}
	for i, attr := range attrs {
		attrNs := attr.Name.Space
		// namespace with prefix
		if attr.Name.Space == xmlnsPrefix ||
			// default namespace
			attr.Name.Space == "" && attr.Name.Local == xmlnsPrefix {
			attrNs = ""
			if attr.Name.Space != "" {
				next.prefix[attr.Value] = attr.Name.Local
			}
			next.ns[attr.Name.Local] = attr.Value
		}
		var attrPrefix string
		if stack != nil {
			attrPrefix = stack.getPrefix(attrNs)
		}
		newAttr := &Node{
			Type:         AttributeNode,
			ParentNode:   node,
			Name:         attr.Name.Local,
			Value:        attr.Value,
			NamespaceURI: attrNs,
			Prefix:       attrPrefix,
		}
		if i > 0 {
			newAttr.PrevSibling = node.Attrs[i-1]
			newAttr.PrevSibling.NextSibling = newAttr
		}
		node.Attrs = append(node.Attrs, newAttr)
	}
	next.prev = stack
	if stack != nil {
		stack.next = next
	}
	return next
}

func (stack *xmlStack) popLast() *xmlStack {
	if stack != nil {
		prev := stack.prev
		prev.next = nil
		return prev
	}
	return nil
}

func (stack *xmlStack) getPrefix(ns string) string {
	for p := stack; p != nil; p = p.prev {
		if prefix, ok := p.prefix[ns]; ok {
			return prefix
		}
	}
	return ""
}

func (stack *xmlStack) getNs(prefix string) string {
	for p := stack; p != nil; p = p.prev {
		if ns, ok := p.ns[prefix]; ok {
			return ns
		}
	}
	return ""
}

func (stack *xmlStack) parseStrAttr(text string) []*Node {
	buf := strings.TrimSpace(text)
	var attrs []*Node
	for len(buf) > 0 {
		cut := strings.IndexFunc(buf, func(r rune) bool {
			return !unicode.IsSpace(r)
		})
		if cut > -1 {
			buf = buf[cut:]
		}
		var key string
		var value string
		cut = strings.Index(buf, "=")
		if cut > -1 {
			key = buf[0:cut]
			buf = buf[cut+1:]
			ptoken := buf[0:1]
			buf = buf[1:]
			cut = strings.Index(buf, ptoken)
			value = buf[0:cut]
			buf = buf[cut+1:]
		} else {
			cut = strings.IndexFunc(buf, unicode.IsSpace)
			if cut > -1 {
				key = buf[0:cut]
				buf = buf[cut+1:]
			} else {
				key = buf
				buf = ""
			}
		}
		if key != "" {
			attr := &Node{Type: AttributeNode, Value: value}
			if cut = strings.Index(key, ":"); cut > -1 {
				attr.Name = key[cut+1 : 0]
				attr.Prefix = key[:cut]
				attr.NamespaceURI = stack.getNs(attr.Prefix)
			} else {
				attr.Name = key
			}
			attrs = append(attrs, attr)
		}
	}
	return attrs
}

type xmlParser struct {
	*bufio.Reader
	ci      int
	isCData bool
}

func newXmlParser(reader io.Reader) *xmlParser {
	return &xmlParser{Reader: bufio.NewReader(reader)}
}

func Parse(reader io.Reader) (*Node, error) {
	return newXmlParser(reader).parse()
}

func (parser *xmlParser) ReadByte() (byte, error) {
	bt, err := parser.Reader.ReadByte()
	if err == nil {
		if !parser.isCData && bt == cdataOpen[parser.ci] {
			parser.ci++
			if parser.ci == len(cdataOpen) {
				parser.isCData = true
			}
		} else {
			parser.ci = 0
		}
	}
	return bt, err
}

func (parser *xmlParser) decoder() *xml.Decoder {
	decoder := xml.NewDecoder(parser)
	decoder.CharsetReader = charset.NewReaderLabel
	return decoder
}

func (parser *xmlParser) parse() (*Node, error) {
	root := &Node{
		Type: DocumentNode,
		Name: "document",
	}
	var current *xmlStack
	current = current.pushNext(root, nil)
	decoder := parser.decoder()
	for {
		parser.isCData = false
		xtk, err := decoder.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			logx.Error(err.Error())
			return nil, err
		}
		var node *Node
		parent := current
		switch el := xtk.(type) {
		case xml.StartElement:
			node = &Node{
				Type:         ElementNode,
				Name:         el.Name.Local,
				NamespaceURI: el.Name.Space,
				Prefix:       parent.getPrefix(el.Name.Space),
			}
			current = current.pushNext(node, el.Attr)
		case xml.EndElement:
			current = current.popLast()
		case xml.CharData:
			text := string(el)
			node = &Node{
				Type:  TextNode,
				Name:  "text",
				Value: text,
			}
			if parser.isCData {
				node.Type = CDataSectionNode
				node.Name = ""
			}
		case xml.Comment:
			text := string(el)
			node = &Node{Type: CommentNode, Name: "comment", Value: text}
		case xml.ProcInst:
			node = &Node{Type: ProcessingInstructionNode, Name: el.Target, Value: string(el.Inst)}
			node.Attrs = parent.parseStrAttr(string(el.Inst))
		case xml.Directive:
			text := string(el)
			node = &Node{Type: DirectiveNode}
			cut := strings.IndexFunc(text, unicode.IsSpace)
			if cut > -1 {
				node.Name = text[0:cut]
				node.Value = strings.TrimSpace(text[cut:])
			} else {
				node.Name = text
				node.Value = ""
			}
			if strings.ToUpper(node.Name) == "DOCTYPE" {
				node.Type = DocumentTypeNode
				cut = strings.IndexFunc(node.Value, unicode.IsSpace)
				if cut > -1 {
					node.Name = node.Value[0:cut]
					node.Value = node.Value[cut:]
				} else {
					node.Name = node.Value
					node.Value = ""
				}
			}
		}
		parent.Node.AppendChild(node)
	}
	return root, nil
}
