package river

import (
	"log"
	"net/http"
	"reflect"
	"runtime/debug"
)

type Config struct {
	Addr               string
	UploadMaxFileSize  int64
	MultipartMaxMemory int64
	Session            *SessionConfig
	middlebrows        []MiddlewareFunc
	renders            []RenderFunc
	errorHandlers      map[string]ErrorHandlerFunc
	resolvers          []ArgumentResolverFunc
	sessionManager     SessionManager
	viewConfig         *ViewConfig
	viewEngine         ViewEngine
	Debug              bool
}

type Controller interface{}

type ErrorHandlerFunc func(req *Request, resp *Response, err IError)

type MiddlewareFunc func(req *Request, resp *Response, next func()) error

var (
	config  *Config
	app     *Application
)

type Application struct {
	http.Handler
	Config *Config
}

func (app *Application) ViewEngine(engineFunc func(config *ViewConfig) ViewEngine) *Application {
	config.viewConfig = &ViewConfig{
		Dir:        "./",
		Prefix:     "views/",
		Suffix:     ".html",
		DelimLeft:  "{{",
		DelimRight: "}}",
		funcMap:    map[string]interface{}{},
		Cache:      false,
	}
	config.viewEngine = engineFunc(config.viewConfig)
	return app
}

func (app *Application) recoverError(request *Request, response *Response) {
	if pr := recover(); pr != nil {
		switch pr.(type) {
		case string:
			app.errorHandler(request, response, Error500(pr.(string)))
		default:
			prType := reflect.TypeOf(pr)
			if prType.Implements((reflect.TypeOf((*IError)(nil))).Elem()) {
				err := pr.(IError)
				app.errorHandler(request, response, err)
			} else if prType.Implements((reflect.TypeOf((*error)(nil))).Elem()) {
				err := pr.(error)
				app.errorHandler(request, response, Error500(err.Error()))
			}
		}
		if config.Debug {
			debug.PrintStack()
		}
	}
	request = nil
	response = nil
}

func (app *Application) executeMiddleware(request *Request, response *Response, index int) {
	if index < len(config.middlebrows) {
		middleware := config.middlebrows[index]
		middleware(request, response, func() {
			app.executeMiddleware(request, response, index+1)
		})
	}
}

func (app *Application) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

	request := &Request{
		ResponseWriter: resp,
		Request:        req,
		attrMap:        make(map[string]interface{}),
	}
	response := &Response{resp}
	defer app.recoverError(request, response)
	parseErr := request.ParseForm()
	if parseErr != nil {
		panic(parseErr)
		return
	}
	app.executeMiddleware(request, response, 0)

}

func App() *Application {
	if app != nil {
		return app
	}
	initConfig()
	initEnv()
	app = &Application{
		Config: config,
	}
	return app
}

func (app *Application) AddMiddleware(add func(func(middlewareFunc MiddlewareFunc))) *Application {

	add(func(middlewareFunc MiddlewareFunc) {
		app.Use(middlewareFunc)
	})

	return app
}

func (app *Application) Route(prefix string) (route *Router) {

	app.Use(RouteMiddleware(prefix,func(router *Router) {
		route = router
	}))
	return route
}

func (app *Application) Use(middlewareFunc MiddlewareFunc) *Application {
	config.middlebrows = append(config.middlebrows, middlewareFunc)
	return app
}

func (app *Application) errorHandler(req *Request, resp *Response, err IError) {
	for key, handler := range config.errorHandlers {
		if key == reflect.TypeOf(err).String() {
			handler(req, resp, err)
		}
	}
}

func (app *Application) SetDefaultErrorHandler(handlerFunc ErrorHandlerFunc) *Application {
	return app.SetErrorHandler("river.DefaultError", handlerFunc)
}

func (app *Application) SetErrorHandler(name string, handlerFunc ErrorHandlerFunc) *Application {
	config.errorHandlers[name] = handlerFunc
	return app
}

func (app *Application) AddRender(renderFunc RenderFunc) *Application {
	config.renders = append(config.renders, renderFunc)
	return app
}

func (app *Application) AddResolver(resolverFunc ArgumentResolverFunc) *Application {
	config.resolvers = append(config.resolvers, resolverFunc)
	return app
}

func (Application) SessionManager(manager SessionManager) {
	config.sessionManager = manager
}

func (app *Application) Run(addrs ...string) {
	var err error
	addr := Env.GetString("server.addr")
	if len(addrs) > 0 {
		addr = addrs[0]
	}
	defer func() {
		if err != nil {
			log.Println("[River] ", "Error  ", err.Error())
		}
	}()
	log.Println("[River] ", "Listening and serving HTTP on ", addr)
	err = http.ListenAndServe(addr, app)
}

func initConfig()  {
	if config !=  nil{
		return
	}
	config =  &Config{
		UploadMaxFileSize:  1024 * 1024 * 10,
		MultipartMaxMemory: 32 << 20,
		Session:            &SessionConfig{},
		resolvers: []ArgumentResolverFunc{
			requestResolverFunc,
			responseResolverFunc,
			sessionResolverFunc,
			urlValuesResolverFunc,
			formResolverFunc,
			requestBodyResolverFunc,
			multipartFormResolverFunc,
		},
		renders: []RenderFunc{
			ViewRenderFunc,
			RedirectRenderFunc,
		},
		errorHandlers: map[string]ErrorHandlerFunc{
			"river.DefaultError": defaultErrorHandler,
		},
	}
}



