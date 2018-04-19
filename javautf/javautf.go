package javautf // import "github.com/kreativka/gostempel/javautf"

import (
	"bytes"
	"encoding/binary"
	"io"
)

const surrogateMin = 0xD800
const surrogateMax = 0xDFFF

// ReadUTF datastream
func ReadUTF(reader io.Reader) (string, error) {
	var utfLen uint16
	err := binary.Read(reader, binary.BigEndian, &utfLen)
	if err != nil {
		return "", err
	}
	return readUTFBytes(reader, int(utfLen))
}

func readUTFBytes(reader io.Reader, utfLen int) (string, error) {
	byteArr := make([]byte, utfLen)
	_, err := reader.Read(byteArr)
	if err != nil {
		return "", err
	}
	var charArr bytes.Buffer
	charArr.Grow(utfLen / 2)
	var count int
	var c byte

	for count < utfLen {
		c = byteArr[count]
		if c > 127 {
			break
		}
		count++
		charArr.WriteByte(c)
	}
	var surrogateHi int32
	for count < utfLen {
		c = byteArr[count]
		msb4Bits := c >> 4
		if msb4Bits < 8 {
			count++
			charArr.WriteByte(c)
		} else if msb4Bits == 12 || msb4Bits == 13 {
			count += 2
			if count > utfLen {
				return charArr.String(), &errorString{"truncated character"}
			}
			c2 := byteArr[count-1]
			if c2&0xc0 != 0x80 {
				return charArr.String(), &errorString{"malformed input"}
			}
			r := rune((int32(c&0x1f) << 6) | int32(c2&0x3f))
			charArr.WriteRune(r)
		} else if msb4Bits == 14 {
			count += 3
			if count > utfLen {
				return charArr.String(), &errorString{"truncated character"}
			}
			c2 := byteArr[count-2]
			c3 := byteArr[count-1]
			if c2&0xc0 != 0x80 || c3&0xc0 != 0x80 {
				return charArr.String(), &errorString{"malformed input"}
			}
			r := rune((int32(c&0x0f) << 12) | (int32(c2&0x3f) << 6) | int32(c3&0x3f))
			if r < surrogateMin || r > surrogateMax {
				charArr.WriteRune(r)
			} else {
				if r < 0xDBFF {
					surrogateHi = r
				} else {
					surrogateLo := r
					r = (((surrogateHi - 0xD800) << 10) | (surrogateLo - 0xDC00)) + 0x10000
					charArr.WriteRune(r)
				}
			}
		} else {
			return charArr.String(), &errorString{"malformed input, maybe it's regular UTF-8?"}
		}
	}
	return charArr.String(), nil
}

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
