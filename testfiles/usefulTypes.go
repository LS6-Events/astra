package testfiles

import (
	"fmt"
	"github.com/ls6-events/gengo/testfiles/otherpkg1"
	"strings"
)

type MyStruct struct {
	Name string
}

func (m *MyStruct) SayHello() {
	fmt.Println("Hello from", strings.Join([]string{"MyStruct", m.Name}, " "))
}

func (m *MyStruct) ExternalPackage() {
	_ = otherpkg1.Foo{}
}
