package memory

import (
	"clap/staging/session"
	"container/list"
	"sync"
	"time"
)

var Pder = &Provider{list: list.New()}

type SessionStore struct {
	sid          string                      //session id唯一标示
	timeAccessed time.Time                   //最后访问时间
	value        map[interface{}]interface{} //session里面存储的值
}
type Provider struct {
	lock     sync.Mutex               //用来锁
	sessions map[string]*list.Element //用来存储在内存
	list     *list.List               //用来做gc
}

func (st *SessionStore) Set(key, value interface{}) error {
	st.value[key] = value
	Pder.SessionUpdate(st.sid)
	return nil
}

func (st *SessionStore) Get(key interface{}) interface{} {
	Pder.SessionUpdate(st.sid)
	if v, ok := st.value[key]; ok {
		return v
	} else {
		return nil
	}
	return nil
}

func (st *SessionStore) Delete(key interface{}) error {
	delete(st.value, key)
	Pder.SessionUpdate(st.sid)
	return nil
}

func (st *SessionStore) SessionID() string {
	return st.sid
}

func (Pder *Provider) SessionInit(sid string) (session.Session, error) {
	Pder.lock.Lock()
	defer Pder.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	newsess := &SessionStore{sid: sid, timeAccessed: time.Now(), value: v}
	element := Pder.list.PushBack(newsess)
	Pder.sessions[sid] = element
	return newsess, nil
}

func (Pder *Provider) SessionRead(sid string) (session.Session, error) {
	if element, ok := Pder.sessions[sid]; ok {
		return element.Value.(*SessionStore), nil
	} else {
		sess, err := Pder.SessionInit(sid)
		return sess, err
	}
	return nil, nil
}
func (Pder *Provider) SessionDestory(sid string) error {
	if element, ok := Pder.sessions[sid]; ok {
		delete(Pder.sessions, sid)
		Pder.list.Remove(element)
		return nil
	}
	return nil
}
func (Pder *Provider) SessionGC(maxlifetime int64) {
	Pder.lock.Lock()
	defer Pder.lock.Unlock()

	for {
		element := Pder.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*SessionStore).timeAccessed.Unix() + maxlifetime) < time.Now().Unix() {
			Pder.list.Remove(element)
			delete(Pder.sessions, element.Value.(*SessionStore).sid)
		} else {
			break
		}
	}
}

func (Pder *Provider) SessionUpdate(sid string) error {
	Pder.lock.Lock()
	defer Pder.lock.Unlock()
	if element, ok := Pder.sessions[sid]; ok {
		element.Value.(*SessionStore).timeAccessed = time.Now()
		Pder.list.MoveToFront(element)
		return nil
	}
	return nil
}

func init() {
	Pder.sessions = make(map[string]*list.Element, 0)
	session.Register("memory", Pder)
}
