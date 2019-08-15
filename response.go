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
		Status(200).
		ContentType("text/html;charset=utf-8").
		WriteString(html)
	return resp

}

func (resp *Response) Redirect(url string) {

	resp.SetHeader("Location",url)
}

func (resp *Response) WriteString(data string)  {

	resp.Write([]byte(data))
}

func (resp *Response) Cookie(cookie *http.Cookie) *Response {

	http.SetCookie(resp,cookie)
	return resp
}

func (resp *Response) Json(data interface{}) error  {

	json,err :=json.Marshal(data);
	if err != nil {
		return nil
	}
	resp.ContentType("application/json;charset=utf-8")
	resp.Write(json)
	return err


}
