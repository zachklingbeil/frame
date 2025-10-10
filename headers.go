package frame

import (
	"sync"
	"sync/atomic"
)

type Headers struct {
	Map      map[string]*atomic.Value
	watchers map[string][]chan any
	*sync.RWMutex
	*sync.Cond
}

func NewHeaders() *Headers {
	h := &Headers{
		Map:      make(map[string]*atomic.Value),
		watchers: make(map[string][]chan any),
		RWMutex:  &sync.RWMutex{},
	}
	h.Cond = sync.NewCond(h.RWMutex)
	return h
}

func (h *Headers) Add(key string, val int) {
	h.Lock()
	defer h.Unlock()
	if _, exists := h.Map[key]; !exists {
		h.Map[key] = &atomic.Value{}
	}
	h.Map[key].Store(val)
	for _, ch := range h.watchers[key] {
		select {
		case ch <- val:
		default:
		}
	}
	h.Broadcast()
}

func (h *Headers) Observe(key string) <-chan int {
	ch := make(chan int, 1)
	h.RLock()
	if v, exists := h.Map[key]; exists {
		if i, ok := v.Load().(int); ok {
			ch <- i
		}
	}
	h.RUnlock()
	return ch
}

func (h *Headers) Subtract(key string) {
	h.Lock()
	defer h.Unlock()
	delete(h.Map, key)
	for _, ch := range h.watchers[key] {
		select {
		case ch <- 0:
		default:
		}
	}
	h.Broadcast()

}

// WaitForCondition waits until the provided condition function returns true.
// It must be called with the lock held.
func (h *Headers) WaitForCondition(condFn func() bool) {
	h.Lock()
	defer h.Unlock()
	for !condFn() {
		h.Wait()
	}
}

// SignalOne signals one waiting goroutine.
func (h *Headers) SignalOne() {
	h.Lock()
	h.Signal()
	h.Unlock()
}

// BroadcastAll signals all waiting goroutines.
func (h *Headers) BroadcastAll() {
	h.Lock()
	h.Broadcast()
	h.Unlock()
}
