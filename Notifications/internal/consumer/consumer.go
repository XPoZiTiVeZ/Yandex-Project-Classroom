package consumer

import (
	"Classroom/Notifications/pkg/events"
	"Classroom/Notifications/pkg/logger"
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type NotificationsService interface {
	UserEnrolled(ctx context.Context, userID, courseID string) error
	UserExpelled(ctx context.Context, userID, courseID string) error
	LessonCreated(ctx context.Context, userID, courseID string) error
	TaskCreated(ctx context.Context, userID, courseID string) error
}

type EventHandler func(ctx context.Context, msg *sarama.ConsumerMessage)

type consumer struct {
	master   sarama.Consumer
	svc      NotificationsService
	handlers map[string]EventHandler
}

func MustNew(brokers []string, svc NotificationsService) *consumer {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("failed to create consumer: %v", err)
	}

	consumer := &consumer{master: master}

	consumer.handlers = map[string]EventHandler{
		events.CourseEnrolledTopic: consumer.handleUserEnrolled,
		events.CourseExpelledTopic: consumer.handleUserExpelled,
		events.LessonCreatedTopic:  consumer.handleLessonCreated,
		events.TaskCreatedTopic:    consumer.handleTaskCreated,
	}

	return consumer
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
					handler, ok := c.handlers[msg.Topic]
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

func (c *consumer) handleUserEnrolled(ctx context.Context, msg *sarama.ConsumerMessage) {
	var payload events.UserEnrolled
	if err := decodeMessage(msg, &payload); err != nil {
		logger.Error(ctx, "invalid user enrolled payload")
		return
	}
	logger.Info(ctx, "Получено сообщение", "topic", msg.Topic, "partition", msg.Partition, "offset", msg.Offset, "key", string(msg.Key), "value", string(msg.Value))
}

func (c *consumer) handleUserExpelled(ctx context.Context, msg *sarama.ConsumerMessage) {
	var payload events.UserExpelled
	if err := decodeMessage(msg, &payload); err != nil {
		logger.Error(ctx, "invalid user expelled payload")
		return
	}
	logger.Info(ctx, "Получено сообщение", "topic", msg.Topic, "partition", msg.Partition, "offset", msg.Offset, "key", string(msg.Key), "value", string(msg.Value))
}

func (c *consumer) handleLessonCreated(ctx context.Context, msg *sarama.ConsumerMessage) {
	var payload events.LessonCreated
	if err := decodeMessage(msg, &payload); err != nil {
		logger.Error(ctx, "invalid lesson created payload")
		return
	}
	logger.Info(ctx, "Получено сообщение", "topic", msg.Topic, "partition", msg.Partition, "offset", msg.Offset, "key", string(msg.Key), "value", string(msg.Value))
}

func (c *consumer) handleTaskCreated(ctx context.Context, msg *sarama.ConsumerMessage) {
	var payload events.TaskCreated
	if err := decodeMessage(msg, &payload); err != nil {
		logger.Error(ctx, "invalid task created payload")
		return
	}
	logger.Info(ctx, "Получено сообщение", "topic", msg.Topic, "partition", msg.Partition, "offset", msg.Offset, "key", string(msg.Key), "value", string(msg.Value))
}

func decodeMessage(msg *sarama.ConsumerMessage, dest any) error {
	return json.Unmarshal(msg.Value, dest)
}
