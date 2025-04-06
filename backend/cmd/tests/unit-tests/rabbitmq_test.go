package rabbitmq

import (
	"log"
	"testing"

	"github.com/rabbitmq/amqp091-go"
)

// MockChannel is a mock implementation of the amqp091.Channel interface
type MockChannel struct {
	PublishFunc func(exchange, key string, mandatory, immediate bool, msg amqp091.Publishing) error
	ConsumeFunc func(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp091.Table) (<-chan amqp091.Delivery, error)
}

func (m *MockChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp091.Publishing) error {
	return m.PublishFunc(exchange, key, mandatory, immediate, msg)
}

func (m *MockChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp091.Table) (<-chan amqp091.Delivery, error) {
	return m.ConsumeFunc(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}

// TestPublish tests the publishing of messages
func TestPublish(t *testing.T) {
	mockChannel := &MockChannel{
		PublishFunc: func(exchange, key string, mandatory, immediate bool, msg amqp091.Publishing) error {
			if msg.Body == nil {
				return amqp091.ErrEmptyBody
			}
			return nil
		},
	}

	err := mockChannel.Publish("test_exchange", "test_key", false, false, amqp091.Publishing{Body: []byte("test message")})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = mockChannel.Publish("test_exchange", "test_key", false, false, amqp091.Publishing{Body: nil})
	if err == nil {
		t.Error("Expected error for empty message body, got none")
	}
}

// TestConsume tests the consumption of messages
func TestConsume(t *testing.T) {
	mockChannel := &MockChannel{
		ConsumeFunc: func(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp091.Table) (<-chan amqp091.Delivery, error) {
			deliveries := make(chan amqp091.Delivery)
			go func() {
				deliveries <- amqp091.Delivery{Body: []byte("test message")}
			}()
			return deliveries, nil
		},
	}

	deliveries, err := mockChannel.Consume("test_queue", "test_consumer", true, false, false, false, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	msg := <-deliveries
	if string(msg.Body) != "test message" {
		t.Errorf("Expected 'test message', got '%s'", msg.Body)
	}
}
