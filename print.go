package xmldom

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
)

func stringifyProcInst(pi *xml.ProcInst) string {
	if pi == nil {
		return ""
	}
	return fmt.Sprintf("<?%s %s?>", pi.Target, string(pi.Inst))
}

func stringifyDirective(directive *xml.Directive) string {
	if directive == nil {
		return ""
	}
	return fmt.Sprintf("<!%s>", string(*directive))
}

type printer struct{}

func (p *printer) printXML(buf *bytes.Buffer, n *Node, level int, indent string) {
	pretty := len(indent) > 0

	if pretty {
		buf.WriteString(strings.Repeat(indent, level))
	}

	space := n.GetNamespace()

	buf.WriteByte('<')
	if space != nil {
		buf.WriteString(space.Name.Local)
		buf.WriteByte(':')
	}
	buf.WriteString(n.Name.Local)

	if len(n.Attributes) > 0 {
		for _, attr := range n.Attributes {

			if attr.Name.Space == "" {
				buf.WriteByte(' ')
				buf.WriteString(attr.Name.Local)
				buf.WriteByte('=')
				buf.WriteByte('"')
				xml.Escape(buf, []byte(attr.Value))
				buf.WriteByte('"')
			} else if attr.Name.Space == "xmlns" {
				if space != nil && attr.Name.Local == space.Name.Local {
					buf.WriteByte(' ')
					buf.WriteString(attr.Name.Space)
					buf.WriteByte(':')
					buf.WriteString(attr.Name.Local)
					buf.WriteByte('=')
					buf.WriteByte('"')
					xml.Escape(buf, []byte(attr.Value))
					buf.WriteByte('"')
				}

			}
		}
	}

	if len(n.Children) == 0 && len(n.Text) == 0 {
		buf.WriteString(" />")
		if pretty {
			buf.WriteByte('\n')
		}
		return
	}

	buf.WriteByte('>')

	if len(n.Children) > 0 {
		if pretty {
			buf.WriteByte('\n')
		}
		for _, c := range n.Children {
			p.printXML(buf, c, level+1, indent)
		}
	}
	if len(n.Text) > 0 {
		xml.EscapeText(buf, []byte(n.Text))
	}

	if len(n.Children) > 0 && len(indent) > 0 {
		buf.WriteString(strings.Repeat(indent, level))
	}
	buf.WriteString("</")
	if space != nil {
		buf.WriteString(space.Name.Local)
		buf.WriteByte(':')
	}
	buf.WriteString(n.Name.Local)
	buf.WriteByte('>')

	if pretty {
		buf.WriteByte('\n')
	}
}
