package utils

import (
	"fmt"
	"testing"
)

func TestInterfaceToString(t *testing.T) {
	m := map[string]interface{}{
		"mothers": "children",
	}
	println(InterfaceToString(1))
	c := InterfaceToString(4)
	fmt.Printf("%T\n", c)

	one := SafeGetFromInterface(1, "$")
	println(InterfaceToString(one))
	println(m)
}

func GetExternalServices() {

}
