package river

import (
	"io/ioutil"
	"net/http"
	"strconv"
)

type Request struct {
	*http.Request
	App *Application
	handlerNext bool
	attrMap map[string]interface{}


}

func (req *Request) Next()  {
	req.handlerNext = true
}

func (req *Request) Param(name string,defaultValue ...string) string{
	value := req.Form.Get(name)
	if value == "" && len(defaultValue) > 0{
		return defaultValue[0]
	}
	return value
}

func (req *Request) ParamArray(name string) []string  {
	return req.Form[name]
}

func (req *Request) ParamInt(name string,defaultValue ...int) (int,error)   {
	value,err := strconv.Atoi(req.Param(name))
	if err != nil && len(defaultValue) > 0{
		return defaultValue[0],err
	}
	return value,err
}


func (req *Request) ParamInt64(name string) (int64,error)   {
	return strconv.ParseInt(req.Param(name),10,64)
}

func (req *Request) ParamIntBool(name string) (bool,error)   {
	return strconv.ParseBool(name)
}

func (req *Request) ParamValues(name string) []string  {

	return req.Form[name]
}
func (req *Request) ParamStringValues(name string) []string  {

	//values:= req.ParamsMap[name];
	//stringValues

	return nil
}

func (req *Request) StringBody() string  {
	body,err :=ioutil.ReadAll(req.Body)
	if err != nil {
		return ""
	}
	return string(body)
}

func (req *Request) Json()  {

}

func (req *Request) GetAttr(name string) interface{}  {
	return req.attrMap[name]
}

func (req *Request) SetAttr(name string,value interface{}) *Request {
	req.attrMap[name] = value
	return req
}

func (rep *Request) SetAttrs(attrs map[string]interface{}) *Request{
	for key,value:=range attrs{
		rep.SetAttr(key,value)
	}
	return rep

}
