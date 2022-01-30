package org

import "io"

// Node is a interface that represents a parsed node of the document.
type Node interface {
	Write(w io.Writer) error
}
