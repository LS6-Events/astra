package testfiles

// Example call expression with ranging number of arguments and return types

import (
	"fmt"

	"strings"

	"github.com/ls6-events/astra/astTraversal/testfiles/otherpkg1"
)

func NoArgs() {
	fmt.Println("Function with no arguments.")
}

// nolint:unused
func callExpr() {
	// Direct function calls
	fmt.Println("Hello, World!")

	// Variable declaration & method call on standard library types
	str := "LS6 Events"
	strings.Contains(str, "LS6")

	// Regular function call with multiple arguments
	TestFunction(5, "LS6")

	// Method calls on user-defined types
	ms := MyStruct{Name: "MyStruct"}
	ms.SayHello()

	// Call to a function with no arguments
	NoArgs()

	// Call to a function with a response
	result := TestFunction(10, "OpenAI")

	// Call to a function with more than 1 response
	result2, err := otherpkg1.GetFoo()

	fmt.Println(result, result2, err)
}

func TestFunction(a int, b string) string {
	return fmt.Sprintf("%d - %s", a, b)
}
