package river

import (
	"io/ioutil"
	"os"
	"regexp"
)

func getFileContent(fileName string) ([]byte,error)  {
	file,err := os.Open(fileName)
	if err != nil {
		return []byte{},err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

func StaticFileMiddleware(dir string,patterns...string) func(request *Request,response *Response,next func())  {

	return func(request *Request,response *Response,next func()) {
		path := request.URL.Path
		for _,pattern:= range patterns{
			if match,err :=regexp.Match(pattern,[]byte(path));err == nil && match {
				filePath := string(dir)+path
				data,err := getFileContent(filePath)
				if err != nil {
					panic(Error500(err.Error()))
					return
				}
				response.Status(200)
				response.Write(data)
				return

			}
		}
		next()

	}
}
