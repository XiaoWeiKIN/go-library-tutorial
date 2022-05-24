package main

import (
	"fmt"
	"reflect"
)

type MyFloat64 float64

func main() {
	var x float64 = 3.4
	fmt.Println("type:", reflect.TypeOf(x))
	v := reflect.ValueOf(x)
	fmt.Println("value:", v)
	fmt.Println("type:", v.Type())
	fmt.Println("kind is float64:", v.Kind() == reflect.Float64)
	fmt.Println("value:", v.Float())

	fmt.Println("MyFloat-----------")

	var k MyFloat64 = 3.5
	v1 := reflect.ValueOf(k)
	fmt.Println("kind is float64:", v1.Kind() == reflect.Float64)

	fmt.Println("SetInt-----------")
	Set()
	fmt.Println("SetStruct-----------")
	setStruct()
	fmt.Println("create-----------")
	create()
	fmt.Println("createMap-----------")
	createMap()
	fmt.Println("callFunc-----------")
	callFunc()
}

func Set() {
	var x float64 = 3.4
	p := reflect.ValueOf(&x) // Note: take the address of x.
	v := p.Elem()
	fmt.Println("settability of v:", v.CanSet())
	v.SetFloat(7.1)
	fmt.Println(v.Interface().(float64))
	fmt.Println(x)

}

type T struct {
	A int
	B string
}

func setStruct() {
	t := T{23, "skidoo"}
	s := reflect.ValueOf(&t).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
	s.Field(0).SetInt(77)
	s.Field(1).SetString("Sunset Strip")
	fmt.Println("t is now", t)
}

func create() {
	var t T
	typeOfA := reflect.TypeOf(t)
	// 创建这个类型的实例值，值以 reflect.Value 类型返回
	aIns := reflect.New(typeOfA)
	// aIns 的类型为 *main.T，种类为指针
	fmt.Println(aIns.Type(), aIns.Kind())
}

func createMap() {
	var m1 map[string]int
	var m2 = &m1
	v2 := reflect.ValueOf(m2)

	var key = "key1"
	var value = 123
	//
	mapType := reflect.MapOf(reflect.TypeOf(key), reflect.TypeOf(value))

	mapValue := reflect.MakeMap(mapType)
	mapValue.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
	mapValue.SetMapIndex(reflect.ValueOf("2"), reflect.ValueOf(1))

	keys := mapValue.MapKeys()
	for _, k := range keys {
		ck := k.Convert(mapValue.Type().Key())
		cv := mapValue.MapIndex(ck)
		fmt.Println(ck, cv)

	}

	v2.Elem().Set(mapValue)

	fmt.Println(m2)

}

func add(a, b int) int {
	return a + b
}

func callFunc() {
	// 将函数包装为反射值对象
	funcValue := reflect.ValueOf(add)
	// 构造函数参数, 传入两个整型值
	paramList := []reflect.Value{reflect.ValueOf(10), reflect.ValueOf(20)}
	// 反射调用函数
	retList := funcValue.Call(paramList)
	// 获取第一个返回值, 取整数值
	fmt.Println(retList[0].Int())
}
