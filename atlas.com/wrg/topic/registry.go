package topic

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type registry struct {
	topics map[string]string
	lock   sync.RWMutex
}

var once sync.Once
var r *registry

func GetRegistry() *registry {
	once.Do(func() {
		r = &registry{
			topics: make(map[string]string, 0),
			lock:   sync.RWMutex{},
		}
	})
	return r
}

func (r *registry) Get(l logrus.FieldLogger, token string) string {
	r.lock.RLock()
	if val, ok := r.topics[token]; ok {
		r.lock.RUnlock()
		return val
	} else {
		r.lock.RUnlock()
		r.lock.Lock()
		if val, ok = r.topics[token]; ok {
			r.lock.Unlock()
			return val
		}
		td, err := GetTopic(l)(token)
		if err != nil {
			r.lock.Unlock()
			l.WithError(err).Fatalf("Unable to locate topic for token %s.", token)
			return ""
		}
		r.topics[token] = td.Attributes.Name
		r.lock.Unlock()
		return td.Attributes.Name
	}
}
