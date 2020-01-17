package eventbus

// MessageListener an alias type for a function which accepts a string
type MessageListener struct {
	ID      string
	Handler func(message Message)
}

// EventBus the interface which exposes functionality available via EventBus implementations
type EventBus interface {
	// CreateConsumer registers the provided lister with the specified topicId.
	// The listener will be invoked with messages received on the topicId
	CreateConsumer(topicID string, listener MessageListener) error

	// DeleteConsumer removes a registered consumer and its associated MessageHandler from the topic.
	DeleteConsumer(topicID string, messageListenerID string) error

	// PublishMessage publishes the message to all consumers registered to the topicId
	PublishMessage(message Message) error

	// SendMessage delivers the message to one matching consumer.
	SendMessage(message Message) error

	// GetCacheValue gets the value associated with the specified key
	GetCacheValue(key string) (interface{}, error)

	// SetCacheValue sets the value associated with the specified key
	SetCacheValue(key string, value interface{}) error
}

// Message a struct which contains the components needed to send a message via the EventBus
type Message struct {
	// MessageID the unique identifier for the message
	MessageID string

	// Topic the targeted topic for the message
	Topic string

	// Payload the payload of the message
	Payload string
}
