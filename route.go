package river

import (
	"log"
	"reflect"
	"strings"
	"unicode"
)

type Router struct {

	routeMap map[string]RouteInfo
}



type RouteInfo struct {
	path string
	method string
	controller interface{}
	hanlderMethod reflect.Method
}

type GroupFunc func(router *Router)

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
			requestMethod = method;
			rname := strings.Replace(methodName,method,"",-1)
			if rname != ""{
				if strings.HasPrefix(rname,"Path") && len(rname) > 4 {
					path += "/:"
					pathValue := strings.Replace(rname, "Path", "", -1)
					for j, v := range pathValue {
						if j == 0 {
							path += strings.ToLower(string(v))
						} else {
							path += string(v)
						}
					}
				}else{
					for _,r:=range rname{
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

func (router *Router) Controllers(controllers ...Controller)  *Router {
	for _,controller:= range  controllers{
		router.Group("",controller)
	}
	return router
}

func (router *Router) Controller(path string,controller Controller)  *Router{

	newPath := string(path)
	if len(newPath) > 1 && strings.HasSuffix(newPath,"/") {
		newPath = strings.TrimRight(newPath,"/")
	}
	rv :=reflect.TypeOf(controller)
	methodNum := rv.NumMethod()
	for i:=0; i < methodNum; i++ {
		method:=rv.Method(i)
		newSubPath := string(newPath)
		subPath,requestMethod := parseRoutePath(method.Name)
		if newSubPath != "" && requestMethod != "" {
			if subPath != "/"{
				newSubPath+=subPath
			}
			router.route(newSubPath,requestMethod,controller,method)
		}
	}

	return router
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
		router.Controller(subPath,controller)
	}
	return  router;
}

func (router *Router) route(path string,method string,controller Controller,handler reflect.Method) *Router {
	key:=strings.ToUpper(method)
	if len(path) > 1 && strings.HasSuffix(path,"/") {
		path = strings.TrimRight(path,"/")
	}
	key+=path
	router.routeMap[strings.Replace(key,"/","",-1)] = RouteInfo{
		path:path,
		method:strings.ToUpper(method),
		controller:controller,
		hanlderMethod:handler,
	}
	logStr := strings.ToUpper(method)
	for i:=len(logStr);i<=8;i++{
		logStr +=" "
	}
	logStr +=path
	logStr +="       "
	logStr +=reflect.TypeOf(controller).String()
	logStr +="."+handler.Name+"()"
	log.Println("[River] ","Route",logStr)
	return  router
}

func (router *Router) find(req *Request) (RouteInfo,bool) {
	key := strings.ToUpper(req.Method)+strings.Replace(req.URL.Path,"/","",-1)
	route,exits:=router.routeMap[key]
	if !exits {
		for _,value:=range router.routeMap{
			if value.method == req.Method && strings.Index(value.path,":") != -1 {
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





