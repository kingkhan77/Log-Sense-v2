package service

import (
	"time"

	"github.com/kingkhan77/log-sense/internal/repository"
)

type DashboardSummary struct {
	OpenAlerts      int64            `json:"open_alerts"`
	AlertsLast24h   int64            `json:"alerts_last_24h"`
	EnabledRules    int64            `json:"enabled_rules"`
	Services        int64            `json:"services"`
	BySeverity      map[string]int64 `json:"by_severity"`
}

type DashboardService struct {
	alertRepo   *repository.AlertRepository
	ruleRepo    *repository.RuleRepository
	serviceRepo *repository.ServiceRepository
}

func NewDashboardService(
	alertRepo *repository.AlertRepository,
	ruleRepo *repository.RuleRepository,
	serviceRepo *repository.ServiceRepository,
) *DashboardService {
	return &DashboardService{
		alertRepo:   alertRepo,
		ruleRepo:    ruleRepo,
		serviceRepo: serviceRepo,
	}
}

func (s *DashboardService) Summary(tenantID string) (*DashboardSummary, error) {
	open, err := s.alertRepo.CountOpenAlerts(tenantID)
	if err != nil {
		return nil, err
	}

	since := time.Now().Add(-24 * time.Hour)
	last24, err := s.alertRepo.CountSince(tenantID, since)
	if err != nil {
		return nil, err
	}

	openAlerts, err := s.alertRepo.GetOpenAlertsByTenant(tenantID)
	if err != nil {
		return nil, err
	}
	bySeverity := map[string]int64{}
	for _, a := range openAlerts {
		bySeverity[a.Severity]++
	}

	rules, err := s.ruleRepo.List(tenantID)
	if err != nil {
		return nil, err
	}
	var enabled int64
	for _, r := range rules {
		if r.IsEnabled {
			enabled++
		}
	}

	services, err := s.serviceRepo.List(tenantID)
	if err != nil {
		return nil, err
	}

	return &DashboardSummary{
		OpenAlerts:    open,
		AlertsLast24h: last24,
		EnabledRules:  enabled,
		Services:      int64(len(services)),
		BySeverity:    bySeverity,
	}, nil
}
