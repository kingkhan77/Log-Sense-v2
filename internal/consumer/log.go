package consumer

import (
	"bytes"
	"context"

	"github.com/IBM/sarama"
	"github.com/kingkhan77/log-sense/pkg"
	"github.com/opensearch-project/opensearch-go"
	opensearchapi "github.com/opensearch-project/opensearch-go/opensearchapi"
	"github.com/rs/zerolog/log"
)

type LogConsumer struct {
	os      *opensearch.Client
	topic   string
	groupID string
}

func (c *LogConsumer) Setup(
	sarama.ConsumerGroupSession,
) error {
	return nil
}

func (c *LogConsumer) Cleanup(
	sarama.ConsumerGroupSession,
) error {
	return nil
}

func NewLogConsumer(
	os *opensearch.Client,
	cfg *pkg.Config,
) *LogConsumer {

	return &LogConsumer{
		os:      os,
		topic:   cfg.Kafka.Topics.Logs,
		groupID: cfg.Kafka.ConsumerGroups.Logs,
	}
}

func (c *LogConsumer) Start(
	ctx context.Context,
	brokers []string,
) {

	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0

	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
	}
	group, err := sarama.NewConsumerGroup(
		brokers,
		c.groupID,
		config,
	)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to create consumer group")
	}

	defer group.Close()

	for {

		err := group.Consume(
			ctx,
			[]string{c.topic},
			c,
		)

		if err != nil {
			log.Error().Err(err).Msg("log consumer error")
		}

		if ctx.Err() != nil {
			return
		}
	}
}
func (c *LogConsumer) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {

	for msg := range claim.Messages() {

		req := opensearchapi.IndexRequest{
			Index: "logs",
			Body:  bytes.NewReader(msg.Value),
		}

		_, err := req.Do(
			session.Context(),
			c.os,
		)

		if err != nil {
			log.Error().
				Err(err).
				Msg("opensearch index failed")
			continue
		}

		session.MarkMessage(msg, "")
	}

	return nil
}
