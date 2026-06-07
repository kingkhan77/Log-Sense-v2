package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/kingkhan77/log-sense/internal/constants"
	"github.com/kingkhan77/log-sense/internal/models"
	"github.com/kingkhan77/log-sense/internal/repository"
	"github.com/kingkhan77/log-sense/internal/service"
	"github.com/kingkhan77/log-sense/pkg"
	"github.com/google/uuid"
	"github.com/opensearch-project/opensearch-go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RuleEngine struct {
	ruleRepo      *repository.RuleRepository
	alertRepo     *repository.AlertRepository
	os            *opensearch.Client
	redis         *redis.Client
	alertProducer *service.AlertProducer
	interval      time.Duration
	workerCount   int
}

func NewRuleEngine(
	ruleRepo *repository.RuleRepository,
	alertRepo *repository.AlertRepository,
	os *opensearch.Client,
	redis *redis.Client,
	alertProducer *service.AlertProducer,
	cfg *pkg.Config,
) *RuleEngine {
	sec := cfg.Alerting.EvaluationIntervalSeconds
	if sec <= 0 {
		sec = 60
	}
	return &RuleEngine{
		ruleRepo:      ruleRepo,
		alertRepo:     alertRepo,
		os:            os,
		redis:         redis,
		alertProducer: alertProducer,
		interval:      time.Duration(sec) * time.Second,
		workerCount:   cfg.Alerting.WorkerCount,
	}
}

func (e *RuleEngine) Start(ctx context.Context) {
	ticker := time.NewTicker(e.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			e.evaluateRules()
		}
	}
}

func (e *RuleEngine) evaluateRules() {
	rules, err := e.ruleRepo.ListEnabled()
	if err != nil {
		log.Error().Err(err).Msg("failed to list enabled rules")
		return
	}

	jobs := make(chan models.AlertRule)
	var wg sync.WaitGroup

	for i := 0; i < e.workerCount; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for rule := range jobs {
				e.evaluateRule(rule)
			}
		}()
	}

	for _, rule := range rules {
		jobs <- rule
	}

	close(jobs)
	wg.Wait()

	e.resolveDisabledRuleAlerts(rules)
}

// resolveDisabledRuleAlerts finds open alerts whose rule is now disabled and resolves them.
func (e *RuleEngine) resolveDisabledRuleAlerts(enabledRules []models.AlertRule) {
	enabledIDs := make(map[string]struct{}, len(enabledRules))
	for _, r := range enabledRules {
		enabledIDs[r.ID] = struct{}{}
	}

	openAlerts, err := e.alertRepo.GetOpenAlerts()
	if err != nil {
		return
	}

	for _, alert := range openAlerts {
		if _, enabled := enabledIDs[alert.RuleID]; enabled {
			continue
		}
		rule, err := e.ruleRepo.FindByID(alert.RuleID)
		if err != nil {
			continue
		}
		if rule.IsEnabled {
			continue
		}
		// Rule is explicitly disabled — resolve the open alert.
		e.resolveAlertIfNeeded(*rule)
	}
}

func (e *RuleEngine) evaluateRule(rule models.AlertRule) {
	count, err := e.countMatchingLogs(rule)
	if err != nil {
		log.Error().Err(err).Str("rule_id", rule.ID).Msg("opensearch count failed")
		return
	}

	if count >= rule.Threshold {
		e.createAlertIfNeeded(rule, count)
	} else {
		e.resolveAlertIfNeeded(rule)
	}
}

func (e *RuleEngine) countMatchingLogs(
	rule models.AlertRule,
) (int, error) {

	query, err := buildCountQuery(rule)
	if err != nil {
		return 0, err
	}

	body, err := json.Marshal(query)
	if err != nil {
		return 0, err
	}

	res, err := e.os.Count(
		e.os.Count.WithIndex("logs"),
		e.os.Count.WithBody(bytes.NewReader(body)),
	)

	if err != nil {
		return 0, err
	}

	defer res.Body.Close()

	if res.IsError() {
		return 0, fmt.Errorf(
			"opensearch count failed: %s",
			res.Status(),
		)
	}

	var response struct {
		Count int `json:"count"`
	}

	if err := json.NewDecoder(
		res.Body,
	).Decode(&response); err != nil {
		return 0, err
	}

	return response.Count, nil
}

func (e *RuleEngine) createAlertIfNeeded(rule models.AlertRule, count int) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("alert:rule:%s", rule.ID)

	now := time.Now()
	alert := &models.Alert{
		ID:           uuid.NewString(),
		TenantID:     rule.TenantID,
		ServiceID:    rule.ServiceID,
		RuleID:       rule.ID,
		Title:        rule.Name,
		Description:  rule.Description,
		Severity:     rule.Severity,
		Status:       constants.AlertOpen,
		Threshold:    rule.Threshold,
		CurrentCount: count,
		TriggeredAt:  now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	ok, err := e.redis.SetNX(
		ctx,
		cacheKey,
		"pending",
		time.Duration(rule.WindowMinutes*2)*time.Minute,
	).Result()

	if !ok || err != nil {
		return
	}

	if err := e.alertProducer.Publish(alert); err != nil {
		log.Error().Err(err).Msg("failed to publish alert")
	}
}

func (e *RuleEngine) resolveAlertIfNeeded(rule models.AlertRule) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("alert:rule:%s", rule.ID)

	alert, err := e.alertRepo.GetOpenAlertByRule(rule.ID)
	if err != nil {
		return
	}

	now := time.Now()
	alert.Status = constants.AlertResolved
	alert.ResolvedAt = &now
	alert.UpdatedAt = now

	if err := e.alertRepo.Update(alert); err != nil {
		return
	}

	_ = e.redis.Del(ctx, cacheKey).Err()
}

func (e *RuleEngine) WarmAlertCache() {
	ctx := context.Background()
	alerts, err := e.alertRepo.GetOpenAlerts()
	if err != nil {
		return
	}

	for _, alert := range alerts {
		cacheKey := fmt.Sprintf("alert:rule:%s", alert.RuleID)
		_ = e.redis.Set(ctx, cacheKey, alert.ID, 0).Err()
	}
}
