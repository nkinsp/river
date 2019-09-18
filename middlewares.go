package river

import (
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)



//静态文件
func StaticFileMiddleware(dir string,patterns...string) MiddlewareFunc  {
	return func(request *Request,response *Response,next func()) {
		path := request.URL.Path
		for _,pattern:= range patterns{
			if match,err :=regexp.Match(pattern,[]byte(path));err == nil && match {
				filePath := strings.Join([]string{dir,path},"")
				data,err := ioutil.ReadFile(filePath)
				if err != nil {
					panic(Error404(strings.Join([]string{"Not Found ",path},"")))
					return
				}
				_, _ = response.Write(data)
				data = nil
				return

			}
		}
		next()

	}
}

type CorsConfig struct {
	patterns []string
	data     map[string]string
}

func (config *CorsConfig) AddMapping(pattern string) *CorsConfig {

	config.patterns = append(config.patterns,pattern)

	return config
}

func (config *CorsConfig) AllowCredentials(credentials bool) *CorsConfig {
	config.data["Access-Control-Allow-Credentials"] = strconv.FormatBool(credentials)
	return config
}

func (config *CorsConfig) AllowHeaders(headers ...string) *CorsConfig {
	config.data["Access-Control-Allow-Headers"] = strings.Join(headers,",")
	return config
}

func (config *CorsConfig) AllowMethods(methods ...string) *CorsConfig {
	config.data["Access-Control-Allow-Methods"] = strings.Join(methods,",")
	return config
}

func (config *CorsConfig) AllowOrigins(origins ...string) *CorsConfig {
	config.data["Access-Control-Allow-Origin"] = strings.Join(origins,",")
	return config
}

func (config *CorsConfig) ExposeHeaders(headers ...string) *CorsConfig {
	config.data["Access-Control-Expose-Headers"] = strings.Join(headers,",")
	return config
}


//跨域
func  CorsMiddleware(configFunc func(config *CorsConfig)) MiddlewareFunc {
	reg := &CorsConfig{patterns: []string{},data: map[string]string{}}
	configFunc(reg)
	return func(req *Request, resp *Response, next func()) {
		path := req.URL.Path
		for _,pattern:= range reg.patterns{
			match,err :=regexp.Match(pattern,[]byte(path))
			if err == nil && match {
				for key, value := range reg.data {
					resp.SetHeader(key,value)
				}
				break
			}
		}
		next()

	}
}

func NotFountMappingMiddleware() MiddlewareFunc  {
	return func(req *Request, resp *Response, next func()) {
		panic(Error404(strings.Join([]string{"Not Fount Mapping ",req.Method," ",req.URL.Path},"")))
	}
}


type ConfigRoueFunc func(router *Router)


func RouteMiddleware(prefix string,roueFunc ConfigRoueFunc) MiddlewareFunc  {
	router := &Router{
		routeMap:make(map[string]*RouteInfo),
		interceptorRegister:&InterceptorRegister{[]*InterceptorConfig{}},
	}
	router.Prefix(prefix)
	roueFunc(router)
	return func(req *Request, resp *Response, next func()) {
		if strings.HasPrefix(req.URL.Path,router.prefix) &&  matchRouteHandler(router,req,resp) {
			return
		}
		next()

	}
}