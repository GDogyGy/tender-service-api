package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type Producer struct {
	Producer sarama.SyncProducer
	Topic    string
}

func NewEventProducer(brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		Producer: producer,
		Topic:    topic,
	}, nil
}

func (p *Producer) SendEvent(eventType string, model interface{}) error {
	const op = "kafka.producer.SendEvent"
	eventData, err := json.Marshal(map[string]interface{}{
		"event_type": eventType,
		"data":       model,
		"timestamp":  time.Now().UTC(),
	})
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Partition: int32(2),
		Topic:     p.Topic,
		Value:     sarama.ByteEncoder(eventData),
	}

	_, _, err = p.Producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	return nil
}

func (p *Producer) SendAsync(eventType string, model interface{}) {
	go func() {
		if err := p.SendEvent(eventType, model); err != nil {
			log.Printf("Async Kafka send failed: %v", err)
		}
	}()
}

func (p *Producer) Close() error {
	return p.Producer.Close()
}
