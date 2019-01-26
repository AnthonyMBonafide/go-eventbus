package local

import (
	"github.com/AnthonyMBonafide/go-eventbus"
	"sync"
)

type InMemoryEventBus struct {
	topicRegistry      map[string][]eventbus.MessageListener
	topicRegistryMutex *sync.Mutex
	localCache         map[string]interface{}
	localCacheMutex    *sync.Mutex
}

func New() eventbus.EventBus {
	//log.Printf("Creating new InMemoryEventBus")
	return InMemoryEventBus{
		topicRegistry:      make(map[string][]eventbus.MessageListener),
		topicRegistryMutex: &sync.Mutex{},
		localCache:         make(map[string]interface{}),
		localCacheMutex:    &sync.Mutex{},
	}
}

func (imeb InMemoryEventBus) CreateConsumer(topicId string, listener eventbus.MessageListener) {
	imeb.topicRegistryMutex.Lock()
	defer imeb.topicRegistryMutex.Unlock()

	n := append(imeb.topicRegistry[topicId], listener)
	imeb.topicRegistry[topicId] = n
}

func (imeb InMemoryEventBus) DeleteConsumer(topicId string, messageListenerId string) {
	imeb.topicRegistryMutex.Lock()
	defer imeb.topicRegistryMutex.Unlock()

	listeners, ok := imeb.topicRegistry[topicId]
	if ok {
		for index, listener := range listeners {
			if listener.ID == messageListenerId {
				// Slice trick which removes the element from the slice
				listeners = append(listeners[:index], listeners[index+1:]...)
				imeb.topicRegistry[topicId] = listeners
			}
		}
	}

}

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

func (imeb InMemoryEventBus) GetCacheValue(key string) interface{} {
	imeb.localCacheMutex.Lock()
	defer imeb.localCacheMutex.Unlock()

	return imeb.localCache[key]
}

func (imeb InMemoryEventBus) SetCacheValue(key string, value interface{}) {
	imeb.localCacheMutex.Lock()
	defer imeb.localCacheMutex.Unlock()
	imeb.localCache[key] = value
}

func (InMemoryEventBus) GetEventBusId() string {
	return "InMemoryEventBus"
}
