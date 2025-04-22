package consumer

import (
	"Classroom/Notifications/pkg/logger"
	"context"
	"log"

	"github.com/IBM/sarama"
)

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
					logger.Info(ctx, "Получено сообщение", "topic", msg.Topic, "partition", msg.Partition, "offset", msg.Offset, "key", string(msg.Key), "value", string(msg.Value))
				case err := <-pc.Errors():
					logger.Error(ctx, "Получена ошибка", "err", err)
				case <-ctx.Done():
					logger.Info(ctx, "Получен сигнал завершения.")
					return
				}
			}
		}(pc)
	}
}
