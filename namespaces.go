package xmldom

import (
	"encoding/xml"
)

type Namespaces []xml.Attr

func (n *Namespaces) GetName(space string) *xml.Attr {
	for _, attr := range *n {
		if attr.Value == space {
			return &attr
		}
	}
	return nil
}