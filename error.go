package river

type IError interface {
	error
	GetCode() int
	GetError() string
}

type DefaultError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err DefaultError) GetCode() int  {
	return err.Code
}

func (err DefaultError) GetError() string  {
	return err.Message
}

func (err DefaultError) Error() string  {
	return err.Message
}

func Error404(err string) IError  {

	return DefaultError{Code:404, Message:err}
}

func Error500(err string) IError  {
	return DefaultError{Code:500, Message:err}
}