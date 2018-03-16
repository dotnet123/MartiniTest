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
."fasthttptest/models"
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
	user:= &User{}
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

  gorun(ctx)

}


func test_Rpc1(me *User) (int64, string, error) {
	return me.Id+78, me.Name, nil
}