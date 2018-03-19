package main

import (

"fmt"
"reflect"


"encoding/json"
"github.com/dotnet123/fasthttptest/models"
)



func main(){

	//server := NewServer()
	//fmt.Println(server.Register(new(models.Assert)))
	//server.AllInfo()
	//server.Start(":9000")
	m:=&models.UserHandler{&models.UserDal{}}
	r:=NewRoute()
	r.Register(m)
	var data = `{"Name":"Xiao mi 6","Id":10}`

	//objType := reflect.TypeOf(&models.User{}).Elem()
	//obj := reflect.New(objType).Interface()
	//err := json.Unmarshal([]byte(data), obj)
	routehandler,_:=r.allRouteHandlers["UserHandler/Create"]
	obj := reflect.New(routehandler.parType).Interface()
	err := json.Unmarshal([]byte(data), obj)

	//var i models.IEvent= new(models.UserHandler)
	a,b:= routehandler.Handler.Create(obj)

	fmt.Print(a,b,err)
}
