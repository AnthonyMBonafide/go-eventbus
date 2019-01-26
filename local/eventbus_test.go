package local

import (
	"fmt"
	"github.com/AnthonyMBonafide/go-eventbus"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	eb := New()
	eb.SetCacheValue("integer", 1)
	result := eb.GetCacheValue("integer")
	if 1 != result {
		t.Fail()
	}
}

func TestCacheConcurrency(t *testing.T) {
	eb := New()

	for i := 0; i <= 200; i++ {
		go eb.SetCacheValue("integer", 1)
	}
}

func TestMessageHandler(t *testing.T) {
	eb := New()
	done := make(chan struct{})
	eb.CreateConsumer("TestTopic", eventbus.MessageListener{
		ID: "test-listener",
		Handler: func(message eventbus.Message) {
			log.Printf("Handler received message '%s' with payload '%s'", message.MessageID, message.Payload)
			close(done)
		},
	})

	eb.SendMessage(eventbus.Message{
		MessageID: "TestMessage",
		Topic:     "TestTopic",
		Payload:   "TestPayload",
	})

	select {
	case <-done:
		log.Println("MessageListener 'done' signal received")
	case <-time.After(time.Second):
		{
			t.Log("The message was not received by the MessageListener in time")
			t.Fail()
		}
	}

}

func TestRemoveConsumer(t *testing.T) {
	eb := New()
	done := make(chan struct{})
	eb.CreateConsumer("TestTopic", eventbus.MessageListener{
		ID: "test-listener",
		Handler: func(message eventbus.Message) {
			log.Printf("Handler received message '%s' with payload '%s'", message.MessageID, message.Payload)
			close(done)
		},
	})

	eb.DeleteConsumer("TestTopic", "test-listener")

	eb.SendMessage(eventbus.Message{
		MessageID: "TestMessage",
		Topic:     "TestTopic",
		Payload:   "TestPayload",
	})

	select {
	case <-done:
		{
			t.Log("Unexpected 'done' signal received from MessageListener")
			t.Fail()
		}

	case <-time.After(time.Second):
		log.Println("The message was not received by the deleted MessageListener")
	}
}

func BenchmarkInMemoryEventBus_CreateConsumer(b *testing.B) {
	b.ReportAllocs()
	eb := New()
	for i := 0; i < b.N; i++ {
		eb.CreateConsumer("TestingTopic", eventbus.MessageListener{
			ID:      fmt.Sprintf("TestConsumer#%d", i),
			Handler: func(message eventbus.Message) {},
		})
	}
}

func BenchmarkInMemoryEventBus_DeleteConsumer(b *testing.B) {
	eb := New()
	for i := 0; i < b.N; i++ {
		eb.CreateConsumer("TestingTopic", eventbus.MessageListener{
			ID:      fmt.Sprintf("TestConsumer#%d", i),
			Handler: func(message eventbus.Message) {},
		})
	}

	b.ResetTimer()
	b.ReportAllocs()
	// Remove all consumers starting with the last to simulate the worst case scenario where we are forced to iterate
	// All the consumers in the registry
	for i := b.N; i >= 0; i-- {
		eb.DeleteConsumer("TestingTopic", fmt.Sprintf("TestConsumer#%d", i))
	}
}

func BenchmarkInMemoryEventBus_SendMessage(b *testing.B) {
	eb := New()
	for i := 0; i < b.N; i++ {
		eb.CreateConsumer("TestingTopic", eventbus.MessageListener{
			ID:      fmt.Sprintf("TestConsumer#%d", i),
			Handler: func(message eventbus.Message) {},
		})
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		eb.SendMessage(eventbus.Message{
			Topic:     fmt.Sprintf("TestConsumer#%d", i),
			MessageID: fmt.Sprintf("Message#%d", i),
			Payload:   "Message",
		})
	}
}

func BenchmarkInMemoryEventBus_SetCacheValue(b *testing.B) {
	eb := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		eb.SetCacheValue(fmt.Sprintf("CacheKey#%d", i), fmt.Sprintf("CacheValue#%d", i))
	}
}

func BenchmarkInMemoryEventBus_GetCacheValue(b *testing.B) {
	eb := New()
	for i := 0; i < b.N; i++ {
		eb.SetCacheValue(fmt.Sprintf("CacheKey#%d", i), fmt.Sprintf("CacheValue#%d", i))
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		eb.GetCacheValue(fmt.Sprintf("CacheKey#%d", i))
	}
}

func BenchmarkInMemoryEventBus_SetCacheValueLock(b *testing.B) {
	eb := New()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			eb.SetCacheValue(fmt.Sprintf("CacheKey#%d", rand.Int()), fmt.Sprintf("CacheValue#%d", rand.Int()))
		}
	})
}

func BenchmarkInMemoryEventBus_GetCacheValueLock(b *testing.B) {
	eb := New()
	for i := 0; i < b.N; i++ {
		eb.SetCacheValue(fmt.Sprintf("CacheKey#%d", i), fmt.Sprintf("CacheValue#%d", i))
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			eb.GetCacheValue(fmt.Sprintf("CacheKey#%d", rand.Intn(b.N)))
		}
	})
}
