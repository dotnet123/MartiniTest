package main

import (

"fmt"
"reflect"
"net/http"
"strings"


"encoding/json"
	"github.com/dotnet123/fasthttptest/models"
)
type Server struct {
	allStruct map[string]*StructDefind
}

type StructDefind struct {
	name   string
	rvalue reflect.Value
	rtype  reflect.Type

	methods map[string]*Method
}

type Method struct {
	method     reflect.Method
	haveReturn bool
}

// construct method default
func NewServer() *Server {
	server := new(Server)
	server.allStruct = make(map[string]*StructDefind)

	return server
}

func (this *Server) Start(port string) error {
	return http.ListenAndServe(port, this)
}

func (this *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello bs!")
	for _, s_define := range this.allStruct {
		for m_name, m := range s_define.methods {
			methodAccept := "/" + s_define.name + "/" + m_name
			if strings.EqualFold(methodAccept, r.URL.Path) {
				//fmt.Fprintln(w, "true")
				m.method.Func.Call([]reflect.Value{s_define.rvalue, reflect.ValueOf(w), reflect.ValueOf(r)})
			}
		}
	}
	//fmt.Fprintf(w, "finish")
}

func (this *Server) Register(object interface{}) error {
	register := new(StructDefind)
	register.methods = make(map[string]*Method)
	register.rvalue = reflect.ValueOf(object)
	register.rtype = reflect.TypeOf(object)
	register.name = reflect.Indirect(register.rvalue).Type().Name()
	if register.name == "" {
		return fmt.Errorf("no service name for type ")
	}
	for m := 0; m < register.rtype.NumMethod(); m++ {
		method := register.rtype.Method(m)
		mtype := method.Type
		mname := method.Name

		if mtype.NumIn() != 3 {
			return fmt.Errorf("method %s has wrong number of ins: %d", mname, mtype.NumIn())
		}
		for in := 0; in < mtype.NumIn(); in++ {
			fmt.Println(mtype.In(in))
		}
		reply := mtype.In(1)
		if reply.String() != "http.ResponseWriter" {
			return fmt.Errorf("%s argument type not exported: %s", mname, reply)
		}
		arg := mtype.In(2)
		if arg.String() != "*http.Request" {
			return fmt.Errorf("%s argument type not exported: %s", mname, arg)
		}
		register.methods[mname] = &Method{method, false}

	}
	this.allStruct[register.name] = register
	return nil
}

func (this *Server) AllInfo() {
	for _, s_define := range this.allStruct {
		for m_name := range s_define.methods {
			fmt.Println("struct name:" + s_define.name + " method name:" + m_name)
		}
	}
}
func main(){

	//server := NewServer()
	//fmt.Println(server.Register(new(models.Assert)))
	//server.AllInfo()
	//server.Start(":9000")

	var data = `{"Name":"Xiao mi 6","Id":10}`

	objType := reflect.TypeOf(&models.User{}).Elem()
	obj := reflect.New(objType).Interface()

	err := json.Unmarshal([]byte(data), obj)

	var i models.IEvent
	i= new(models.UserHandler)
	a,b:= i.Create(obj)

	fmt.Print(a,b,err)
}
