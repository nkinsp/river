package river

import "io/ioutil"

type FormFile struct {
	Name string
	Path string
	Size int64
}

func (formFile *FormFile) GetBytes() ([]byte, error)  {

	return ioutil.ReadFile(formFile.Path)


}