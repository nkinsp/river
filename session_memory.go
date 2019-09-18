package river

import (
	"sync"
	"time"
)

type MemorySession struct {

	id   	   string
	data       map[string]interface{}
	expireTime int64
	sync.Mutex


}


type MemorySessionManager struct {

	data map[string]Session
	lock sync.Mutex
}

func (manager *MemorySessionManager) remove(id string)  {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	delete(manager.data,id)
}

func (manager *MemorySessionManager) startExpireMonitor() * MemorySessionManager{
	go func() {
		for{
			time.Sleep(time.Millisecond*5000)
			timeNow := time.Now().Unix()
			for _,session := range  manager.data{
				m := session.(*MemorySession)
				if m.expireTime < timeNow {
					manager.remove(m.id)
				}
			}
		}
	}()
	return manager

}




func (manager *MemorySessionManager) Get(req *Request) Session  {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie := getSessionCookie(req)
	session,exists := manager.data[cookie.Value]
	if !exists {
		expireTime := time.Now().Unix()
		expireTime += int64(config.Session.ExpireTime)
		session = &MemorySession{id:cookie.Value,data:make(map[string]interface{}),expireTime:expireTime}
		manager.data[cookie.Value] = session

	}
	return session
}


func (session *MemorySession) Id() string  {

	return session.id
}

func (session *MemorySession) Get(name string) interface{}  {
	value,exists := session.data[name]
	if exists {
		return value
	}
	return nil
}

func (session *MemorySession) Set(name string,value interface{})  {

	session.Lock()
	defer  session.Unlock()
	session.data[name] = value

}

func (session *MemorySession) Remove(name ...string)  {
	session.Lock()
	defer  session.Unlock()
	for _,key:=range name{
		delete(session.data,key)
	}
}



func NewMemorySessionManager() *MemorySessionManager {
	manager := &MemorySessionManager{data:make(map[string]Session),}
	manager.startExpireMonitor()
	return manager
}