package river

import (
	"fmt"
	"net/http"
	"sync"
	"github.com/satori/go.uuid"
)

type Session interface {

	Id() string
	Get(name string) (interface{},bool)
	Set(name string,value interface{})
	Remove(name ...string)

}

type SessionManager interface {


	Get(request *Request) Session



}

type ConfigSessionFunc func(request *Request)



type MemorySession struct {

	id string
	data map[string]interface{}

}

type MemorySessionManager struct {

	data map[string]interface{}
	lock sync.Mutex
}

func (manager *MemorySessionManager) Get(req *Request) Session  {

	fmt.Println("rhost::",req.Host)
	manager.lock.Lock()
	defer manager.lock.Unlock()
	config := req.App.Config
	cookie,err := req.Cookie(config.SessionName)

	if err != nil {
		cookie = &http.Cookie{
			Name:config.SessionName,
			Path:"/",
			Value:uuid.NewV4().String(),
			HttpOnly:true,

		}
		http.SetCookie(req.ResponseWriter,cookie)
	}
	session,exists := manager.data[cookie.Value]
	if !exists {
		session = &MemorySession{cookie.Value, map[string]interface{}{}}
		manager.data[cookie.Value] = session
	}
	return session.(Session)
}


func (session *MemorySession) Id() string  {

	return session.id
}

func (session *MemorySession) Get(name string) ( interface{},bool)  {
	value,exists := session.data[name]
	return value,exists
}

func (session *MemorySession) Set(name string,value interface{})  {

	session.data[name] = value

}

func (session *MemorySession) Remove(name ...string)  {

	for _,key:=range name{
		delete(session.data,key)
	}

}



