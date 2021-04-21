package topics

import (
	"atlas-wrg/rest/requests"
	"github.com/sirupsen/logrus"
	"sync"
)

type Registry struct {
	topics map[string]string
	lock   sync.RWMutex
}

var once sync.Once
var registry *Registry

func GetRegistry() *Registry {
	once.Do(func() {
		registry = &Registry{
			topics: make(map[string]string, 0),
			lock:   sync.RWMutex{},
		}
	})
	return registry
}

func (r *Registry) Get(l logrus.FieldLogger, token string) string {
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
		td, err := requests.Topic(l).GetTopic(token)
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
