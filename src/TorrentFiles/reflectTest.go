package main

import (
	"fmt"
	"reflect"
)

//反射提供一种机制,也就是可以在代码运行时获得变量的类型和值
//reflect.ValueOf可以得到reflect.Value
//reflect.TypeOf可以得到reflect.Type

type Addr struct {
	Province  string
	City      string
	Telephone string
}

func main() {
	addr := Addr{
		Province:  "SICHUAN",
		City:      "SUINING",
		Telephone: "12345",
	}

	var obj interface{} = &addr
	reflectValue := reflect.ValueOf(obj)
	v := reflectValue.Elem()
	//reflectType := reflect.TypeOf(addr)
	fmt.Printf("%T\n", v)
	fmt.Println(reflectValue.Type())
}
