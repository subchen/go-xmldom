// XML DOM processing for Golang, supports xpath query
package xmldom

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	xmlPrefix   = "xml"
	xmlUrl      = "http://www.w3.org/XML/1998/namespace"
	xmlnsPrefix = "xmlns"
	xmlnsUrl    = "http://www.w3.org/2000/xmlns"
	xlinkPrefix = "xlink"
	xlinkUrl    = "http://www.w3.org/1999/xlink"
	xsiPrefix   = "xsi"
	xsiUrl      = "http://www.w3.org/2001/XMLSchema-instance"
)

func Must(doc *Document, err error) *Document {
	if err != nil {
		panic(err)
	}
	return doc
}

func ParseXML(s string) (*Document, error) {
	return Parse(strings.NewReader(s))
}

func ParseFile(filename string) (*Document, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return Parse(file)
}

func Parse(r io.Reader) (*Document, error) {
	p := xml.NewDecoder(r)
	t, err := p.Token()
	if err != nil {
		return nil, err
	}

	doc := new(Document)
	var e *Node
	for t != nil {
		switch token := t.(type) {
		case xml.StartElement:
			// a new node
			el := new(Node)
			el.Document = doc
			el.Parent = e
			el.Name = token.Name.Local
			for _, attr := range token.Attr {
				var name, ns string
				if attr.Name.Space != "" {
					ns = attr.Name.Space
					switch ns {
					case xmlnsUrl:
						name = fmt.Sprintf("%s:%s", xmlnsPrefix, attr.Name.Local)
					case xmlUrl:
						name = fmt.Sprintf("%s:%s", xmlPrefix, attr.Name.Local)
					case xlinkUrl:
						name = fmt.Sprintf("%s:%s", xlinkPrefix, attr.Name.Local)
					case xsiUrl:
						name = fmt.Sprintf("%s:%s", xsiPrefix, attr.Name.Local)
					default:
						name = fmt.Sprintf("%s:%s", attr.Name.Space, attr.Name.Local)
					}
				} else {
					name = attr.Name.Local
				}
				el.Attributes = append(el.Attributes, &Attribute{
					Name:  name,
					Value: attr.Value,
				})
			}
			if e != nil {
				e.Children = append(e.Children, el)
			}
			e = el

			if doc.Root == nil {
				doc.Root = e
			}
		case xml.EndElement:
			e = e.Parent
		case xml.CharData:
			// text node
			if e != nil {
				e.Text = string(bytes.TrimSpace(token))
			}
		case xml.ProcInst:
			doc.ProcInst = stringifyProcInst(&token)
		case xml.Directive:
			doc.Directives = append(doc.Directives, stringifyDirective(&token))
		}

		// get the next token
		t, err = p.Token()
	}

	// Make sure that reading stopped on EOF
	if err != io.EOF {
		return nil, err
	}

	// All is good, return the document
	return doc, nil
}
