package river

import (
	"log"
	"reflect"
	"strings"
)

//默认的错误处理
func defaultErrorHandler(req *Request,resp *Response ,err IError)  {

	log.Println("[River] ","defaultErrorHandler ",err.GetError())
	resp.Status(err.GetCode())
	resp.WriteString(err.GetError())
}

//静态文件处理
func staticFileHandler (req *Request,resp *Response){

	if !strings.HasSuffix(req.URL.Path,".ioc"){
		req.Next()
	}
}

//跨域处理
func crossHandler (req *Request,resp *Response){


	//fmt.Println("crossHandler::=>",req.handlerNext)


	req.Next()



}



func matchRouteHandler(req *Request,resp *Response)  {
	router := req.App.Router
	route,exists := router.find(req)
	if !exists {
		panic(Error404("No Find Mapping "+req.URL.Path))
		return
	}
	values := []reflect.Value{reflect.ValueOf(route.controller)}
	typeNumIn := route.hanlderMethod.Type.NumIn()
	for i :=0;i< typeNumIn;i++ {
		if i > 0{
			for _,resolver:= range req.App.ArgumentResolversConfig.resolvers{
				value,flag :=resolver(&ResolverChain{
					Request:req,
					Response:resp,
					ParamType:route.hanlderMethod.Type.In(i),
				})
				if flag {
					values = append(values,value)
				}
			}
		}
	}

	for _,interceptor := range  req.App.InterceptorRegister.interceptors{
		if interceptor.match(req.URL.Path) {
			if !interceptor.interceptor(req,resp,route.controller,route.hanlderMethod) {
				return
			}
		}
	}

	resultValues := route.hanlderMethod.Func.Call(values)
	if len(resultValues) == 0 {

	}






}