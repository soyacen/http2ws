package event

import "sync"

var events = make(map[string][]*emitter)
var lock sync.RWMutex

//var eventsCount int
//var maxListeners int

type emitter struct {
	C         chan interface{}
	EventName string
}

func On(eventName string) *emitter {
	lock.Lock()
	defer lock.Unlock()
	e := &emitter{
		C:         make(chan interface{}),
		EventName: eventName,
	}
	emitters, exists := events[eventName]
	if !exists {
		emitters = make([]*emitter, 0)
	}
	emitters = append(emitters, e)
	return e
}

func Emit(eventName string, msg interface{}) {
	lock.RLock()
	defer lock.Unlock()
	if emitters, exists := events[eventName]; exists {
		for _, e := range emitters {
			e.C <- msg
		}
	}
}

func (this *emitter) Destroy(eventName string) (err error) {
	if this == nil {
		return
	}
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	lock.Lock()
	defer lock.Unlock()
	close(this.C)
	emitters, exists := events[eventName]
	if !exists {
		return
	}
	index := -1
	for i, e := range emitters {
		if this == e {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}
	emitters = append(emitters[:index], emitters[index+1:]...)
	events[eventName] = emitters
	return
}
