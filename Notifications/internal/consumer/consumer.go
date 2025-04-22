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
	LessonCreated(ctx context.Context, lessonID, courseID string) error
	TaskCreated(ctx context.Context, taskID, courseID string) error
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

	consumer := &consumer{master: master, svc: svc}

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

	if err := c.svc.UserEnrolled(ctx, payload.UserID, payload.CourseID); err != nil {
		logger.Error(ctx, "failed to notify user enrolled", "user_id", payload.UserID, "err", err)
		return
	}

	logger.Debug(ctx, "notified user enrolled", "user_id", payload.UserID)
}

func (c *consumer) handleUserExpelled(ctx context.Context, msg *sarama.ConsumerMessage) {
	var payload events.UserExpelled
	if err := decodeMessage(msg, &payload); err != nil {
		logger.Error(ctx, "invalid user expelled payload")
		return
	}

	if err := c.svc.UserExpelled(ctx, payload.UserID, payload.CourseID); err != nil {
		logger.Error(ctx, "failed to notify user expelled", "user_id", payload.UserID, "err", err)
		return
	}

	logger.Debug(ctx, "notified user expelled", "user_id", payload.UserID)
}

func (c *consumer) handleLessonCreated(ctx context.Context, msg *sarama.ConsumerMessage) {
	var payload events.LessonCreated
	if err := decodeMessage(msg, &payload); err != nil {
		logger.Error(ctx, "invalid lesson created payload")
		return
	}

	if err := c.svc.LessonCreated(ctx, payload.LessonID, payload.CourseID); err != nil {
		logger.Error(ctx, "failed to notify lesson created", "lesson_id", payload.LessonID, "err", err)
		return
	}

	logger.Debug(ctx, "notified lesson created", "lesson_id", payload.LessonID)
}

func (c *consumer) handleTaskCreated(ctx context.Context, msg *sarama.ConsumerMessage) {
	var payload events.TaskCreated
	if err := decodeMessage(msg, &payload); err != nil {
		logger.Error(ctx, "invalid task created payload")
		return
	}

	if err := c.svc.TaskCreated(ctx, payload.TaskID, payload.CourseID); err != nil {
		logger.Error(ctx, "failed to notify task created", "task_id", payload.TaskID, "err", err)
		return
	}

	logger.Debug(ctx, "notified task created", "task_id", payload.TaskID)
}

func decodeMessage(msg *sarama.ConsumerMessage, dest any) error {
	return json.Unmarshal(msg.Value, dest)
}
