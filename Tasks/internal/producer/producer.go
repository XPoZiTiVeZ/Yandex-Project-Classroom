package producer

import (
	"Classroom/Tasks/pkg/events"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type kafkaProducer struct {
	producer sarama.SyncProducer
}

func MustNewProducer(brokers []string) *kafkaProducer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("failed to create producer: %v", err)
	}

	return &kafkaProducer{producer: producer}
}

func (p *kafkaProducer) Close() error {
	return p.producer.Close()
}

func (p *kafkaProducer) PublishTaskCreated(msg events.TaskCreated) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	kafkaMsg := &sarama.ProducerMessage{
		Topic: events.TaskCreatedTopic,
		Value: sarama.ByteEncoder(data),
	}

	_, _, err = p.producer.SendMessage(kafkaMsg)
	return err
}
