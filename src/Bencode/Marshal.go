package Bencode

import (
	"io"
	"reflect"
	"strings"
)

func Marshal(w io.Writer, s interface{}) int {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem() // 判断s是否为指针，如果是指针则取它的值(方便处理)
	}

	return marshalValue(w, v)
}

func marshalValue(w io.Writer, v reflect.Value) int {
	len := 0
	switch v.Kind() {
	case reflect.String:
		len += EncodeString(w, v.String())
	case reflect.Int:
		len += EncodeInt(w, int(v.Int()))
	case reflect.Slice:
		len += marshalList(w, v)
	case reflect.Struct:
		len += marshalDict(w, v)
	}

	return len
}

func marshalList(w io.Writer, val reflect.Value) int {
	len := 2
	w.Write([]byte{'l'})
	for i := 0; i < val.Len(); i++ {
		eval := val.Index(i)
		len += marshalValue(w, eval) //对于List中每个具体的元素，递归调用
	}
	w.Write([]byte{'e'})
	return len
}

func marshalDict(w io.Writer, val reflect.Value) int {
	len := 2
	w.Write([]byte{'d'})
	for i := 0; i < val.NumField(); i++ {
		fv := val.Field(i)
		ft := val.Type().Field(i)
		len += EncodeString(w, strings.ToLower(ft.Name))
		len += marshalValue(w, fv)
	}
	w.Write([]byte{'e'})
	return len
}
