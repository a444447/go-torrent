package Bencode

import (
	"bufio"
	"errors"
	"io"
)

type Btype uint8

//表示Bencode的四种基本类型
const (
	BtypeofINT  Btype = 0x01
	BtypeofSTR  Btype = 0x02
	BtypeofLIST Btype = 0x03
	BtypeofDICT Btype = 0x04
)

//因为Bvalue可以是任意的类型，所以用interface{}
type Bvalue interface{}

type Bobject struct {
	type_ Btype
	val_  Bvalue
}

func (o *Bobject) INT() (int, error) {
	if o.type_ != BtypeofINT {
		return 0, errors.New("type Error : Need int")
	}
	return o.val_.(int), nil
}

func (o *Bobject) STR() (string, error) {
	if o.type_ != BtypeofSTR {
		return "", errors.New("type Error : Need string")
	}
	return o.val_.(string), nil
}

func (o *Bobject) LIST() ([]*Bobject, error) {

	if o.type_ != BtypeofLIST {
		return nil, errors.New("type Error : Need list")
	}
	return o.val_.([]*Bobject), nil
}

func (o *Bobject) DICT() (map[string]*Bobject, error) {
	if o.type_ != BtypeofDICT {
		return nil, errors.New("type Error : Need dict")
	}
	return o.val_.(map[string]*Bobject), nil
}

func (o *Bobject) Bencode(w io.Writer) int {
	bw, ok := w.(*bufio.Writer)
	if !ok {
		bw = bufio.NewWriter(w)
	}
	wLen := 0
	switch o.type_ {
	case BtypeofINT:
		val, _ := o.INT()
		wLen += EncodeInt(bw, val)

	case BtypeofSTR:
		val, _ := o.STR()
		wLen += EncodeString(bw, val)

	case BtypeofLIST:
		val, _ := o.LIST()
		bw.WriteByte('l')
		wLen++
		for _, v := range val {
			wLen += v.Bencode(bw)
		}
		bw.WriteByte('e')
		wLen++

	case BtypeofDICT:
		val, _ := o.DICT()
		bw.WriteByte('d')
		wLen++
		for k, v := range val {
			wLen += EncodeString(bw, k)
			wLen += v.Bencode(bw)
		}
	}

	_ = bw.Flush()

	return wLen

}
