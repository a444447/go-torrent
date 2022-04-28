package Bencode

import "testing"

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
