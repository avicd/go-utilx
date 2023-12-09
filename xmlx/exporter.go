package xmlx

import (
	"fmt"
	"strings"
)

type Exporter struct {
	DropComment               bool
	DropEmpty                 bool
	DropCDataSection          bool
	DropProcessingInstruction bool
	DropDocumentType          bool
	DropDirective             bool
	IncludeSelf               bool
	node                      *Node
}

func (et *Exporter) ApplyOn(node *Node) string {
	var exporter *Exporter
	if et == nil {
		exporter = &Exporter{node: node, IncludeSelf: true}
	} else {
		exporter = &Exporter{
			node:                      node,
			DropComment:               et.DropComment,
			DropEmpty:                 et.DropEmpty,
			DropCDataSection:          et.DropCDataSection,
			DropProcessingInstruction: et.DropProcessingInstruction,
			DropDocumentType:          et.DropDocumentType,
			DropDirective:             et.DropDirective,
			IncludeSelf:               et.IncludeSelf,
		}
	}
	return exporter.export(node)
}

func (et *Exporter) stringifyAttrs(attrs []*Node) string {
	buffer := &strings.Builder{}
	for _, attr := range attrs {
		buffer.WriteString(fmt.Sprintf(" %s=\"%s\"", attr.NameWithPrefix(), attr.Value))
	}
	return buffer.String()
}

func (et *Exporter) export(node *Node) string {
	buffer := &strings.Builder{}
	switch node.Type {
	case ElementNode, DocumentNode:
		withSelf := node.Type != DocumentNode && (et.IncludeSelf || et.node != node)
		var selfName string
		if withSelf {
			selfName = node.NameWithPrefix()
			attrsStr := et.stringifyAttrs(node.Attrs)
			buffer.WriteString(fmt.Sprintf("<%s%s>", selfName, attrsStr))
		}
		for _, p := range node.ChildNodes {
			buffer.WriteString(et.export(p))
		}
		if withSelf {
			buffer.WriteString(fmt.Sprintf("</%s>", selfName))
		}
	case TextNode:
		text := node.Value
		if et.DropEmpty {
			text = strings.TrimSpace(node.Value)
		}
		buffer.WriteString(text)
	case CDataSectionNode:
		if !et.DropCDataSection {
			buffer.WriteString(fmt.Sprintf("<![CDATA[%s]]>", node.Value))
		}
	case ProcessingInstructionNode:
		if !et.DropProcessingInstruction {
			buffer.WriteString(fmt.Sprintf("<?%s%s?>", node.Name, et.stringifyAttrs(node.Attrs)))
		}
	case CommentNode:
		if !et.DropComment {
			buffer.WriteString(fmt.Sprintf("<!--%s-->", et.node.Value))
		}
	case DocumentTypeNode:
		if !et.DropDocumentType {
			buffer.WriteString(fmt.Sprintf("<!DOCTYPE %s%s>", node.Name, node.Value))
		}
	case DirectiveNode:
		if !et.DropDirective {
			buffer.WriteString(fmt.Sprintf("<!%s %s>", node.Name, node.Value))
		}
	}
	return buffer.String()
}
