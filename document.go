package xmldom

import "bytes"

type Document struct {
	ProcInst   string
	Directives []string
	Root       *Node
}

func (d *Document) XML() string {
	buf := new(bytes.Buffer)
	buf.WriteString(d.ProcInst)
	for _, directive := range d.Directives {
		buf.WriteString(directive)
	}
	buf.WriteString(d.Root.XML())
	return buf.String()
}

func (d *Document) XMLPretty() string {
	buf := new(bytes.Buffer)
	if len(d.ProcInst) > 0 {
		buf.WriteString(d.ProcInst)
		buf.WriteByte('\n')
	}
	for _, directive := range d.Directives {
		buf.WriteString(directive)
		buf.WriteByte('\n')
	}
	buf.WriteString(d.Root.XMLPretty())
	return buf.String()
}
