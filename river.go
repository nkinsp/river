package river

import (
	"log"
	"net/http"
	"reflect"
	"runtime/debug"
)


type Config struct {
	UploadMaxFileSize  int64
	MultipartMaxMemory int64
	SessionName        string
	SessionExpireTime  int
}

type Controller interface{}

type ErrorHandlerFunc func(req *Request, resp *Response, err IError)

type MiddlewareFunc func(req *Request, resp *Response, next func())

type Application struct {
	http.Handler
	Config                  *Config
	Router                  *Router
	Debug					bool
	ArgumentResolversConfig *ArgumentResolversConfig
	InterceptorRegister     *InterceptorRegister
	SessionManager          SessionManager
	errorHandlers           map[string]ErrorHandlerFunc
	middlewares             []MiddlewareFunc
	viewConfig              *ViewConfig
	viewEngine              ViewEngine
	renders                 []RenderFunc

}
type configRoueFunc func(router *Router)

type configArgumentResolversFunc func(config *ArgumentResolversConfig)

type configInterceptorsFunc func(reigter *InterceptorRegister)

// config routes
func (app *Application) ConfigRoute(fun configRoueFunc) *Application {
	fun(app.Router)
	return app
}

func (app *Application) ViewEngine(engineFunc func(config *ViewConfig) ViewEngine) *Application {
	app.viewConfig = &ViewConfig{
		Dir:        "./",
		Prefix:     "views/",
		Suffix:     ".html",
		DelimLeft:  "{{",
		DelimRight: "}}",
		funcMap:    map[string]interface{}{},
		Cache:      false,

	}
	app.viewEngine = engineFunc(app.viewConfig)
	return app
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
	request := &Request{
		ResponseWriter: resp,
		Request:        req,
		attrMap:        map[string]interface{}{},
		App:            app,
		handlerNext:    true,
	}
	response := &Response{resp}
	defer func() {
		if pr := recover(); pr != nil {
			switch pr.(type) {
			case string:
				app.ErrorHandler(request, response, Error500(pr.(string)))
			default:
				prType := reflect.TypeOf(pr)
				if prType.Implements((reflect.TypeOf((*IError)(nil))).Elem()) {
					err := pr.(IError)
					app.ErrorHandler(request, response, err)
				} else if prType.Implements((reflect.TypeOf((*error)(nil))).Elem()) {
					err := pr.(error)
					app.ErrorHandler(request, response, Error500(err.Error()))
				}
			}
			if app.Debug {
				debug.PrintStack()
			}
		}
	}()
	parseErr := request.ParseForm()
	if parseErr != nil{
		panic(parseErr)
		return
	}
	//中间件
	nextIndex := 0
	hasNext := true
	for hasNext && nextIndex < len(app.middlewares) {
		hasNext = false
		middleware := app.middlewares[nextIndex]
		middleware(request, response, func() {
			hasNext = true
			nextIndex++
		})
	}
	//路由匹配
	if hasNext {
		matchRouteHandler(request, response)
	}

}

func App() *Application {
	app := &Application{
		Config: &Config{
			UploadMaxFileSize:  1024 * 1024 * 10,
			MultipartMaxMemory: 32 << 20,
			SessionName:        "session",
			SessionExpireTime:  1800,
		},
		Router: &Router{map[string]RouteInfo{}},
		Debug:true,
		ArgumentResolversConfig: &ArgumentResolversConfig{[]ArgumentResolverFunc{
			requestResolverFunc,
			responseResolverFunc,
			urlValuesResolverFunc,
			formResolverFunc,
			requestBodyResolverFunc,
			multipartFormResolverFunc,
		}},
		InterceptorRegister: &InterceptorRegister{[]*InterceptorConfig{}},
		errorHandlers:       map[string]ErrorHandlerFunc{"river.DefaultError": defaultErrorHandler},
		SessionManager:      &MemorySessionManager{data: map[string]interface{}{}},
		middlewares:         []MiddlewareFunc{},
		renders:			 []RenderFunc{
			ViewRenderFunc,
			RedirectRenderFunc,
		},
	}
	return app
}

func (app *Application) Route(path string, controller Controller) *Application {
	app.Router.Add(path, controller)
	return app

}

func (app *Application) RouteGroup(path string, controllers ...Controller) *Application {
	app.Router.Group(path, controllers)
	return app

}

func (app *Application) Use(middlewareFunc MiddlewareFunc) *Application {
	app.middlewares = append(app.middlewares, middlewareFunc)
	return app
}

func (app *Application) ErrorHandler(req *Request, resp *Response, err IError) {
	for key, handler := range app.errorHandlers {
		if key == reflect.TypeOf(err).String() {
			handler(req, resp, err)
		}
	}
}

func (app *Application) SetErrorHandler(name string, handlerFunc ErrorHandlerFunc) *Application {
	app.errorHandlers[name] = handlerFunc
	return app
}

func (app *Application) Run(addr string) {
	log.Println("[River] ", "Listening and serving HTTP on ", addr)
	err := http.ListenAndServe(addr,app)
	if err != nil {
		log.Println("[River] Error ", err)
	}
}
