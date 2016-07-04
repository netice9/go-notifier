package notifier

import "sync"

// Notifier is a event broadcaster
type Notifier struct {
	sync.Mutex
	listeners        []chan interface{}
	lastNotification interface{}
}

// NewNotifier creates a new Notifier with initial notfication value
func NewNotifier(firstNotification interface{}) *Notifier {
	return &Notifier{
		lastNotification: firstNotification,
	}
}

func nonBlockingSendToChannel(chn chan interface{}, val interface{}) {
	// recover in the case of sending to closed channel
	defer func() {
		if r := recover(); r != nil {
			// ignore?
		}
	}()

	select {
	case chn <- val:
		// everything is ok
	default:
		// previous value is blocking the channel, remove it!
		select {
		case <-chn:
			// removed value, all clear to send!
			chn <- val
		default:
			// receiver read it, send it now!
			chn <- val
		}
	}

}

// Notify notifies current value to all listeners
func (n *Notifier) Notify(value interface{}) {
	n.Lock()
	defer n.Unlock()

	n.lastNotification = value
	for _, listener := range n.listeners {
		nonBlockingSendToChannel(listener, value)
	}

}

// AddListener creats a new listener channel
func (n *Notifier) AddListener(capacity int) chan interface{} {
	if capacity == 0 {
		capacity = 1
	}
	listenerChannel := make(chan interface{}, capacity)
	n.Lock()
	defer n.Unlock()
	n.listeners = append(n.listeners, listenerChannel)
	last := n.lastNotification
	listenerChannel <- last
	return listenerChannel
}

// RemoveListener removes and closes an existing listener channel
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
	close(listenerChannel)
}

// Close closes and removes all listeners
func (n *Notifier) Close() {
	n.Lock()
	defer n.Unlock()
	for _, listener := range n.listeners {
		close(listener)
	}
	n.listeners = []chan interface{}{}
}

// NumberOfListeners returns the current count of open listeners
func (n *Notifier) NumberOfListeners() int {
	n.Lock()
	defer n.Unlock()
	return len(n.listeners)
}
