package river

import (
	"log"
	"reflect"
	"runtime"
	"strings"
	"unicode"
)

type Router struct {
	prefix                   string
	routeMap 				 map[string]*RouteInfo
	interceptorRegister     *InterceptorRegister

}


type RouteInfo struct {
	path         string
	httpMethod   string
	isFunc       bool
	controller   interface{}
	handleMethod reflect.Method
	handleFunc   RouteFunc
}


type GroupFunc func(router *Router)

type RouteFunc func(req *Request,resp *Response)

func nameToPath(name string) string {
	path :=""
	if name != ""{
		for _,r:=range name{
			if unicode.IsUpper(r){
				path+="/"
			}
			path+=strings.ToLower(string(r))
		}
	}
	return path

}

func parseRoutePath(methodName string) (string,string) {
	path := ""
	requestMethod :=""
	methods :=[7]string{"Get","Post","Put","Delete","Options","Patch","Trace"}
	for _,method:=range methods  {
		if strings.HasPrefix(methodName,method){
			requestMethod = method
			rename := strings.Replace(methodName,method,"",-1)
			if rename != ""{
				if strings.HasPrefix(rename,"Path") && len(rename) > 4 {
					path += "/:"
					pathValue := strings.Replace(rename, "Path", "", -1)
					for j, v := range pathValue {
						if j == 0 {
							path += strings.ToLower(string(v))
						} else {
							path += string(v)
						}
					}
				}else{
					for _,r:=range rename {
						if unicode.IsUpper(r){
							path+="/"
						}
						path+=strings.ToLower(string(r))
					}
				}
			}
			if len(path) == 0{
				path+="/"
			}

		}
	}
	return path,requestMethod
}

func (router *Router) Get(path string,routeFunc RouteFunc) *Router  {
	return router.route(&RouteInfo{
		path:path,
		httpMethod:"GET",
		isFunc:true,
		handleFunc:routeFunc,
	})
}

func (router *Router) Post(path string,routeFunc RouteFunc) *Router  {
	return router.route(&RouteInfo{
		path:path,
		httpMethod:"POST",
		isFunc:true,
		handleFunc:routeFunc,
	})
}

func (router *Router) Put(path string,routeFunc RouteFunc) *Router  {
	return router.route(&RouteInfo{
		path:path,
		httpMethod:"PUT",
		isFunc:true,
		handleFunc:routeFunc,
	})
}

func (router *Router) Delete(path string,routeFunc RouteFunc) *Router  {
	return router.route(&RouteInfo{
		path:path,
		httpMethod:"DELETE",
		isFunc:true,
		handleFunc:routeFunc,
	})
}

func (router *Router) Options(path string,routeFunc RouteFunc) *Router  {
	return router.route(&RouteInfo{
		path:path,
		httpMethod:"OPTIONS",
		isFunc:true,
		handleFunc:routeFunc,
	})
}

func (router *Router) Patch(path string,routeFunc RouteFunc) *Router  {
	return router.route(&RouteInfo{
		path:path,
		httpMethod:"PATCH",
		isFunc:true,
		handleFunc:routeFunc,
	})
}

func (router *Router) Trace(path string,routeFunc RouteFunc) *Router  {
	return router.route(&RouteInfo{
		path:path,
		httpMethod:"TRACE",
		isFunc:true,
		handleFunc:routeFunc,
	})
}

func (router *Router) Any(path string,routeFunc RouteFunc) *Router  {
	router.Get(path,routeFunc)
	router.Post(path,routeFunc)
	router.Put(path,routeFunc)
	router.Delete(path,routeFunc)
	router.Options(path,routeFunc)
	router.Patch(path,routeFunc)
	router.Trace(path,routeFunc)
	return router
}


func (router *Router) Add(path string,controller interface{})  *Router{

	newPath := string(path)
	if len(newPath) > 1 && strings.HasSuffix(newPath,"/") {
		newPath = strings.TrimRight(newPath,"/")
	}
	rv :=reflect.TypeOf(controller)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return router
	}
	methodNum := rv.NumMethod()
	for i:=0; i < methodNum; i++ {
		method:=rv.Method(i)
		newSubPath := string(newPath)
		subPath,requestMethod := parseRoutePath(method.Name)
		if newSubPath != "" && requestMethod != "" {
			if subPath != "/" && newSubPath == "/" {
				newSubPath = string(subPath)
			} else if subPath != "/" && newSubPath != "/"{
				newSubPath+=subPath
			}
			router.route(&RouteInfo{
				path:newSubPath,
				httpMethod:strings.ToUpper(requestMethod),
				isFunc:false,
				controller:controller,
				handleMethod:method,
			})
		}
	}

	return router
}

func (router *Router) Interceptor(interceptor Interceptor) *InterceptorConfig {
	return router.interceptorRegister.Interceptor(interceptor)
}


func (router *Router) Group(path string,controllers ...Controller) *Router  {
	newPath := string(path)
	if len(newPath) > 1 && strings.HasSuffix(path,"/") {
		newPath = strings.TrimRight(path,"/")
	}
	for _,controller:=range controllers{
		subPath := string(newPath)
		rv :=reflect.TypeOf(controller)
		subPath += nameToPath(strings.Replace(rv.Elem().Name(),"Controller","",-1))
		router.Add(subPath,controller)
	}
	return  router;
}

func (router *Router) Prefix(prefix string) *Router {
	if prefix == ""{
		router.prefix = "/"
		return  router
	}
	if !strings.HasPrefix(prefix,"/") {
		prefix = strings.Join([]string{"/",prefix},"")
	}
	if strings.HasSuffix(prefix,"/") {
		prefix = strings.TrimRight(prefix,"/")
	}
	router.prefix = prefix
	return  router
}

func (router *Router) route(routeInfo *RouteInfo) *Router {
	if router.prefix != "/" {
		routeInfo.path = strings.Join([]string{router.prefix,routeInfo.path},"")
	}
	key:=strings.ToUpper(routeInfo.httpMethod)
	if len(routeInfo.path) > 1 && strings.HasSuffix(routeInfo.path,"/") {
		routeInfo.path = strings.TrimRight(routeInfo.path,"/")
	}
	key+= routeInfo.path
	pathKey := strings.Replace(key,"/","",-1)
	_,exists := router.routeMap[pathKey]
	if exists {
		log.Println("[River] ","Route Error",pathKey," It already exists.")
		runtime.Goexit()
	}

	router.routeMap[pathKey] = routeInfo
	logStr := strings.ToUpper(routeInfo.httpMethod)
	for i:=len(logStr);i<=8;i++{
		logStr +=" "
	}
	logStr += routeInfo.path
	logStr +="       "
	if !routeInfo.isFunc {
		logStr +=reflect.TypeOf(routeInfo.controller).String()
		logStr +="."+ routeInfo.handleMethod.Name+"()"
	}else{
		logStr += reflect.TypeOf(routeInfo.handleFunc).String()
	}
	log.Println("[River] ","Route",logStr)
	return  router
}

func (router *Router) find(req *Request) (*RouteInfo,bool) {
	key := strings.Join([]string{strings.ToUpper(req.Method),strings.Replace(req.URL.Path,"/","",-1)},"")
	route,exits:=router.routeMap[key]
	if !exits {
		for _,value:=range router.routeMap{
			if value.httpMethod == req.Method && strings.Index(value.path,":") != -1 {
				reqPathArr := strings.Split(req.URL.Path,"/")
				pathArr := strings.Split(value.path,"/")
				if len(reqPathArr) == len(pathArr) {
					isMatch := true
					for i,s:=range pathArr{
						if strings.HasPrefix(s,":")  {
							req.Form.Set(strings.Replace(s,":","",-1),reqPathArr[i])
						}else if  s != reqPathArr[i]{
							isMatch = false
							break
						}
					}
					if isMatch {
						return value,true
					}
				}
			}
		}
	}
	return route,exits

}





