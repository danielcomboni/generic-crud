package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"log"
)

func InterfaceToString(i interface{}) string {
	marshal, _ := json.Marshal(i)
	return string(marshal)
}

func AnyToString(i any) string {
	marshal, _ := json.Marshal(i)
	return string(marshal)
}

func SafeGetFromInterfaceGeneric[t any](i interface{}, selector string) any {

	defer func() {

		if err := recover(); err != nil {
			log.Println("invalid selector: " + selector)
		}

	}()

	jsonString := oj.JSON(i)
	obj, err := oj.ParseString(jsonString)

	if err != nil {
		log.Println("failed to parse incoming jsonString for SafeGet operations")
		log.Println("jsonString: " + jsonString)
		log.Println("selector: " + selector)
		log.Println("err: " + err.Error())

		return *new(t)
	} else {
		expression, err := jp.ParseString(selector)
		if err != nil {
			log.Println("failed to parse selector")
			log.Println("selector: " + selector)
			log.Println("err: " + err.Error())
			return *new(t)
		}

		data := expression.Get(obj)

		if IsNullOrEmpty(data) {
			return *new(t)
		}

		if IsNullOrEmpty(data[0]) {
			return *new(t)
		}
		return data[0]
	}

}

func SafeGetFromInterfaceGenericAndDeserialize[T any](i interface{}, selector string) T {

	defer func() {

		if err := recover(); err != nil {
			log.Println("invalid selector: " + selector)
		}

	}()

	jsonString := oj.JSON(i)
	obj, err := oj.ParseString(jsonString)

	if err != nil {
		log.Println("failed to parse incoming jsonString for SafeGet operations")
		log.Println("jsonString: " + jsonString)
		log.Println("selector: " + selector)
		log.Println("err: " + err.Error())

		return *new(T)
	} else {
		expression, err := jp.ParseString(selector)
		if err != nil {
			log.Println("failed to parse selector")
			log.Println("selector: " + selector)
			log.Println("err: " + err.Error())
			return *new(T)
		}

		data := expression.Get(obj)

		if IsNullOrEmpty(data) {
			return *new(T)
		}

		if IsNullOrEmpty(data[0]) {
			return *new(T)
		}
		marshal, _ := json.Marshal(data[0])
		var t T
		_ = json.Unmarshal(marshal, &t)
		return t
	}

}

func SafeGetFromInterface(i interface{}, selector string) interface{} {

	defer func() {

		if err := recover(); err != nil {
			log.Println("invalid selector: " + selector)
		}

	}()

	jsonString := oj.JSON(i)
	obj, err := oj.ParseString(jsonString)

	if err != nil {
		log.Println("failed to parse incoming jsonString for SafeGet operations")
		log.Println("jsonString: " + jsonString)
		log.Println("selector: " + selector)
		log.Println("err: " + err.Error())
		return nil
	} else {
		expression, err := jp.ParseString(selector)
		if err != nil {
			log.Println("failed to parse selector")
			log.Println("selector: " + selector)
			log.Println("err: " + err.Error())
			return nil
		}

		data := expression.Get(obj)

		if IsNullOrEmpty(data) {
			return nil
		}

		if IsNullOrEmpty(data[0]) {
			return nil
		}
		return data[0]
	}

}

func SafeGetFromInterfaceErrorCaught(i interface{}, selector string) (interface{}, error) {

	jsonString := oj.JSON(i)
	obj, err := oj.ParseString(jsonString)

	if err != nil {
		log.Println("failed to parse incoming jsonString for SafeGet operations")
		log.Println("jsonString: " + jsonString)
		log.Println("selector: " + selector)
		log.Println("err: " + err.Error())
		return nil, err
	} else {
		expression, err := jp.ParseString(selector)
		if err != nil {
			log.Println("failed to parse selector")
			log.Println("selector: " + selector)
			log.Println("err: " + err.Error())
			return nil, err
		}
		data := expression.Get(obj)
		return data[0], errors.New("")
	}

}

func SafeGet(jsonString string, selector string) interface{} {

	obj, err := oj.ParseString(jsonString)

	if err != nil {
		log.Println("failed to parse incoming jsonString for SafeGet operations")
		log.Println("jsonString: " + jsonString)
		log.Println("selector: " + selector)
		log.Println("err: " + err.Error())
		return nil
	} else {
		expression, err := jp.ParseString(selector)
		if err != nil {
			log.Println("failed to parse selector")
			log.Println("selector: " + selector)
			log.Println("err: " + err.Error())
			return nil
		}
		data := expression.Get(obj)
		if IsNullOrEmpty(data) {
			return nil
		}

		if IsNullOrEmpty(data[0]) {
			return nil
		}
		return data[0]
	}

}

func SafeGetMarshalled(jsonString string, selector string) []byte {

	obj, err := oj.ParseString(jsonString)

	if err != nil {
		log.Println("failed to parse incoming jsonString for SafeGet operations")
		log.Println("jsonString: " + jsonString)
		log.Println("selector: " + selector)
		log.Println("err: " + err.Error())
		return nil
	} else {
		expression, err := jp.ParseString(selector)
		if err != nil {
			log.Println("failed to parse selector")
			log.Println("selector: " + selector)
			log.Println("err: " + err.Error())
			return nil
		}
		data := expression.Get(obj)
		if IsNullOrEmpty(data) {
			return nil
		}

		if IsNullOrEmpty(data[0]) {
			return nil
		}
		marshal, _ := json.Marshal(data[0])
		return marshal
	}

}

func SafeGetToString(jsonString string, selector string) string {

	obj, err := oj.ParseString(jsonString)

	if err != nil {
		log.Println("failed to parse incoming jsonString for SafeGet operations")
		log.Println("jsonString: " + jsonString)
		log.Println("selector: " + selector)
		log.Println("err: " + err.Error())
		return ""
	} else {
		expression, err := jp.ParseString(selector)
		if err != nil {
			log.Println("failed to parse selector")
			log.Println("selector: " + selector)
			log.Println("err: " + err.Error())
			return ""
		}
		data := expression.Get(obj)
		if IsNullOrEmpty(data) {
			return ""
		}

		if IsNullOrEmpty(data[0]) {
			return ""
		}
		marshal, _ := json.Marshal(data[0])
		return string(marshal)
	}

}

func GetType(i interface{}) string {
	t := fmt.Sprintf("%T", i)
	return t
}
