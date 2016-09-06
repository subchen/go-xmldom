package xmldom

type Document struct {
	ProcInst string
	Root     *Node
}

func (d *Document) XML() string {
	return d.ProcInst + d.Root.XML()
}

func (d *Document) XMLPretty() string {
	if len(d.ProcInst) > 0 {
		return d.ProcInst + "\n" + d.Root.XMLPretty()
	}
	return d.Root.XMLPretty()
}
