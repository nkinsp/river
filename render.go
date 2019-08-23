package river

import (
	"reflect"
	"strings"
)

type RenderFunc func(req *Request,resp *Response,v reflect.Value) (support bool)




//
func ViewRenderFunc(req *Request,resp *Response,v reflect.Value) (support bool) {
	strValue := v.String()
	if strings.HasPrefix(strValue,"view:"){
		view := strings.Replace(strValue,"view:","",-1)
		engine := req.App.viewEngine
		if engine != nil {
			req.SetAttr("Form",req.Form)
			engine.Render(resp.ResponseWriter,view,req.attrMap)
			return true
		}
	}
	return false
}


func RedirectRenderFunc(req *Request,resp *Response,v reflect.Value) (support bool) {
	strValue := v.String()
	if strings.HasPrefix(strValue,"redirect:"){
		url := strings.Replace(strValue,"redirect:","",-1)
		resp.Redirect(url)
		return true
	}
	return false
}