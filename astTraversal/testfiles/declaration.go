package testfiles

import "fmt"

var (
	// MyVar1 is a string variable.
	MyVar1 = "MyVar1"

	// MyVar2 is a float variable.
	MyVar2 = 3.14
)

const (
	// MyConst1 is a constant.
	MyConst1 = "MyConst1"

	// MyConst2 is a constant.
	MyConst2 = 1234
)

func MyFunc1() {
	assignStmt := "Hello World"
	var1, var2 := "var1", 123

	fmt.Print(assignStmt)
	fmt.Print(var1, var2)
}

func MyFunc2(param1 MyStruct, param2 MyInt) (string, *MyStruct) {
	return fmt.Sprintf("%s %d", param1.Name, param2), nil
}

var MyFunc3 = func() (int, error) {
	return 42, nil
}
