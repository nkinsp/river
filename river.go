package river

import (
	"log"
	"net/http"
	"reflect"
	"runtime/debug"
)

const  (
	DEBUG_MODE = true
)


type Config struct {
	UploadMaxFileSize int64
}

type Controller interface {}

type HandlerFunc func(req *Request,resp *Response)

type ErrorHandlerFunc func(req *Request,resp *Response ,err IError)

type ResultRender func(req *Request,resp *Response,result interface{})


type HandlerInfo struct {
	Name string
	Func HandlerFunc
}

type Application struct {
	http.Handler
	Config *Config
	Router *Router
	ArgumentResolversConfig *ArgumentResolversConfig
	InterceptorRegister *InterceptorRegister
	HandlerMap map[string]HandlerFunc
	Handlers []string
	ErrorHandlers map[string]ErrorHandlerFunc
}
type configRoueFunc func(router *Router)

type configArgumentResolversFunc func(config *ArgumentResolversConfig)

type configInterceptorsFunc func(reigter *InterceptorRegister)

// config routes
func (app *Application) ConfigRoute(fun configRoueFunc) *Application {
	fun(app.Router)
	return app
}

func (app *Application) ConfigView()  {

}

func (app *Application) ConfigArgumentResolvers(fun configArgumentResolversFunc) *Application {
	fun(app.ArgumentResolversConfig)
	return app
}

func (app *Application) ConfigInterceptors(fun configInterceptorsFunc) *Application {
	fun(app.InterceptorRegister)
	return app
}


func (app *Application) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	request :=&Request{req,app,true,make(map[string]interface{})}
	response :=&Response{resp}
	defer func() {
		if pr := recover(); pr != nil {
			switch pr.(type) {
				case string:
					app.ErrorHandler(request,response,Error500(pr.(string)))
				default:
					 prType := reflect.TypeOf(pr)
					 if prType.Implements((reflect.TypeOf((*error)(nil))).Elem()) {
						err :=pr.(error)
						app.ErrorHandler(request,response,Error500(err.Error()))
					 }
			}
			if DEBUG_MODE {
				debug.PrintStack()
			}
		}
	}()
	request.ParseForm()
	handlerIndex :=0
	for request.handlerNext && handlerIndex < len(app.Handlers) {
		request.handlerNext = false
		handler :=app.HandlerMap[app.Handlers[handlerIndex]]
		handler(request,response)
		handlerIndex++
	}
}

func App() *Application {
	app := &Application{
		Config:&Config{UploadMaxFileSize:1024*1024*10},
		Router:&Router{map[string]RouteInfo{}},
		ArgumentResolversConfig:&ArgumentResolversConfig{[]ArgumentResolverFunc{requestBodyArgumentResolverFunc}},
		InterceptorRegister:&InterceptorRegister{[]*InterceptorConfig{}},
		HandlerMap: map[string]HandlerFunc{},
		Handlers:[]string{},
		ErrorHandlers: map[string]ErrorHandlerFunc{"river.DefaultError":defaultErrorHandler},
	}
	app.SetHandler("staticFileHandler",staticFileHandler)
	app.SetHandler("crossHandler",crossHandler)
	app.SetHandler("matchRouteHandler",matchRouteHandler)
	return app
}

func (app *Application) SetHandler(name string,handlerFunc HandlerFunc) *Application {
	app.HandlerMap[name] = handlerFunc
	app.Handlers = append(app.Handlers,name)
	log.Println("[River] ","Handler",name)
	return app
}

func (app *Application) ErrorHandler(req *Request,resp *Response, err IError)  {
	for key,handler:= range app.ErrorHandlers{
		if key == reflect.TypeOf(err).String() {
			handler(req,resp,err)
		}
	}
}

func (app *Application) SetErrorHandler(name string,handlerFunc ErrorHandlerFunc) *Application  {
	app.ErrorHandlers[name] = handlerFunc
	return  app
}


func (app *Application) Run(addr string)  {
	log.Println("[River] ","Listening and serving HTTP on ",addr )
	http.ListenAndServe(addr,app)
}

