package local

import (
	"math/rand"
	"sync"

	"github.com/AnthonyMBonafide/go-eventbus"
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
func (imeb InMemoryEventBus) CreateConsumer(topicID string, listener eventbus.MessageListener) error {
	imeb.topicRegistryMutex.Lock()
	defer imeb.topicRegistryMutex.Unlock()

	n := append(imeb.topicRegistry[topicID], listener)
	imeb.topicRegistry[topicID] = n
	return nil
}

// DeleteConsumer removes the specified consumer from receiving messages published to the same instance of
// InMemoryEventBus
func (imeb InMemoryEventBus) DeleteConsumer(topicID string, messageListenerID string) error {
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

	return nil
}

// PublishMessage published a message to all registered consumers listening to the same topicID on the same instance of
// InMemoryEventBus
func (imeb InMemoryEventBus) PublishMessage(message eventbus.Message) error {
	listeners, ok := imeb.topicRegistry[message.Topic]

	if !ok {
		return nil
	}

	sentMessages := 0
	for _, listener := range listeners {
		listener.Handler(message)
		sentMessages++
	}

	return nil
}

func (imeb InMemoryEventBus) SendMessage(message eventbus.Message) error {
	listeners, ok := imeb.topicRegistry[message.Topic]

	if !ok {
		return nil
	}

	randomIndex := rand.Intn(len(listeners) - 1)
	listeners[randomIndex].Handler(message)

	return nil
}

// GetCacheValue gets the value from InMemoryEventBus' cache associated with the specified key
func (imeb InMemoryEventBus) GetCacheValue(key string) (interface{}, error) {
	imeb.localCacheMutex.Lock()
	defer imeb.localCacheMutex.Unlock()

	return imeb.localCache[key], nil
}

// SetCacheValue updates the InMemoryEventBus' cache with the specified key value pair
func (imeb InMemoryEventBus) SetCacheValue(key string, value interface{}) error {
	imeb.localCacheMutex.Lock()
	defer imeb.localCacheMutex.Unlock()
	imeb.localCache[key] = value

	return nil
}

// GetEventBusID gets the event bus ID
func (InMemoryEventBus) GetEventBusID() string {
	return "InMemoryEventBus"
}
