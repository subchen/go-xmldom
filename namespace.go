package xmldom

type Namespace struct {
	Name string
	URI  string
}

func (ns *Namespace) IsDefault() bool {
	return len(ns.Name) == 0
}
