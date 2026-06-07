package consumer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"github.com/kingkhan77/log-sense/internal/models"
	"github.com/kingkhan77/log-sense/internal/service"
	"github.com/kingkhan77/log-sense/pkg"
	"github.com/rs/zerolog/log"
)

type NotificationConsumer struct {
	notificationService *service.NotificationService
	topic               string
	groupID             string
}

func (c *NotificationConsumer) Setup(
	sarama.ConsumerGroupSession,
) error {
	return nil
}

func (c *NotificationConsumer) Cleanup(
	sarama.ConsumerGroupSession,
) error {
	return nil
}

func NewNotificationConsumer(
	notificationService *service.NotificationService,
	cfg *pkg.Config,
) *NotificationConsumer {
	return &NotificationConsumer{
		notificationService: notificationService,
		topic:               cfg.Kafka.Topics.Alerts,
		groupID:             cfg.Kafka.ConsumerGroups.Alerts,
	}
}

func (c *NotificationConsumer) Start(
	ctx context.Context,
	brokers []string,
) {
	config := sarama.NewConfig()

	config.Version = sarama.V2_8_0_0

	config.Consumer.Group.Rebalance.GroupStrategies =
    []sarama.BalanceStrategy{
        sarama.NewBalanceStrategyRoundRobin(),
    }

	group, err := sarama.NewConsumerGroup(
		brokers,
		c.groupID,
		config,
	)

	if err != nil {
		log.Fatal().
			Err(err).
			Msg("notification consumer group failed")
	}

	defer group.Close()

	for {
		err := group.Consume(
			ctx,
			[]string{c.topic},
			c,
		)

		if err != nil {
			log.Error().
				Err(err).
				Msg("notification consume failed")
		}

		if ctx.Err() != nil {
			return
		}
	}
}

func (c *NotificationConsumer) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {

	for msg := range claim.Messages() {

		var alert models.Alert

		if err := json.Unmarshal(
			msg.Value,
			&alert,
		); err != nil {

			log.Error().
				Err(err).
				Msg("invalid alert payload")

			continue
		}

		var notifyErr error

		for attempt := 1; attempt <= 3; attempt++ {

			notifyErr = c.notificationService.Notify(alert)

			if notifyErr == nil {
				break
			}

			log.Warn().
				Err(notifyErr).
				Int("attempt", attempt).
				Str("alert_id", alert.ID).
				Msg("notification retry")

			time.Sleep(
				time.Duration(attempt) * time.Second,
			)
		}

		if notifyErr != nil {
			log.Error().
				Err(notifyErr).
				Str("alert_id", alert.ID).
				Msg("notification permanently failed, skipping message")
		}

		// Always mark the message so the offset advances. Skipping it on
		// permanent failure is intentional: a stuck partition would block all
		// subsequent alerts. Operators should monitor the error logs above.
		session.MarkMessage(msg, "")

	}

	return nil
}
