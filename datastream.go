package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

// readInt returns int32 from reader
func readInt(r io.Reader) int32 {
	var v int32
	err := binary.Read(r, binary.BigEndian, &v)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	return v
}

// readBool returns boolean from reader
func readBool(r io.Reader) bool {
	var v bool
	err := binary.Read(r, binary.BigEndian, &v)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	return v
}

//FIXME
// readChar returns char(2bytes) from reader
func readChar(r io.Reader) rune {
	var v int16
	err := binary.Read(r, binary.BigEndian, &v)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	//FIXME!!!
	return rune(v)
}
