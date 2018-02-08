package javautf

import (
	"testing"
	"bytes"
	"strings"
	"encoding/hex"
)

func TestAscii(t *testing.T) {
	inp := "hello"
	result, err := readUTFBytes(bytes.NewReader([]byte(inp)), len(inp))
	if err != nil {
		t.Error("expected", inp, " error ", err)
	}
	if strings.Compare(result, inp) != 0 {
		t.Error("expected", inp, "got", result)
	}
}

func TestCyr(t *testing.T) {
	inp := "привет"
	inpBytes, ok := hex.DecodeString("000CD0BFD180D0B8D0B2D0B5D182")
	if ok != nil {
		t.Error("bad hex")
	}
	result, err := ReadUTF(bytes.NewReader(inpBytes))
	if err != nil {
		t.Error("expected", inp, " error ", err)
	}
	if strings.Compare(result, inp) != 0 {
		t.Error("expected", inp, "got", result)
	}
}

func Test3Bytes(t *testing.T) {
	inp := "法轮功"
	inpBytes, ok := hex.DecodeString("0009E6B395E8BDAEE58A9F")
	if ok != nil {
		t.Error("bad hex")
	}
	result, err := ReadUTF(bytes.NewReader(inpBytes))
	if err != nil {
		t.Error("expected", inp, " error ", err)
	}
	if strings.Compare(result, inp) != 0 {
		t.Error("expected", inp, "got", result)
	}
}
func TestDollarEuro(t *testing.T) {
	inp := "$€"
	inpBytes, ok := hex.DecodeString("000424E282AC")
	if ok != nil {
		t.Error("bad hex")
	}
	result, err := ReadUTF(bytes.NewReader(inpBytes))
	if err != nil {
		t.Error("expected", inp, " error ", err)
	}
	if strings.Compare(result, inp) != 0 {
		t.Error("expected", inp, "got", result)
	}
}

func TestSurrogate(t *testing.T) {
	inp := "𤭢"
	inpBytes, ok := hex.DecodeString("0006EDA192EDBDA2")
	if ok != nil {
		t.Error("bad hex")
	}
	result, err := ReadUTF(bytes.NewReader(inpBytes))
	if err != nil {
		t.Error("expected", inp, " error ", err)
	}
	if strings.Compare(result, inp) != 0 {
		t.Error("expected", inp, "got", result)
	}
}

func TestAll(t *testing.T) {
	inp := "𤭢0€轮功한국йЯz-\x00𐐷"
	inpBytes, ok := hex.DecodeString("0024EDA192EDBDA230E282ACE8BDAEE58A9FED959CEAB5ADD0B9D0AF7A2DC080EDA081EDB0B7")
	if ok != nil {
		t.Error("bad hex")
	}
	result, err := ReadUTF(bytes.NewReader(inpBytes))
	if err != nil {
		t.Error("expected", inp, " error ", err)
	}
	if strings.Compare(result, inp) != 0 {
		t.Error("expected", inp, "got", result)
	}
}

func TestAllRegularUtf8(t *testing.T) {
	inpBytes, ok := hex.DecodeString("F0A4ADA230E282ACE8BDAEE58A9FED959CEAB5ADD0B9D0AF7A2D00F09090B7")
	if ok != nil {
		t.Error("bad hex")
	}
	_, err := readUTFBytes(bytes.NewReader(inpBytes), len(inpBytes))
	if err == nil {
		t.Error("expected to throw error since it's regular UTF created by str.getBytes(StandardCharsets.UTF_8)")
	}
	if strings.Compare(string(inpBytes), "𤭢0€轮功한국йЯz-\x00𐐷") != 0 {
		t.Error("unexpected input")
	}
}
