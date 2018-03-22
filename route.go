package main

import (
	"reflect"

	"github.com/dotnet123/fasthttptest/inject"
)
type Martini struct {
	inject.Injector
	Router map[string]*Route
}
func New() *Martini {
	  m:=&Martini{inject.New(), make(map[string]*Route)}
	return m
}


type Route struct {
	inject.Injector
	parType reflect.Type
	Action  interface{}
	Handler interface{}

}



