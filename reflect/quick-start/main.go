package main

import (
	"fmt"
	"reflect"
)

type Enum int

const (
	Zero Enum = 0
)

type cat struct {
	Name string
	// 带有结构体tag的字段
	Type int `json:"type" id:"100"`
}

func main() {
	var a int
	typeOfA := reflect.TypeOf(a)
	fmt.Println(typeOfA.Name(), typeOfA.Kind())

	// 获取结构体实例的反射类型对象
	typeOfCat := reflect.TypeOf(cat{})
	// 显示反射类型对象的名称和种类
	fmt.Println(typeOfCat.Name(), typeOfCat.Kind())

	// 获取Zero常量的反射类型对象
	typeOfZ := reflect.TypeOf(Zero)
	// 显示反射类型对象的名称和种类
	fmt.Println(typeOfZ.Name(), typeOfZ.Kind())

	ptr()
	fmt.Println("-------")

	filed()
}

func ptr() {
	// 创建cat的实例
	ins := &cat{}
	// 获取结构体实例的反射类型对象
	typeOfCat := reflect.TypeOf(ins)
	// 显示反射类型对象的名称和种类
	// 指针变量的类型名称是空
	fmt.Printf("name:'%v' kind:'%v'\n", typeOfCat.Name(), typeOfCat.Kind())
	// 取指针类型的元素
	typeOfCat = typeOfCat.Elem()
	// 显示反射类型对象的名称和种类
	fmt.Printf("element name: '%v', element kind: '%v'\n", typeOfCat.Name(), typeOfCat.Kind())
}

func filed() {
	// 创建cat的实例
	ins := cat{Name: "mimi", Type: 1}
	// 获取结构体实例的反射类型对象
	typeOfCat := reflect.TypeOf(ins)

	// 遍历结构体所有成员
	for i := 0; i < typeOfCat.NumField(); i++ {
		// 获取每个成员的结构体字段类型
		fieldType := typeOfCat.Field(i)
		// 输出成员名和tag
		fmt.Printf("name: %v  tag: '%v'\n", fieldType.Name, fieldType.Tag)
	}

	// 通过字段名, 找到字段类型信息
	if catType, ok := typeOfCat.FieldByName("Type"); ok {
		// 从tag中取出需要的tag
		fmt.Println(catType.Tag.Get("json"), catType.Tag.Get("id"))
	}

}
