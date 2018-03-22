package util

import "reflect"

var Handlers []interface{}

func init() {
	Handlers = make([]interface{},0)
}

func InitHandler(i interface{}) error  {
	Handlers= append(Handlers,i)
	return nil
}
func T(i interface{}) interface{} {
	if reflect.TypeOf(i).Kind() == reflect.Ptr {
		return reflect.New(reflect.TypeOf(i).Elem()).Interface()
	}

	val := reflect.New(reflect.TypeOf(i))

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return val.Interface()
}