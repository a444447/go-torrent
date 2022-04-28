package Bencode

import (
	"bytes"
	"testing"
)

func TestInt(t *testing.T) {
	buf := new(bytes.Buffer)
	val := 2200
	wLen := EncodeInt(buf, val)
	if wLen != 6 {
		t.Fatal("\tEncode Int Error :", wLen, ballotX)
	}
	t.Log("\tEncode Result:", buf)
	decodeVal, _ := DecodeInt(buf)
	if decodeVal != val {
		t.Fatal("\tDecode Int Fail :", decodeVal, ballotX)
	}
	t.Log("\tTest Int Success", decodeVal, checkMark)
}

func TestString(t *testing.T) {
	buf := new(bytes.Buffer)

	val := "a44447"
	//valLen := len(val)
	wLen := EncodeString(buf, val)
	if wLen != 8 {
		t.Fatal("\tEncode String Fail:", buf, ballotX)
	}
	t.Log("Encode String :", buf)
	decodeVar, _ := DecodeString(buf)
	if decodeVar != val {
		t.Fatal("\t Decode String Fail:", decodeVar, ballotX)
	}
	t.Log("\t Success:", decodeVar, checkMark)
}
