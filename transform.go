package xmldom

import (
	"fmt"
	"bytes"
	"errors"
)

type TransformFunc func([]byte) ([]byte, error)

func Transform(n *Node, f TransformFunc) (*Node, error) {
	switch {
	case n == nil:
		return nil, errors.New("empty node")
	case f == nil:
		return nil, errors.New("empty transform func")
	}
	data, err := f([]byte(n.XML()))
	if err != nil {
		return nil, err
	}
	doc, err := Parse(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("parse transform xml error: %v", err)
	}
	root := doc.Root
	root.Document = n.Document
	root.Parent = n.Parent
	return root, err
}

func (n *Node) Transform(f TransformFunc) error {
	root, err := Transform(n, f)
	if err != nil {
		return err
	}
	*n = *root
	return nil
}
