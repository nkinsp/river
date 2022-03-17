package river

import (
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

/**
 * 参数解析
 */
type ArgumentResolversConfig struct {
	resolvers []ArgumentResolverFunc
}

type ResolverChain struct {
	Request    *Request
	Response   *Response
	Controller Controller
	Method     reflect.Method
	ParamType  reflect.Type
}

type ArgumentResolverFunc func(chain *ResolverChain) (reflect.Value, bool, error)

type Form struct {
	data url.Values
}

func (form *Form) String(name string) string {
	return form.data.Get(name)
}
func (form *Form) Strings(name string) []string {
	return form.data[name]
}
func (form *Form) Int(name string, defaultValues ...int) (int, error) {
	value, err := strconv.Atoi(form.String(name))

	if err != nil {
		return value, err
	}
	return value, nil
}

func (arc *ArgumentResolversConfig) Add(resolverFunc ArgumentResolverFunc) *ArgumentResolversConfig {
	arc.resolvers = append(arc.resolvers, resolverFunc)
	return arc
}

var (
	responseWriterType = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
	sessionType        = reflect.TypeOf((*Session)(nil)).Elem()
)

//requestResolverFunc
func requestResolverFunc(chain *ResolverChain) (reflect.Value, bool, error) {
	switch chain.ParamType {
	case reflect.TypeOf(chain.Request):
		return reflect.ValueOf(chain.Request), true, nil
	case reflect.TypeOf(chain.Request).Elem():
		return reflect.ValueOf(chain.Request).Elem(), true, nil
	case reflect.TypeOf(chain.Request.Request):
		return reflect.ValueOf(chain.Request.Request), true, nil
	case reflect.TypeOf(chain.Request.Request).Elem():
		return reflect.ValueOf(chain.Request.Request).Elem(), true, nil
	default:
		return reflect.ValueOf(nil), false, nil
	}
}

//responseResolverFunc
func responseResolverFunc(chain *ResolverChain) (reflect.Value, bool, error) {
	switch chain.ParamType {
	case reflect.TypeOf(chain.Response):
		return reflect.ValueOf(chain.Response), true, nil
	case reflect.TypeOf(chain.Response).Elem():
		return reflect.ValueOf(chain.Request).Elem(), true, nil
	case responseWriterType:
		return reflect.ValueOf(chain.Response.ResponseWriter), true, nil
	default:
		return reflect.ValueOf(nil), false, nil
	}
}

func urlValuesResolverFunc(chain *ResolverChain) (reflect.Value, bool, error) {
	if chain.ParamType == reflect.TypeOf(chain.Request.Form) {
		return reflect.ValueOf(chain.Request.Form), true, nil
	}
	return reflect.ValueOf(nil), false, nil
}

func requestBodyResolverFunc(chain *ResolverChain) (reflect.Value, bool, error) {

	name := chain.ParamType.Name()
	isPtr := chain.ParamType.Kind() == reflect.Ptr
	if isPtr {
		name = chain.ParamType.Elem().Name()
	}
	if strings.HasSuffix(name, "JsonBody") {
		var value interface{}
		if isPtr {
			value = reflect.New(chain.ParamType.Elem()).Interface()
		} else {
			value = reflect.New(chain.ParamType).Interface()
		}
		jsonErr := chain.Request.BindJsonBody(value)
		if jsonErr != nil {
			return reflect.ValueOf(nil), false, jsonErr
		}
		if isPtr {
			return reflect.ValueOf(value), true, nil
		}
		return reflect.ValueOf(value).Elem(), true, nil
	}
	return reflect.ValueOf(nil), false, nil
}

func formResolverFunc(chain *ResolverChain) (reflect.Value, bool, error) {
	name := chain.ParamType.Name()
	isPtr := chain.ParamType.Kind() == reflect.Ptr
	if isPtr {
		name = chain.ParamType.Elem().Name()
	}
	if strings.HasSuffix(name, "Form") {
		var v interface{}
		if !isPtr {
			v = reflect.New(chain.ParamType).Interface()
		} else {
			v = reflect.New(chain.ParamType.Elem()).Interface()
		}
		err := chain.Request.BindForm(v)
		if err != nil {
			return reflect.ValueOf(v), false, err
		}
		if isPtr {
			return reflect.ValueOf(v), true, nil
		}
		return reflect.ValueOf(v).Elem(), true, nil
	}
	return reflect.ValueOf(nil), false, nil
}

func multipartFormResolverFunc(chain *ResolverChain) (reflect.Value, bool, error) {
	switch chain.ParamType {
	case reflect.TypeOf(chain.Request.MultipartForm):
		err := chain.Request.ParseMultipartForm(config.MultipartMaxMemory)
		if err != nil {
			return reflect.ValueOf(nil), false, err
		}
		return reflect.ValueOf(chain.Request.MultipartForm), true, nil
	case reflect.TypeOf(chain.Request.MultipartForm).Elem():
		err := chain.Request.ParseMultipartForm(config.MultipartMaxMemory)
		if err != nil {
			return reflect.ValueOf(nil), false, err
		}
		return reflect.ValueOf(chain.Request.MultipartForm).Elem(), true, nil
	default:
		return reflect.ValueOf(nil), false, nil
	}
}

func sessionResolverFunc(chain *ResolverChain) (reflect.Value, bool, error) {
	if chain.ParamType == sessionType {
		return reflect.ValueOf(chain.Request.Session()), true, nil
	}
	return reflect.ValueOf(nil), false, nil
}
