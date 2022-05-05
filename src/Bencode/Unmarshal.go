package Bencode

import (
	"errors"
	"io"
	"reflect"
	"strings"
)

// Unmarshal r是我们读入的buffer, s是我们传入的结构，我们要将值填在其中的field中,因此需要传入指针
func Unmarshal(r io.Reader, s interface{}) error {
	o, err := Parse(r) //将buffer中的文本流 处理为Bobject
	if err != nil {
		return errors.New("parse Error")
	}
	p := reflect.ValueOf(s)      //传入的s必须是一个指针，返回一个s类型的reflect.value
	if p.Kind() != reflect.Ptr { //判断s是否是指针
		return errors.New("second parameter must be a pointer")
	}
	switch o.type_ { //判断我们得到的Bobject类型
	case BtypeofLIST:
		list, _ := o.LIST()
		l := reflect.MakeSlice(p.Elem().Type(), len(list), len(list)) // p.Elem()类似于指针取值的*操作，p.Elem().Type()得到的是 s{}的类型
		p.Elem().Set(l)                                               //将新建的slice l 分配给p
		err := unmarshalList(p, list)                                 //我们用unmarshalList将list的中的值设置到p中
		if err != nil {
			return err
		}
	case BtypeofDICT:
		dict, _ := o.DICT()
		err = unmarshalDict(p, dict)
		if err != nil {
			return err
		}
	default:
		return errors.New("first parameter must be struct or slice")
	}

	return nil
}

func unmarshalList(p reflect.Value, list []*Bobject) error {
	if p.Kind() != reflect.Ptr || p.Elem().Type().Kind() != reflect.Slice { //判断p是否是指针
		return errors.New("second parameter must be pointer to slice")
	}

	v := p.Elem() // 取值操作
	if len(list) == 0 {
		return nil
	}
	switch list[0].type_ {
	case BtypeofSTR:
		for i, o := range list {
			val, err := o.STR()
			if err != nil {
				return err
			}
			v.Index(i).SetString(val)
		}
	case BtypeofINT:
		for i, o := range list {
			val, err := o.INT()
			if err != nil {
				return err
			}
			v.Index(i).SetInt(int64(val))
		}
	case BtypeofLIST:
		for i, o := range list {
			val, err := o.LIST()
			if err != nil {
				return err
			}
			if v.Type().Elem().Kind() != reflect.Slice {
				return errors.New("type error")
			}
			lp := reflect.New(v.Type().Elem())
			ls := reflect.MakeSlice(v.Type().Elem(), len(val), len(val))
			lp.Elem().Set(ls)
			err = unmarshalList(lp, val)
			if err != nil {
				return err
			}
			v.Index(i).Set(lp.Elem())
		}
	case BtypeofDICT:
		for i, o := range list {
			val, err := o.DICT()
			if err != nil {
				return err
			}
			if v.Type().Elem().Kind() != reflect.Struct {
				return errors.New("type error")
			}
			dp := reflect.New(v.Type().Elem())
			err = unmarshalDict(dp, val)
			if err != nil {
				return err
			}
			v.Index(i).Set(dp.Elem())
		}
	}

	return nil
}

func unmarshalDict(p reflect.Value, dict map[string]*Bobject) error {
	if p.Kind() != reflect.Ptr || p.Elem().Type().Kind() != reflect.Struct {
		return errors.New("first parameter must be pointer")
	}
	v := p.Elem()
	for i, n := 0, v.NumField(); i < n; i++ {
		fv := v.Field(i)  //取struct的第i个field
		if !fv.CanSet() { //判断该field能否被set
			continue
		}
		ft := v.Type().Field(i)
		key := ft.Tag.Get("bencode")
		if key == "" {
			key = strings.ToLower(ft.Name)
		}
		fo := dict[key] //fo是Bobject类型
		if fo == nil {
			continue
		}
		switch fo.type_ {
		case BtypeofSTR:
			if ft.Type.Kind() != reflect.String {
				continue
			}
			val, _ := fo.STR()
			fv.SetString(val)
		case BtypeofINT:
			if ft.Type.Kind() != reflect.Int {
				break
			}
			val, _ := fo.INT()
			fv.SetInt(int64(val))
		case BtypeofLIST:
			if ft.Type.Kind() != reflect.Slice {
				break
			}
			list, _ := fo.LIST()
			lp := reflect.New(ft.Type)
			ls := reflect.MakeSlice(ft.Type, len(list), len(list))
			lp.Elem().Set(ls)
			err := unmarshalList(lp, list)
			if err != nil {
				break
			}
			fv.Set(lp.Elem())
		case BtypeofDICT:
			if ft.Type.Kind() != reflect.Struct {
				break
			}
			dp := reflect.New(ft.Type)
			dict, _ := fo.DICT()
			err := unmarshalDict(dp, dict)
			if err != nil {
				break
			}
			fv.Set(dp.Elem())
		}
	}
	return nil
}
