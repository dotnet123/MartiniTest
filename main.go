package main

import (
	"github.com/go-martini/martini"
	"fmt"
	"reflect"
	"net/http"
	"strings"


	"martinitest/models"

	"encoding/json"
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
func main2()  {


	m := martini.Classic()
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.RunOnAddr(":8080")
}
package main

import (
	"github.com/valyala/fasthttp"
	"fmt"
	"flag"
	"log"
	"reflect"
	"sync"
	"errors"
	"encoding/json"

	"github.com/dotnet123/fasthttptest/models"
	"runtime"
)

var (
	addr     = flag.String("addr", ":8080", "TCP address to listen to")
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")
	locker sync.Mutex
)



var funcs   map[string]reflect.Value

func init() {
	funcs = make(map[string]reflect.Value)
}
func  Register(name string, f interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%s is not callable", name)
		}
	}()

	v := reflect.ValueOf(f)

	//to check f is function
	v.Type().NumIn()

	nOut := v.Type().NumOut()
	if nOut == 0 || v.Type().Out(nOut-1).Kind() != reflect.Interface {
		err = fmt.Errorf("%s return final output param must be error interface", name)
		return
	}

	_, b := v.Type().Out(nOut - 1).MethodByName("Error")
	if !b {
		err = fmt.Errorf("%s return final output param must be error interface", name)
		return
	}

	locker.Lock()
	if _, ok := funcs[name]; ok {
		err = fmt.Errorf("%s has registered", name)
		locker.Unlock()
		return
	}
	funcs[name] = v
	locker.Unlock()
	return
}
func main() {
	cpuNum:=runtime.NumCPU()
	runtime.GOMAXPROCS(cpuNum*2)
	flag.Parse()
	Register("rpc1", test_Rpc1)
	h := requestHandler
	if *compress {
		h = fasthttp.CompressHandler(h)
	}

	if err := fasthttp.ListenAndServe(*addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func gorun(ctx *fasthttp.RequestCtx)  {
	ctx.SetContentType("application/json; charset=utf-8")
	url:=string(ctx.RequestURI());log.Println(url)
	body:=ctx.Request.Body()
	user:= &models.User{}
	json.Unmarshal(body,user)

	f, ok := funcs[url[1:]]
	if !ok {
		//return nil, fmt.Errorf("rpc %s not registered", name)
		ret:=make(map[string]interface{})
		ret["error"]=1
		ret["data"]=nil
		b,_:=json.Marshal(ret)
		ctx.Response.SetBody(b)
		return
	}


	//in:= []reflect.Value{ reflect.ValueOf(1)}
	in:= []interface{ }{user}
	//inArgs := make([]interface{}, len(in))
	//for i := 0; i < len(in); i++ {
	//	inArgs[i] = in[i].Interface()
	//}
	inArgs:=in
	inValues := make([]reflect.Value, len(inArgs))

	for i := 0; i < len(inArgs); i++ {
		if inArgs[i] == nil {
			inValues[i] = reflect.Zero(f.Type().In(i))
		} else {
			inValues[i] = reflect.ValueOf(inArgs[i])
		}
	}

	out := f.Call(inValues)

	outArgs := make([]interface{}, len(out))
	for i := 0; i < len(outArgs); i++ {
		outArgs[i] = out[i].Interface()
	}

	p := out[len(out)-1].Interface()
	if p != nil {
		if e, ok := p.(error); ok {
			outArgs[len(out)-1] = errors.New(e.Error())
		} else {
			//return nil, fmt.Errorf("final param must be error")

		}
	}


	ret:=make(map[string]interface{})
	ret["error"]=0
	ret["data"]=outArgs
	b,_:=json.Marshal(ret)
	ctx.Response.SetBody(b)
	return
}
func requestHandler(ctx *fasthttp.RequestCtx) {

  //gorun(ctx)
  ctx.Response.SetBodyString("test")

}


func test_Rpc1(me *models.User) (int64, string, error) {
	return me.Id+78, me.Name, nil
}