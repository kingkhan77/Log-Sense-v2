package service

import (
	"errors"
	"time"

	"github.com/kingkhan77/log-sense/internal/repository"
	"github.com/kingkhan77/log-sense/pkg"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
	cfg      *pkg.Config
}

func NewAuthService(
	userRepo *repository.UserRepository,
	cfg *pkg.Config,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !user.IsActive {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	hours := s.cfg.JWT.AccessTokenExpiryHours
	if hours <= 0 {
		hours = 24
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id":   user.ID,
			"tenant_id": user.TenantID,
			"role":      user.Role,
			"exp":       time.Now().Add(time.Duration(hours) * time.Hour).Unix(),
		},
	)

	return token.SignedString([]byte(s.cfg.JWT.Secret))
}
