package main

import (
	"fmt"
	"reflect"
)

func a() {
	var x float64 = 3.14
	v := reflect.ValueOf(x)
	fmt.Println("type:", v.Type())
	fmt.Println("kind:", v.Kind())
	fmt.Println("kind is float64:", v.Kind() == reflect.Float64)
	fmt.Println("value:", v.Float())

	// Q: is type and kind comparable?
	//fmt.Println(v.Type() == v.Kind())
	// A: no, kind is a uint and type is a struct

	// Q: is reflect.TypeOf() comparable to reflect.ValueOf().Type()?
	fmt.Println(reflect.TypeOf(x) == v.Type())
	// A: yes, the expression evaluates to true

	// Q: if a max uint8 is modified as a Value, does adding 1 to it overflow?
	//var a uint8 = 255
	//fmt.Println(a + 1)
	//aVal := reflect.ValueOf(a)
	//aVal.Add
	//aVal.SetUint(123)
	// A: no, you can't add to a Value - only Set_ methods are available

	// Q: does Type and/or Value preserve user-defined types?
	type MyInt int
	var b MyInt = 7
	fmt.Println(reflect.TypeOf(b))
	fmt.Println(reflect.TypeOf(b).String())
	fmt.Println(reflect.TypeOf(b).Kind())
	fmt.Println(reflect.ValueOf(b).Type())
	fmt.Println(reflect.ValueOf(b).Type().String())
	fmt.Println(reflect.ValueOf(b).Type().Kind())
	// A: yes, both Type and Value preserve user-defined types. But printing the `Kind` of a user-defined type returns
	// the underlying type, not the user-defined type.

	// Q: can you use reflect.ValueOf() to modify a user-defined type?
	//var c MyInt = 7
	//cVal := reflect.ValueOf(c)
	//cVal.Set(reflect.ValueOf(MyInt(123)))
	//cVal.SetInt(123)
	// A: no, you can't use reflect.ValueOf() to modify a user-defined type, because c is not settable.

	// Q: can you use reflect.ValueOf() to modify value of a user-defined type?
	var d MyInt = 7
	dVal := reflect.ValueOf(&d)
	dVal.Elem().Set(reflect.ValueOf(MyInt(123)))
	fmt.Println(d)
	dVal.Elem().SetInt(124)
	fmt.Println(d)
	//dVal.SetInt(125)
	//fmt.Println(d)
	// A: yes. Its not even necessary to cast the int to the user-defined type. But it is required to use Elem() to get
	// the value of the pointer to the user-defined type.

	// Q. How do primitive pointers get represented when printed in various ways?
	e := true
	ePtr := &e
	ePtrVal := reflect.ValueOf(ePtr)
	fmt.Println(ePtrVal)          // 0xc000020108
	fmt.Println(ePtrVal.Type())   // *bool
	fmt.Println(ePtrVal.Kind())   // ptr
	fmt.Println(ePtrVal.String()) // <*bool Value>

	// Q. How do structs get represented when printed in various ways?
	type MyStruct struct {
		A   int
		Bee string
		See float64
	}
	f := MyStruct{A: 1, Bee: "bee", See: 3.14}
	fVal := reflect.ValueOf(f)
	fmt.Println(fVal)          // {1 bee 3.14}
	fmt.Println(fVal.Type())   // main.MyStruct
	fmt.Println(fVal.Kind())   // struct
	fmt.Println(fVal.String()) // <main.MyStruct Value>

	// Q. How do pointers to structs get represented when printed in various ways?
	g := &MyStruct{A: 1, Bee: "bee", See: 3.14}
	gVal := reflect.ValueOf(g)
	fmt.Println(gVal)          // &{1 bee 3.14}
	fmt.Println(gVal.Type())   // *main.MyStruct
	fmt.Println(gVal.Kind())   // ptr
	fmt.Println(gVal.String()) // <*main.MyStruct Value>

	// Q. Is there a difference between passing a pointer variable vs the address to the value in reflect.ValueOf()?
	h := MyStruct{A: 1, Bee: "bee", See: 3.14}
	hVal := reflect.ValueOf(&h)
	fmt.Println(hVal)          // &{1 bee 3.14}
	fmt.Println(hVal.Type())   // *main.MyStruct
	fmt.Println(hVal.Kind())   // ptr
	fmt.Println(hVal.String()) // <*main.MyStruct Value>
	// A. No

	// Q. How does reflect.ValueOf() handle nil values?
	i := (*MyStruct)(nil)
	iVal := reflect.ValueOf(i)
	fmt.Println(iVal)          // <nil>
	fmt.Println(iVal.Type())   // *main.MyStruct
	fmt.Println(iVal.Kind())   // ptr
	fmt.Println(iVal.String()) // <*main.MyStruct Value>
	fmt.Println(iVal.Elem())   // <invalid reflect.Value>
	//fmt.Println(iVal.Elem().Type()) // panic: reflect: call of reflect.Value.Type on zero Value
	// A. Getting the type from an invalid reflect.Value panics

	// Q. Can I retrieve the function signature of a nil func pointer?
	var j func(int) int
	jVal := reflect.ValueOf(j)
	fmt.Println(jVal)          // <invalid reflect.Value>
	fmt.Println(jVal.String()) // <invalid reflect.Value>
	fmt.Println(jVal.Type())   // func(int) int
	// A. Yes, this works differently than nil pointers to structs.

}
