package transport

import "io"

// Interface represents the transport-interface
type Interface interface {

	// we support reader/writer interface
	io.Reader
	io.Writer
}
