package main

import (
	"fmt"
	"reflect"
)

func main(){
	demo3()
}


// demo1 转化成reflect类型
func demo1(){
	author := "testdemo1"

	// 1. 转化成reflect类型
	fmt.Println("Typeof : ", reflect.TypeOf(author))
	fmt.Println("Valueof : ", reflect.ValueOf(author))

	v := reflect.ValueOf(author)
	fmt.Println("original value:", v.Interface().(string))
}

// demo2 error operation
func demo2(){
	i := 1
	// 值传递
	v := reflect.ValueOf(i)
	v.SetInt(10)
	fmt.Println(i)
}

func demo3(){
	i := 1
	// 值传递，传入的是指针；获取指针对应的 reflect.Value类型
	v := reflect.ValueOf(&i)
	// reflect.Value.Elem: 获取指针指向的变量（reflect.Value类型）
	// reflect.Value.SetInt 更新变量的值
	v.Elem().SetInt(10)
	fmt.Println(i)
}