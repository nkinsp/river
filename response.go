package river

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	http.ResponseWriter
}


func (resp *Response) Status(code int) *Response{
	resp.WriteHeader(code)
	return resp
}



func (resp *Response) SetHeader(name string,value string) *Response  {
	resp.Header().Set(name,value)
	return resp
}

func (resp *Response) ContentType(contentType string) *Response{
	return resp.SetHeader("Content-Type",contentType)
}

func (resp *Response) Html(html string) *Response {

	resp.
		ContentType("text/html;charset=utf-8").
		WriteString(html)
	return resp

}

func (resp *Response) View(name string,data interface{})  {
	if config.viewEngine != nil {
		config.viewEngine.Render(resp.ResponseWriter,name,data)
	}
}

func (resp *Response) Redirect(url string) {
	resp.SetHeader("Location",url)
	resp.Status(301)

}

func (resp *Response) WriteString(data string)  {
	_, _ = resp.Write([]byte(data))
}

func (resp *Response) Cookie(cookie *http.Cookie) *Response {

	http.SetCookie(resp,cookie)
	return resp
}

func (resp *Response) Json(data interface{})   {

	jsonString,err :=json.Marshal(data)
	if err != nil {
		panic(err)
		return
	}
	resp.ContentType("application/json;charset=utf-8")
	_, _ = resp.Write(jsonString)


}
