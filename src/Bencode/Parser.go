package Bencode

import (
	"bufio"
	"io"
)

func Parse(r io.Reader) (*Bobject, error) {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}
	b, err := br.Peek(1) //Peek读取n字节但是不消耗,返回[]byte
	if err != nil {
		return nil, err
	}
	var Bret Bobject
	switch {
	case b[0] >= '0' && b[0] <= '9': //string类型
		val, err := DecodeString(br)
		if err != nil {
			return nil, err
		}
		Bret.type_ = BtypeofSTR
		Bret.val_ = val

	case b[0] == 'i': //判断为int
		val, err := DecodeInt(br)
		if err != nil {
			return nil, err
		}
		Bret.type_ = BtypeofINT
		Bret.val_ = val

	case b[0] == 'l': //判断为list
		_, _ = br.ReadByte() //先消耗第一个'l'
		var list []*Bobject
		for { //循环体内处理Bobject
			if p, _ := br.Peek(1); p[0] == 'e' {
				br.ReadByte()
				break
			}
			e, err := Parse(br)
			if err != nil {
				return nil, err
			}
			list = append(list, e)
		}
		Bret.type_ = BtypeofLIST
		Bret.val_ = list
	case b[0] == 'd':
		br.ReadByte() //消耗掉第一个d
		dict := make(map[string]*Bobject)
		for {
			if p, _ := br.Peek(1); p[0] == 'e' {
				br.ReadByte()
				break
			}
			//先处理key
			key, err := DecodeString(br)
			if err != nil {
				return nil, err
			}
			//再处理val
			val, err := Parse(br)
			if err != nil {
				return nil, err
			}
			dict[key] = val
		}

		Bret.type_ = BtypeofDICT
		Bret.val_ = dict
	}

	return &Bret, err
}
