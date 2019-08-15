package river

import (
	"reflect"
)

type ArgumentResolversConfig struct {
	resolvers []ArgumentResolverFunc
}

type ResolverChain struct {
	Request *Request
	Response *Response
	ParamType reflect.Type
	Next func()
}

type ArgumentResolverFunc func(chain *ResolverChain) (reflect.Value,bool)

func (arc *ArgumentResolversConfig) Add(resolverFunc ArgumentResolverFunc) *ArgumentResolversConfig  {
	arc.resolvers = append(arc.resolvers,resolverFunc)
	return arc
}

//
func requestBodyArgumentResolverFunc(chain *ResolverChain) (reflect.Value,bool)   {

	if chain.ParamType == reflect.TypeOf(chain.Request){
		return reflect.ValueOf(chain.Request),true
	}
	if(chain.ParamType == reflect.TypeOf(chain.Request.Request)){

	}
	if chain.ParamType == reflect.TypeOf(chain.Response) {
		return reflect.ValueOf(chain.Response),true
	}
	return reflect.ValueOf(nil),false

}
