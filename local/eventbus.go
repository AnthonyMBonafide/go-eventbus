package local

import (
	"github.com/AnthonyMBonafide/go-eventbus"
	"sync"
)

// Local implementation of EventBus which uses no networking resources
type InMemoryEventBus struct {
	topicRegistry      map[string][]eventbus.MessageListener
	topicRegistryMutex *sync.Mutex
	localCache         map[string]interface{}
	localCacheMutex    *sync.Mutex
}

// New creates an instance of InMemoryEventBus
func New() eventbus.EventBus {
	return InMemoryEventBus{
		topicRegistry:      make(map[string][]eventbus.MessageListener),
		topicRegistryMutex: &sync.Mutex{},
		localCache:         make(map[string]interface{}),
		localCacheMutex:    &sync.Mutex{},
	}
}

// CreateConsumer creates and registers a consumer which reacts to messages published to the same instance of
// InMemoryEventBus
func (imeb InMemoryEventBus) CreateConsumer(topicID string, listener eventbus.MessageListener) {
	imeb.topicRegistryMutex.Lock()
	defer imeb.topicRegistryMutex.Unlock()

	n := append(imeb.topicRegistry[topicID], listener)
	imeb.topicRegistry[topicID] = n
}

// DeleteConsumer removes the specified consumer from recieving messages published to the same instance of
// InMemoryEventBus
func (imeb InMemoryEventBus) DeleteConsumer(topicID string, messageListenerID string) {
	imeb.topicRegistryMutex.Lock()
	defer imeb.topicRegistryMutex.Unlock()

	listeners, ok := imeb.topicRegistry[topicID]
	if ok {
		for index, listener := range listeners {
			if listener.ID == messageListenerID {
				// Slice trick which removes the element from the slice
				listeners = append(listeners[:index], listeners[index+1:]...)
				imeb.topicRegistry[topicID] = listeners
			}
		}
	}

}

// SendMessage published a message to all registered consumers listening to the same topicID on the same instance of
// InMemoryEventBus
func (imeb InMemoryEventBus) SendMessage(message eventbus.Message) {
	listeners, ok := imeb.topicRegistry[message.Topic]

	if !ok {
		return
	}

	sentMessages := 0
	for _, listener := range listeners {
		listener.Handler(message)
		sentMessages++
	}
}

// GetCacheValue gets the value from InMemoryEventBus' cache associated with the specified key
func (imeb InMemoryEventBus) GetCacheValue(key string) interface{} {
	imeb.localCacheMutex.Lock()
	defer imeb.localCacheMutex.Unlock()

	return imeb.localCache[key]
}

// SetCacheValue updates the InMemoryEventBus' cache with the specified key value pair
func (imeb InMemoryEventBus) SetCacheValue(key string, value interface{}) {
	imeb.localCacheMutex.Lock()
	defer imeb.localCacheMutex.Unlock()
	imeb.localCache[key] = value
}

// GetEventBusID gets the event bus ID
func (InMemoryEventBus) GetEventBusID() string {
	return "InMemoryEventBus"
}
