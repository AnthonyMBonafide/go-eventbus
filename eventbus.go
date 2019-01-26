package eventbus

//MessageListener an alias type for a function which accepts a string
type MessageListener struct {
	Id string
	Handler func(message Message)
}

//EventBus the interface which exposes functionality available via EventBus implementations
type EventBus interface {
	//New creates a new instance of an EventBus
	//New() EventBus

	//CreateConsumer registers the provided lister with the specified topicId.
	//The listener will be invoked with messages received on the topicId
	CreateConsumer(topicId string, listener MessageListener)

	//DeleteConsumer removes a registered consumer and its associated MessageHandler from the topic.
	DeleteConsumer(topicId string, messageListenerId string)

	//SendMessage publishes the message to all consumers registered to the topicId
	SendMessage(message Message)

	//GetCacheValue gets the value associated with the specified key
	GetCacheValue(key string) interface{}

	//SetCacheValue sets the value associated with the specified key
	SetCacheValue(key string, value interface{})

	//GetEventBusId returns the unique identifier for the eventbus
	GetEventBusId() string
}

//Message a struct which contains the components needed to send a message via the EventBus
type Message struct {
	//MessageId the unique identifier for the message
	MessageId string

	//Topic the targeted topic for the message
	Topic     string

	//Payload the payload of the message
	Payload   string
}