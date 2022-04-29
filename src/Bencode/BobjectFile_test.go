package Bencode

import (
	"bytes"
	"testing"
)

const checkMark = "\u2713"
const ballotX = "\u2717"

func TestINT(t *testing.T) {
	ObjectInt := Bobject{
		type_: BtypeofINT,
		val_:  2200,
	}

	res, err := ObjectInt.INT()
	if err != nil {
		t.Fatal("\t\tFalse\t", ballotX, err)
	}
	t.Log("\t\tSuccess\t", res, checkMark)
}

func TestBobject_STR(t *testing.T) {
	ObjectInt := Bobject{
		type_: BtypeofSTR,
		val_:  "hello",
	}

	res, err := ObjectInt.STR()
	if err != nil {
		t.Fatal("\t\tFalse\t", ballotX, err)
	}
	t.Log("\t\tSuccess\t", res, checkMark)
}

func TestBencode(t *testing.T) {
	ObjectSTR := Bobject{
		type_: BtypeofSTR,
		val_:  "hello",
	}

	ObjectInt := Bobject{
		type_: BtypeofINT,
		val_:  2200,
	}

	ObjectList := Bobject{
		type_: BtypeofLIST,
		val_:  []*Bobject{&ObjectInt, &ObjectSTR},
	}
	buf := new(bytes.Buffer)
	LEN := ObjectList.Bencode(buf)
	if LEN != 15 {
		t.Fatal("\tBencode Fail ", LEN, ballotX)
	}
	t.Log("\tBencode Success\n\t\t ", buf, checkMark)
}
