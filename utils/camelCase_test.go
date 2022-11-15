package utils

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"testing"
)
import "github.com/gobeam/stringy"

func TestToCamelCaseLower(t *testing.T) {
	//:= "TheyDontKnow"
	s := stringy.New("HelloThere??")
	r := s.CamelCase("?", "")
	fmt.Println(r) // HelloManHowAreYou

	println(r)

	str := stringy.New("hello__man how-Are you??")
	result := str.CamelCase("?", "")
	fmt.Println(result) // HelloManHowAreYou

	snakeStr := str.SnakeCase("?", "")
	fmt.Println(snakeStr.ToLower()) // hello_man_how_are_you

	kebabStr := str.KebabCase("?", "")
	fmt.Println(kebabStr.ToUpper()) // HELLO-MAN-HOW-ARE-YOU

	rr := strcase.ToLowerCamel("HelloThere")
	println(rr)
}
