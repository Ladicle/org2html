package org

import "io"

// Node is a interface that represents a parsed node of the document.
type Node interface {
	Write(w io.Writer) error
}

func Write(nodes []Node, out io.Writer) error {
	for i := range nodes {
		if err := nodes[i].Write(out); err != nil {
			return err
		}
	}
	return nil
}
