package main

import (
	"reflect"
	"encoding/json"

	"github.com/dotnet123/fasthttptest/util"
	"strings"
	"github.com/labstack/gommon/log"
   _ "github.com/dotnet123/fasthttptest/models"
	"fmt"
	"github.com/dotnet123/fasthttptest/models"

	"flag"
	"github.com/valyala/fasthttp"
	"runtime"

	"github.com/dotnet123/fasthttptest/inject"
)
func init() {
	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu) // 尝试使用所有可用的CPU
}
var (
	addr     = flag.String("addr", ":8080", "TCP address to listen to")
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")
)

var m *Martini

func main() {

	m = New()
	for _, h := range util.Handlers {

		route := &Route{}

		handlerType := reflect.TypeOf(h)
		if handlerType.Kind() == reflect.Ptr {

			handlerType = handlerType.Elem()
		}
		handlerName := handlerType.Name()
		if ! strings.HasSuffix(handlerName, "Handler") {
			log.Fatalf("%s 命名规则不正确", handlerName)
		}

		handler := reflect.New(reflect.TypeOf(h).Elem()).Elem().Interface()
		route.Handler=handler
		//fmt.Printf("1---->%s \r\n",reflect.TypeOf(obj).Elem().String())

		type1 := reflect.TypeOf(handler)
		type2 := reflect.TypeOf(models.UserHandler{})
		fmt.Println(" 一 ", type1.String())
		fmt.Println(" 二 ", type2.String())

		//route.Map(obj)

		//route.SetParent(m.Injector)

		for i := 0; i < handlerType.NumMethod(); i++ {
			method := handlerType.Method(i)
			mtype := method.Type
			mname := method.Name
			pattern :="/"+ strings.ToLower(handlerName[0:len(handlerName)-7] + "/" + mname)
			for i := 0; i < mtype.NumIn(); i++ {
				parType := mtype.In(i)
				if parType.Kind() == reflect.Ptr {

					parType = parType.Elem()
				}
				parStructName := parType.Name()

				if strings.Contains(parStructName, "Dal") {
					newDalStruct := reflect.New(parType).Elem().Interface()
					type1 = reflect.TypeOf(newDalStruct)
					fmt.Printf(" 三 %s", type1.String())
					m.Map(newDalStruct)
				} else {
					route.parType = parType
				}
			}

			route.Action = method.Func.Interface()
			m.Router[pattern] = route

		}

	}

	 for uri,_:= range m.Router{

	 	println("\r",uri)
	 }

	flag.Parse()
	h := requestHandler
	if *compress {
		h = fasthttp.CompressHandler(h)
	}

	if err := fasthttp.ListenAndServe(*addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}


}

func requestHandler(ctx *fasthttp.RequestCtx) {
	path:=strings.TrimSpace(string(ctx.Path()))
	//fmt.Println("\r",path)
	msg := models.ApiMsg{}
	if r,ok:=m.Router[path]; ok {
	    r.Injector= inject.New()
		r.SetParent(m.Injector)
		obj := reflect.New(r.parType).Interface()
		err := json.Unmarshal(ctx.Request.Body(), obj)

		r.Map(r.Handler)
		r.Map(obj)


		values, err := r.Invoke(r.Action)
		if err != nil{
			msg.Error=1
			msg.Msg=err.Error()
			goto return0
		}
		var url1= strings.SplitAfterN(path,"/",3)
		//println(url1[2])
		switch url1[2] {
		case "select":
			list0, count1, err2 := values[0], values[1], values[2]
			log.Debug(list0)
			if err2.Interface() == nil {
				msg.Data = obj
				msg.Count = count1.Int()
			} else {
				if e, ok := err2.Interface().(error); ok {
					msg.Msg = e.Error()
					msg.Error = 1
				}
			}
		case "create":
			int0,  err1 := values[0], values[1]

			if err1.Interface() == nil {
				msg.Data = obj
				msg.Count = int0.Int()
			} else {
				if e, ok := err1.Interface().(error); ok {
					msg.Msg = e.Error()
					msg.Error = 1
				}
			}
		default:
			msg.Msg = "url 不匹配"
			msg.Error = 1
		}

	}
return0:
	b, _ := json.Marshal(msg)
	ctx.Response.SetBody(b)
	ctx.SetContentType("application/json; charset=utf-8")

	return
}


//if v.Kind() == reflect.Ptr {
// // println(i,"->", reflect.Indirect(values[i]).String())
//}
//objType := reflect.TypeOf(&models.User{}).Elem()
//obj := reflect.New(parType.Elem()).Interface()
//err := json.Unmarshal([]byte(data), obj)
//inj1.Map(obj)
//user,ok := obj.(*models.User)

