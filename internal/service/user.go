package service

import (
	"errors"
	"time"

	"github.com/kingkhan77/log-sense/internal/models"
	"github.com/kingkhan77/log-sense/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateDeveloper(tenantID, name, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now()
	user := &models.User{
		ID:           uuid.NewString(),
		TenantID:     tenantID,
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		Role:         "DEVELOPER",
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	return s.userRepo.Create(user)
}

func (s *UserService) ListDevelopers(tenantID string) ([]models.User, error) {
	return s.userRepo.ListDevelopers(tenantID)
}

func (s *UserService) UpdateDeveloper(tenantID, id, name, email string) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	if user.TenantID != tenantID || user.Role != "DEVELOPER" {
		return errors.New("developer not found")
	}
	if name != "" {
		user.Name = name
	}
	if email != "" {
		user.Email = email
	}
	user.UpdatedAt = time.Now()
	return s.userRepo.Update(user)
}

func (s *UserService) ChangePassword(userID, currentPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}
	if len(newPassword) < 6 {
		return errors.New("new password must be at least 6 characters")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hash)
	user.UpdatedAt = time.Now()
	return s.userRepo.Update(user)
}

func (s *UserService) DeactivateDeveloper(tenantID, id string) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	if user.TenantID != tenantID || user.Role != "DEVELOPER" {
		return errors.New("developer not found")
	}
	return s.userRepo.Deactivate(id)
}
