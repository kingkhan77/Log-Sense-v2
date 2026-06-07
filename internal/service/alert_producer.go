package service

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/kingkhan77/log-sense/internal/models"
	"github.com/kingkhan77/log-sense/pkg"
)

type AlertProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewAlertProducer(
	producer sarama.SyncProducer,
	cfg *pkg.Config,
) *AlertProducer {
	return &AlertProducer{
		producer: producer,
		topic:    cfg.Kafka.Topics.Alerts,
	}
}

func (p *AlertProducer) Publish(alert *models.Alert) error {
	payload, err := json.Marshal(alert)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(payload),
	}

	_, _, err = p.producer.SendMessage(msg)
	return err
}
