package utils

import (
	"fmt"
	"strings"
)

func IsNullOrEmpty(i interface{}) bool {

	if i == nil {
		return true
	}

	if strings.TrimSpace(fmt.Sprintf("%v", i)) == "" {
		return true
	}

	if strings.TrimSpace(fmt.Sprintf("%v", i)) == "{}" {
		return true
	}

	// map[string]interface{}{}
	if fmt.Sprintf("%v", i) == "map[]" { // result is map[] and length is 5
		return true
	}

	str := strings.Fields(fmt.Sprintf("%v", i))

	if len(str) == 0 {
		return true
	}

	if len(str) == 2 {
		if str[0] == "{" && str[1] == "}" {
			return true
		}
	}

	return false
}
