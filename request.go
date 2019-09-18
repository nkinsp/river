package river

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Request struct {
	ResponseWriter http.ResponseWriter
	*http.Request
	attrMap map[string]interface{}
}


func (req *Request) Param(name string,defaultValue ...string) string{
	value := req.Form.Get(name)
	if value == "" && len(defaultValue) > 0{
		return defaultValue[0]
	}
	return value
}

func (req *Request) Params(name string) []string  {
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

func (req *Request) GetBody() ([]byte,error)  {
	return ioutil.ReadAll(req.Body)
}


func (req *Request) BindJsonBody(v interface{}) error  {
	data,err := req.GetBody()
	if err != nil {
		return err
	}
	return json.Unmarshal(data,v)
}

func (req *Request) GetAttr(name string) interface{}  {
	return req.attrMap[name]
}

func (req *Request) BindForm(v interface{}) error  {

	return convertFormTo(req.Form,v)

}

func (req *Request) SetAttr(name string,value interface{}) *Request {
	req.attrMap[name] = value
	return req
}

func (req *Request) SetAttrs(attrs map[string]interface{}) *Request{
	for key,value:=range attrs{
		req.SetAttr(key,value)
	}
	return req
}

func (req *Request) Session() Session  {
	if config.sessionManager == nil {
		panic(errors.New("SessionManager is nil"))
	}
	return config.sessionManager.Get(req)

}
