package river

import (
	"reflect"
	"regexp"
	"strings"
)

type InterceptorChain struct {
	Request    *Request
	Response   *Response
	IsFunc     bool
	Controller interface{}
	Method     reflect.Method
	Func       RouteFunc
}

type Interceptor interface {
	Pre(chain *InterceptorChain) bool
}

type InterceptorConfig struct {

	patterns        []string
	excludePatterns [] string
	interceptor     Interceptor
}

type InterceptorRegister struct {

	interceptors []*InterceptorConfig

}

func (config *InterceptorConfig) Patterns(patterns...string) *InterceptorConfig  {
	for _,pattern:= range  patterns{
		config.patterns = append(config.patterns,patternString(pattern))
	}
	return config
}

func patternString(pattern string) string  {
	return strings.ReplaceAll(strings.ReplaceAll(pattern,"/",""),"*",".*")
}

func (config *InterceptorConfig) ExcludePatterns(patterns...string) *InterceptorConfig  {
	for _,pattern:= range  patterns{
		config.excludePatterns = append(config.excludePatterns,patternString(pattern))
	}
	return config
}

func (config *InterceptorConfig) match(path string ) bool   {
	requestPath := strings.ReplaceAll(path,"/","")
	for _,pattern:= range config.patterns{
		match,err := regexp.Match(pattern,[]byte(requestPath))
		if err == nil && match {
			isMatch :=false
			for _,excludePattern := range config.excludePatterns{
				if exMatch,exErr:= regexp.Match(excludePattern,[]byte(requestPath)); exErr == nil && exMatch {
					isMatch = true
				}
			}
			if !isMatch{
				return  true
			}
		}
	}
	return  false
}


func (register *InterceptorRegister) Interceptor(interceptor Interceptor) *InterceptorConfig  {
	interceptorConfig := &InterceptorConfig{
		interceptor:interceptor,
		patterns:[]string{},
		excludePatterns: []string{},
	}
	register.interceptors = append(register.interceptors,interceptorConfig)
	return interceptorConfig
}

