package main

import (
	"reflect"
	"github.com/dotnet123/fasthttptest/models"
	"fmt"
)

func NewRoute() *Route {
	if instance == nil {
		instance = &Route{ make(map[string]*RouteHandler)}
	}
	return instance
}


type Route struct {
	allRouteHandlers map[string]*RouteHandler
}
type RouteHandler struct {
	Url   string
	rvalue reflect.Value
	parType  reflect.Type
	Handler models.IEvent

}

var instance *Route



func (*Route) Register(event models.IEvent) error  {


	routeHandler := new(RouteHandler)

	rvalue := reflect.ValueOf(event).Elem()
	rtype := reflect.TypeOf(event)
	name := reflect.Indirect(rvalue).Type().Name()
	if name == "" {
		return fmt.Errorf("no service name for type ")
	}
	for m := 0; m < rtype.NumMethod(); m++ {
		method := rtype.Method(m)
		mtype := method.Type
		mname := method.Name

		NewRoute().allRouteHandlers[name+"/"+mname] = routeHandler

		//for in := 0; in < mtype.NumIn(); in++ {
		//	fmt.Println(mtype.In(in))
		//}
		n:=mtype.NumIn()
		p:=mtype.In(1).Kind()
		println(p,n)
		routeHandler.parType=reflect.TypeOf(p)
		routeHandler.Handler=event

	}

	return nil
}