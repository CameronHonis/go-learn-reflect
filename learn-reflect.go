package main

import (
	"fmt"
	"reflect"
)

func AbsorbArgs(args []reflect.Value) []reflect.Value {
	fmt.Println("absorbed args: ", args)
	return []reflect.Value{reflect.ValueOf(true)}
}

type MyStruct struct {
	A   int
	Bee string
	See float64
}

func (s *MyStruct) Foo(a int) bool {
	fmt.Println(s.Bee)
	return false
}

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

	// Q. Difference between retrieving method on a Value vs a Type?
	println("")
	k := &MyStruct{A: 1, Bee: "bee", See: 3.14}
	kVal := reflect.ValueOf(k)
	kType := reflect.TypeOf(k)
	valMethod := kVal.Method(0)
	typeMethod := kType.Method(0)
	fmt.Println(valMethod, "|", valMethod.Type(), "|", valMethod.Kind(), "|", valMethod.String())
	fmt.Println(typeMethod, "|", typeMethod.Type, "|", typeMethod.Name, "|", typeMethod.Name)
	fmt.Println(typeMethod.Func)        // 0x4bb1e0
	fmt.Println(typeMethod.Func.Type()) // func(*main.MyStruct, int) bool
	//fmt.Println(valMethod == typeMethod) // does not compile, Value is not comparable to Method
	fmt.Println(valMethod == typeMethod.Func)               // false
	fmt.Println(valMethod.Type() == typeMethod.Func.Type()) // false
	valMethod.Call([]reflect.Value{reflect.ValueOf(12)})
	typeMethod.Func.Call([]reflect.Value{reflect.ValueOf(k), reflect.ValueOf(13)}) // requires passing in the struct
	// A. Main difference is that the Type method is an actual Method vs just a Value. A Method contains a pointer to
	// the underlying method as a Value in its field Func. This Value contains the type info that it's a method on a
	// struct, which the value version of the function pointer only knows the input and output typing. Additionally,
	// calling the valMethod will invoke the method with the underlying struct instance in memory, but calling the
	// typeMethod.Func will require a struct instance to be passed as the first argument (similar to python's "self")

	// Q. Does each struct instance generate each method on the struct at instantiation?
	l := &MyStruct{A: 1, Bee: "bee", See: 3.14}
	m := &MyStruct{A: 1, Bee: "bee", See: 3.14}
	lVal := reflect.ValueOf(l)
	mVal := reflect.ValueOf(m)
	fmt.Println(lVal.Method(0) == lVal.Method(0)) // the control, true
	fmt.Println(lVal.Method(0) == mVal.Method(0)) // false
	// A. Yes, each struct instance generates each method on the struct at instantiation.

	// Q. Does the Method derived from a struct Type point to a shared function definition?
	n := &MyStruct{A: 1, Bee: "bee", See: 3.14}
	o := &MyStruct{A: 1, Bee: "bee", See: 3.14}
	nType := reflect.TypeOf(n)
	oType := reflect.TypeOf(o)
	fmt.Println(nType.Method(0) == nType.Method(0))         // the control, false?
	fmt.Println(nType.Method(0).Func, nType.Method(0).Func) // the control, same address (but compares as false?)
	fmt.Println(nType.Method(0).Func, oType.Method(0).Func) // the same address (but compares as false?)
	// Yes, the Method derived from a struct Type points to a shared function definition.

	// Q. Can I compare reflect.ValueOf(struct.method) to reflect.ValueOf(struct).Method ?
	p := &MyStruct{A: 1, Bee: "bee", See: 3.14}
	pVal := reflect.ValueOf(p)
	pMethod := pVal.Method(0)
	fmt.Println(reflect.ValueOf(p.Foo), reflect.ValueOf(p.Foo))  // control, the same address (but compares as false?)
	fmt.Println(reflect.ValueOf(p.Foo), pMethod)                 // different addresses
	fmt.Println(reflect.ValueOf(p.Foo).Type() == pMethod.Type()) // true
	// A. No. Only the types remain equal.

	// Q. Can I stub methods using Value.Set()?
	q := &MyStruct{A: 1, Bee: "bee", See: 3.14}
	qVal := reflect.ValueOf(q)
	qValMethod := qVal.Method(0)
	qNewMethod := reflect.MakeFunc(qValMethod.Type(), AbsorbArgs)
	fmt.Println(qValMethod.CanInterface()) // true
	fmt.Println(qNewMethod.CanInterface()) // true
	fmt.Println(qValMethod.Type())         // func(int) bool
	//qValMethod.Set(qNewMethod)             // panic: reflect: reflect.Value.Set using unaddressable value

	qMethod := reflect.ValueOf(q.Foo)
	qNewMethod = reflect.MakeFunc(qMethod.Type(), AbsorbArgs)
	fmt.Println(qMethod.Kind()) // func
	//qMethod.Set(qMethod) // panic: reflect: reflect.Value.Set using unaddressable value

	qType := qVal.Type()
	qTypeMethod := qType.Method(0).Func
	qNewMethod = reflect.MakeFunc(qTypeMethod.Type(), AbsorbArgs)
	fmt.Println(qTypeMethod.Kind()) // func
	//qTypeMethod.Set(qNewMethod)     // panic: reflect: reflect.Value.Set using unaddressable value
	fmt.Println(qValMethod.CanSet(), qMethod.CanSet(), qTypeMethod.CanSet()) // false false false
	// A. No, calling Set() on a method results in an unaddressable value panic.

}
