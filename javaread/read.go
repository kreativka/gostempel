package javaread

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Int returns int32 from reader
func Int(r io.Reader) int32 {
	var v int32
	err := binary.Read(r, binary.BigEndian, &v)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	return v
}

// Bool returns boolean from reader
func Bool(r io.Reader) bool {
	var v bool
	err := binary.Read(r, binary.BigEndian, &v)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	return v
}

// Char returns char(2bytes) from reader
func Char(r io.Reader) rune {
	var v int16
	err := binary.Read(r, binary.BigEndian, &v)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	return rune(v)
}
