package river

import (
	"log"
	"reflect"
)

//默认的错误处理
func defaultErrorHandler(req *Request, resp *Response, err IError) {

	log.Println("[River] ", "defaultErrorHandler ", err.GetError())
	resp.Status(err.GetCode())
	resp.WriteString(err.GetError())

}

//参数解析
func resolverHandler(req *Request, resp *Response, route *RouteInfo) []reflect.Value {

	values := []reflect.Value{reflect.ValueOf(route.controller)}
	typeNumIn := route.handleMethod.Type.NumIn()
	for i := 0; i < typeNumIn; i++ {
		if i > 0 {
			chain := &ResolverChain{
				Request:    req,
				Response:   resp,
				Controller: route.controller,
				Method:     route.handleMethod,
				ParamType:  route.handleMethod.Type.In(i),
			}
			isSupport := false
			for _, resolver := range config.resolvers {
				value, flag := resolver(chain)
				if flag {
					values = append(values, value)
					isSupport = true
					break
				}
			}
			chain = nil
			if !isSupport {
				panic(DefaultError{
					400,
					reflect.TypeOf(route.controller).String() + "." + chain.Method.Name + "(...) Param [" + chain.ParamType.String() + "] No matching type was found",
				})
			}
		}
	}
	return values
}

//拦截器
func interceptorHandler(router *Router, chain *InterceptorChain) bool {
	for _, config := range router.interceptorRegister.interceptors {
		if config.match(chain.Request.URL.Path) {
			if !config.interceptor.Pre(chain) {
				return true
			}
		}
	}
	chain = nil
	return false
}

func valueConverter(req *Request, resp *Response, v reflect.Value) {

	switch v.Kind() {
	case reflect.Interface:
		valueConverter(req, resp, v.Elem())
	default:
		for _, render := range config.renders {
			if render(req, resp, v) {
				return
			}
		}
		//执行默认
		switch v.Kind() {
		case reflect.String:
			resp.WriteString(v.String())
		default:
			resp.Json(v.Interface())
		}

	}

}

//结果集映射
func resultValueHandler(req *Request, resp *Response, resultValues []reflect.Value) {

	if len(resultValues) == 0 {
		return
	}

	if len(resultValues) > 1 {
		var data []interface{}
		for _, value := range resultValues {
			data = append(data, value.Interface())
		}
		valueConverter(req, resp, reflect.ValueOf(data))
		return
	}
	valueConverter(req, resp, resultValues[0])
}

//匹配路由
func matchRouteHandler(router *Router, req *Request, resp *Response) (match bool) {
	route, exists := router.find(req)
	if !exists {
		return false
	}
	//拦截器
	if interceptorHandler(router, &InterceptorChain{
		Request:    req,
		Response:   resp,
		IsFunc:     route.isFunc,
		Controller: route.controller,
		Method:     route.handleMethod,
		Func:       route.handleFunc,
	}) {
		return true
	}
	if route.isFunc {
		route.handleFunc(req,resp)
		return true
	}
	//参数映射
	values := resolverHandler(req, resp, route)
	//执行
	resultValues := route.handleMethod.Func.Call(values)
	//结果映射
	resultValueHandler(req, resp, resultValues)
	return true
}
