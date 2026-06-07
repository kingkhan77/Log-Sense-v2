package service

import (
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"github.com/kingkhan77/log-sense/pkg"
)

type IngestionService struct {
	producer sarama.SyncProducer
	topic    string
}

func NewIngestionService(
	producer sarama.SyncProducer,
	cfg *pkg.Config,
) *IngestionService {
	return &IngestionService{
		producer: producer,
		topic:    cfg.Kafka.Topics.Logs,
	}
}

type LogIngestInput struct {
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Timestamp string                 `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata"`
}

func (s *IngestionService) PublishLog(
	tenantID, serviceID string,
	input LogIngestInput,
) error {
	ts := input.Timestamp
	if ts == "" {
		ts = time.Now().UTC().Format(time.RFC3339)
	}

	doc := map[string]interface{}{
		"tenant_id":  tenantID,
		"service_id": serviceID,
		"level":      input.Level,
		"message":    input.Message,
		"timestamp":  ts,
		"metadata":   input.Metadata,
	}
	if doc["metadata"] == nil {
		doc["metadata"] = map[string]interface{}{}
	}

	payload, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: s.topic,
		Value: sarama.ByteEncoder(payload),
	}

	_, _, err = s.producer.SendMessage(msg)
	return err
}
