package river

import (
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

type ArgumentResolversConfig struct {
	resolvers []ArgumentResolverFunc
}

type ResolverChain struct {
	Request   *Request
	Response  *Response
	Controller Controller
	Method reflect.Method
	ParamType reflect.Type
}

type ArgumentResolverFunc func(chain *ResolverChain) (reflect.Value, bool)

type Form struct {
	data url.Values
}

func (form *Form) String(name string) string  {
	return form.data.Get(name)
}
func (form *Form) Strings(name string) []string  {
	return form.data[name]
}
func (form *Form) Int(name string,defaultValues ...int) int  {
	value,err := strconv.Atoi(form.String(name))
	if err != nil {
		if len(defaultValues) == 0 {
			panic(DefaultError{
				400,
				name+" convert to string fail",
			})
		}else {
			return defaultValues[0]
		}
	}
	return value
}

func (arc *ArgumentResolversConfig) Add(resolverFunc ArgumentResolverFunc) *ArgumentResolversConfig {
	arc.resolvers = append(arc.resolvers, resolverFunc)
	return arc
}

var (
	responseWriterType = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
	sessionType = reflect.TypeOf((*Session)(nil)).Elem()
)

//requestResolverFunc
func requestResolverFunc(chain *ResolverChain) (reflect.Value, bool) {
	switch chain.ParamType {
		case reflect.TypeOf(chain.Request):
			return reflect.ValueOf(chain.Request), true
		case reflect.TypeOf(chain.Request).Elem():
			return reflect.ValueOf(chain.Request).Elem(), true
		case reflect.TypeOf(chain.Request.Request):
			return reflect.ValueOf(chain.Request.Request), true
		case reflect.TypeOf(chain.Request.Request).Elem():
			return reflect.ValueOf(chain.Request.Request).Elem(), true
		default:
			return reflect.ValueOf(nil), false
	}
}

//responseResolverFunc
func responseResolverFunc(chain *ResolverChain) (reflect.Value, bool) {
	switch chain.ParamType {
		case reflect.TypeOf(chain.Response):
			return reflect.ValueOf(chain.Response), true
		case reflect.TypeOf(chain.Response).Elem():
			return reflect.ValueOf(chain.Request).Elem(), true
		case responseWriterType:
			return reflect.ValueOf(chain.Response.ResponseWriter),true
		default:
			return reflect.ValueOf(nil), false
	}
}

func urlValuesResolverFunc(chain *ResolverChain) (reflect.Value, bool){
	if chain.ParamType == reflect.TypeOf(chain.Request.Form) {
		return reflect.ValueOf(chain.Request.Form), true
	}
	return reflect.ValueOf(nil), false
}

func requestBodyResolverFunc(chain *ResolverChain) (reflect.Value, bool)  {

	name := chain.ParamType.Name()
	isPtr := chain.ParamType.Kind() == reflect.Ptr
	if isPtr {
		name = chain.ParamType.Elem().Name()
	}
	if strings.HasSuffix(name,"JsonBody") {
		var value interface{}
		if isPtr {
			value = reflect.New(chain.ParamType.Elem()).Interface()
		}else{
			value = reflect.New(chain.ParamType).Interface()
		}
		jsonErr := chain.Request.BindJsonBody(value)
		if jsonErr != nil {
			panic(errors.New(jsonErr.Error()))
		}
		if isPtr {
			return reflect.ValueOf(value),true
		}
		return reflect.ValueOf(value).Elem(),true
	}
	return reflect.ValueOf(nil), false
}

func formResolverFunc(chain *ResolverChain) (reflect.Value, bool)  {
	name := chain.ParamType.Name()
	isPtr := chain.ParamType.Kind() == reflect.Ptr
	if isPtr {
		name = chain.ParamType.Elem().Name()
	}
	if strings.HasSuffix(name,"Form") {
		var v interface{}
		if !isPtr {
			v = reflect.New(chain.ParamType).Interface()
		}else {
			v = reflect.New(chain.ParamType.Elem()).Interface()
		}
		err := chain.Request.BindForm(v)
		if err != nil {
			panic(errors.New(err.Error()))
		}
		if isPtr {
			return reflect.ValueOf(v),true
		}
		return reflect.ValueOf(v).Elem(),true
	}
	return reflect.ValueOf(nil), false
}

func multipartFormResolverFunc(chain *ResolverChain) (reflect.Value, bool)  {
	switch chain.ParamType {
	case reflect.TypeOf(chain.Request.MultipartForm):
		err := chain.Request.ParseMultipartForm(config.MultipartMaxMemory)
		if err != nil {
			panic(errors.New(err.Error()))
		}
		return reflect.ValueOf(chain.Request.MultipartForm), true
	case reflect.TypeOf(chain.Request.MultipartForm).Elem():
		err := chain.Request.ParseMultipartForm(config.MultipartMaxMemory)
		if err != nil {
			panic(errors.New(err.Error()))
		}
		return reflect.ValueOf(chain.Request.MultipartForm).Elem(), true
	default:
		return reflect.ValueOf(nil), false
	}
}

func sessionResolverFunc(chain *ResolverChain) (reflect.Value, bool)  {
	if chain.ParamType == sessionType {
		return reflect.ValueOf(chain.Request.Session()),true
	}
	return reflect.ValueOf(nil), false
}

