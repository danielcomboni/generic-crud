package utils

import (
	"fmt"
	"github.com/Jeffail/gabs"
	"log"
)

func ConvertInterfaceToMapOfStringKey(i interface{}) map[string]interface{} {
	var m map[string]interface{}
	m = i.(map[string]interface{})
	return m
}

func AlterDynamicProperty(prevDynamicValue []byte, dynamicProperty interface{}) map[string]interface{} {
	log.Println("altering dynamic property")
	fmt.Println("init")
	fmt.Println(string(prevDynamicValue))
	fmt.Println("last")
	fmt.Println(dynamicProperty)
	dpParsed, _ := gabs.ParseJSON(prevDynamicValue)
	fmt.Println(dpParsed.Data())

	if dpParsed.Data() != nil {
		log.Println("previous dynamic property is NOT EMPTY")
		prevDp := ConvertInterfaceToMapOfStringKey(dpParsed.Data())
		incomingDp := ConvertInterfaceToMapOfStringKey(dynamicProperty)
		for key, value := range incomingDp {
			prevDp[key] = value
		}
		return prevDp
	}
	log.Println("previous dynamic property is EMPTY")
	return ConvertInterfaceToMapOfStringKey(dynamicProperty)
}
