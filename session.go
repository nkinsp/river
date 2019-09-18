package river

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
)

type Session interface {
	Id() string
	Get(name string) interface{}
	Set(name string, value interface{})
	Remove(name ...string)
}

type SessionConfig struct {
	Name       string
	Domain     string
	Path       string
	HttpOnly   bool
	ExpireTime int
}

type SessionManager interface {
	Get(request *Request) Session
}

type ConfigSessionFunc func(request *Request)

func randSessionId() string {

	b := make([]byte, 48)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		panic(err)
	}
	sum := md5.Sum(b)
	var data []byte
	for _, d := range sum {
		data = append(data, d)
	}
	return hex.EncodeToString(data)
}

func getSessionCookie(req *Request) *http.Cookie {

	cookie, err := req.Cookie(config.Session.Name)
	if err != nil || len(cookie.Value) != 32 {
		cookie = &http.Cookie{
			Name:     config.Session.Name,
			Path:     config.Session.Path,
			Value:    randSessionId(),
			Domain:   config.Session.Domain,
			HttpOnly: config.Session.HttpOnly,
		}

		http.SetCookie(req.ResponseWriter, cookie)
	}
	return cookie

}
