package service

import (
	"fmt"
	"net/smtp"
	"time"

	"github.com/kingkhan77/log-sense/internal/models"
	"github.com/kingkhan77/log-sense/internal/repository"
	"github.com/kingkhan77/log-sense/pkg"
)

type NotificationService struct {
	alertRepo   *repository.AlertRepository
	oncallRepo  *repository.OnCallRepository
	userRepo    *repository.UserRepository
	serviceRepo *repository.ServiceRepository
	cfg         *pkg.Config
}

func NewNotificationService(
	alertRepo *repository.AlertRepository,
	oncallRepo *repository.OnCallRepository,
	userRepo *repository.UserRepository,
	serviceRepo *repository.ServiceRepository,
	cfg *pkg.Config,
) *NotificationService {
	return &NotificationService{
		alertRepo:   alertRepo,
		oncallRepo:  oncallRepo,
		userRepo:    userRepo,
		serviceRepo: serviceRepo,
		cfg:         cfg,
	}
}

func (s *NotificationService) Notify(alert models.Alert) error {
	if alert.NotificationSentAt != nil {
		return nil
	}

	fresh, err := s.alertRepo.FindByID(alert.ID)
	if err != nil {
		return err
	}
	if fresh.NotificationSentAt != nil {
		return nil
	}

	schedules, err := s.oncallRepo.GetAllCurrentOnCall(alert.ServiceID)
	if err != nil {
		return fmt.Errorf("failed to query on-call for service %s: %w", alert.ServiceID, err)
	}
	if len(schedules) == 0 {
		return fmt.Errorf("no on-call engineer for service %s", alert.ServiceID)
	}

	serviceName := alert.ServiceID
	if svc, err := s.serviceRepo.FindByID(alert.ServiceID); err == nil {
		serviceName = svc.Name
	}

	subject := fmt.Sprintf("[%s] %s", alert.Severity, alert.Title)
	body := fmt.Sprintf(
		"Alert: %s\nSeverity: %s\nService: %s\nDescription: %s\nTriggered: %s\nAlert ID: %s",
		alert.Title,
		alert.Severity,
		serviceName,
		alert.Description,
		alert.TriggeredAt.Format(time.RFC3339),
		alert.ID,
	)

	var firstErr error
	for _, schedule := range schedules {
		user, err := s.userRepo.FindByID(schedule.UserID)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if err := s.sendEmail(user.Email, subject, body); err != nil {
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	if firstErr != nil {
		return firstErr
	}

	now := time.Now()
	fresh.NotificationSentAt = &now
	fresh.UpdatedAt = now
	return s.alertRepo.Update(fresh)
}

func (s *NotificationService) sendEmail(to, subject, body string) error {
	cfg := s.cfg.SMTP
	if cfg.Host == "" {
		return fmt.Errorf("smtp not configured")
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		cfg.FromEmail, to, subject, body,
	))

	var auth smtp.Auth
	if cfg.Username != "" {
		auth = smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	}

	return smtp.SendMail(addr, auth, cfg.FromEmail, []string{to}, msg)
}
