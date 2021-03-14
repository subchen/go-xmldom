package xmldom

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
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

	// наследование xmlns от вышестоящего элемента если он есть
	if level == 0 && n.Name.Space != "" && space != nil {
		// для тех случаев когда нужно распечатать часть большого XML, например Signature/SignedInfo
		buf.WriteByte(' ')
		buf.WriteString(space.Name.Space)
		buf.WriteByte(':')
		buf.WriteString(space.Name.Local)
		buf.WriteByte('=')
		buf.WriteByte('"')
		xml.Escape(buf, []byte(space.Value))
		buf.WriteByte('"')
	}
	if len(n.Attributes) > 0 {
		for _, attr := range n.Attributes {
			if attr.Name.Space == "" {
				buf.WriteByte(' ')
				buf.WriteString(attr.Name.Local)
				buf.WriteByte('=')
				buf.WriteByte('"')
				xmlEscape( buf, []byte(attr.Value) )
				buf.WriteByte('"')
			} else if attr.Name.Space == "xmlns" {
				// выставляем xmlns для тех тегов где он был изначально прописан,
				// игнорируем первый тег так как для него уже есть условие
				if level > 0 && space != nil && attr.Name.Local == space.Name.Local {
					buf.WriteByte(' ')
					buf.WriteString(attr.Name.Space)
					buf.WriteByte(':')
					buf.WriteString(attr.Name.Local)
					buf.WriteByte('=')
					buf.WriteByte('"')
					xmlEscape( buf, []byte(attr.Value) )
					buf.WriteByte('"')
				}
			}
		}
	}
	if n.Document.EmptyElementTag {
		if len(n.Children) == 0 && len(n.Text) == 0 {
			buf.WriteString(" />")
			if pretty {
				buf.WriteByte('\n')
			}
			return
		}
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
		if n.Document.TextSafeMode {
			xml.EscapeText(buf, []byte(n.Text))
		} else {
			buf.WriteString(n.Text)
		}
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



func xmlEscape( w io.Writer, value []byte ) {
	var res bytes.Buffer
	xml.Escape(&res, value)

	out := res.String()
	newOut := strings.Replace(out, "&#34;" ,"&quot;" , -1  )

	fmt.Fprintf( w, "%s", newOut  )
}

