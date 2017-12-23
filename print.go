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

func printXML(buf *bytes.Buffer, n *Node, level int, indent string, withNS bool) {
	pretty := len(indent) > 0

	if pretty {
		buf.WriteString(strings.Repeat(indent, level))
	}
	buf.WriteByte('<')
	if n.NS != nil {
		buf.WriteString(n.NS.Name)
		buf.WriteByte(':')
	}
	buf.WriteString(n.Name)

	if withNS {
		for _, ns := range n.Document.NamespaceList {
			buf.WriteString(" xmlns")
			if ns.Name != "" {
				buf.WriteByte(':')
				buf.WriteString(ns.Name)
			}
			buf.WriteByte('=')
			buf.WriteByte('"')
			xml.Escape(buf, []byte(ns.URI))
			buf.WriteByte('"')
		}
	}

	if len(n.Attributes) > 0 {
		for _, attr := range n.Attributes {
			buf.WriteByte(' ')
			if attr.NS != nil {
				buf.WriteString(attr.NS.Name)
				buf.WriteByte(':')
			}
			buf.WriteString(attr.Name)
			buf.WriteByte('=')
			buf.WriteByte('"')
			xml.Escape(buf, []byte(attr.Value))
			buf.WriteByte('"')
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
			printXML(buf, c, level+1, indent, false)
		}
	}
	if len(n.Text) > 0 {
		xml.EscapeText(buf, []byte(n.Text))
	}

	if len(n.Children) > 0 && len(indent) > 0 {
		buf.WriteString(strings.Repeat(indent, level))
	}
	buf.WriteString("</")
	if n.NS != nil {
		buf.WriteString(n.NS.Name)
		buf.WriteByte(':')
	}
	buf.WriteString(n.Name)
	buf.WriteByte('>')

	if pretty {
		buf.WriteByte('\n')
	}
}
