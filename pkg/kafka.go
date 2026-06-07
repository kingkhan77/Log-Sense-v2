package pkg

import (
	"github.com/IBM/sarama"
)

func NewKafkaProducer(
	cfg *Config,
) (sarama.SyncProducer, error) {

	config := sarama.NewConfig()

	config.Producer.Return.Successes = true

	return sarama.NewSyncProducer(
		cfg.Kafka.Brokers,
		config,
	)
}