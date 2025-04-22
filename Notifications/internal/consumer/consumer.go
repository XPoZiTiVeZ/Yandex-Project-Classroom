package consumer

import (
	"Classroom/Notifications/pkg/events"
	"Classroom/Notifications/pkg/logger"
	"context"
	"log"

	"github.com/IBM/sarama"
)

type Service interface{}

type consumer struct {
	master sarama.Consumer
}

func MustNew(brokers []string) *consumer {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("failed to create consumer: %v", err)
	}

	return &consumer{
		master: master,
	}
}

func (c *consumer) Close() error {
	return c.master.Close()
}

func (c *consumer) ConsumeTopic(ctx context.Context, topic string) {
	partitions, err := c.master.Partitions(topic)
	if err != nil {
		log.Fatalf("failed to get partitions: %v", err)
	}

	for _, partition := range partitions {
		pc, err := c.master.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("failed to subscribe partition %d: %v", partition, err)
		}

		go func(pc sarama.PartitionConsumer) {
			defer pc.Close()
			for {
				select {
				case msg := <-pc.Messages():
					handler, ok := handlers[msg.Topic]
					if !ok {
						logger.Error(ctx, "unknown topic", "topic", msg.Topic)
						continue
					}
					handler(ctx, msg)
				case err := <-pc.Errors():
					logger.Error(ctx, "failed to consume message", "err", err)
				case <-ctx.Done():
					return
				}
			}
		}(pc)
	}
}

type EventHandler func(ctx context.Context, msg *sarama.ConsumerMessage)

var handlers = map[string]EventHandler{
	events.CourseEnrolledTopic: handleUserEnrolled,
	events.CourseExpelledTopic: handleUserExpelled,
	events.LessonCreatedTopic:  handleLessonCreated,
	events.TaskCreatedTopic:    handleTaskCreated,
}

func handleUserEnrolled(ctx context.Context, msg *sarama.ConsumerMessage) {
	logger.Info(ctx, "Получено сообщение", "topic", msg.Topic, "partition", msg.Partition, "offset", msg.Offset, "key", string(msg.Key), "value", string(msg.Value))
}

func handleUserExpelled(ctx context.Context, msg *sarama.ConsumerMessage) {
	logger.Info(ctx, "Получено сообщение", "topic", msg.Topic, "partition", msg.Partition, "offset", msg.Offset, "key", string(msg.Key), "value", string(msg.Value))
}

func handleLessonCreated(ctx context.Context, msg *sarama.ConsumerMessage) {
	logger.Info(ctx, "Получено сообщение", "topic", msg.Topic, "partition", msg.Partition, "offset", msg.Offset, "key", string(msg.Key), "value", string(msg.Value))
}

func handleTaskCreated(ctx context.Context, msg *sarama.ConsumerMessage) {
	logger.Info(ctx, "Получено сообщение", "topic", msg.Topic, "partition", msg.Partition, "offset", msg.Offset, "key", string(msg.Key), "value", string(msg.Value))
}
