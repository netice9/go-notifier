package notifier

import "sync"

type Notifier struct {
	sync.Mutex
	listeners []chan interface{}
	lastValue interface{}
}

func NewNotifier() *Notifier {
	return &Notifier{}
}

func (n *Notifier) Notify(value interface{}) {
	n.Lock()
	defer n.Unlock()

	n.lastValue = value
	for _, listener := range n.listeners {

		go func() {
			listener <- value
		}()
	}

}

func (n *Notifier) AddListener(listenerChannel chan interface{}) {
	n.Lock()
	defer n.Unlock()
	for _, existing := range n.listeners {
		if existing == listenerChannel {
			return
		}
	}
	n.listeners = append(n.listeners, listenerChannel)
}

func (n *Notifier) RemoveListener(listenerChannel chan interface{}) {
	n.Lock()
	defer n.Unlock()
	filtered := []chan interface{}{}
	for _, existing := range n.listeners {
		if existing != listenerChannel {
			filtered = append(filtered, existing)
		}
	}

	n.listeners = filtered

}
func (n *Notifier) NumberOfListeners() int {
	n.Lock()
	defer n.Unlock()
	return len(n.listeners)
}
