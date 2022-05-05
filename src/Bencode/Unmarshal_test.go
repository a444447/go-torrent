package Bencode

import (
	"bytes"
	"fmt"
	"testing"
)

type User struct {
	Name string `bencode:"name"`
	Age  int    `bencode:"age"`
}

func TestUnmarshalList(t *testing.T) {
	str := "li85ei90ei95ee"
	l := &[]int{}
	Unmarshal(bytes.NewBufferString(str), l)
	if len(*l) != 3 {
		t.Fatal("FALSE")
	}
	fmt.Println(*l)
}

func TestUnmarshalUser(t *testing.T) {
	str := "d4:name6:archer3:agei29ee"
	u := &User{}
	Unmarshal(bytes.NewBufferString(str), u)
	fmt.Println(u.Name)
}
