package egothor

import (
	"encoding/binary"
	"io"
)

// errBinaryReader is a abstraction layer over binary.Read
type errBinaryReader struct {
	r   io.Reader
	err error
}

// Read reads binary from reader to data interface
func (e *errBinaryReader) Read(data interface{}) {
	if e.err != nil {
		return
	}
	e.err = binary.Read(e.r, binary.BigEndian, data)
}

// Err returns error
func (e *errBinaryReader) Err() error {
	return e.err
}
