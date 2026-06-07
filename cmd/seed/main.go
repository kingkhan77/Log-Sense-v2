package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kingkhan77/log-sense/internal/models"
	"github.com/kingkhan77/log-sense/internal/repository"
	"github.com/kingkhan77/log-sense/internal/service"
	"github.com/kingkhan77/log-sense/pkg"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := pkg.LoadConfig()
	db := pkg.NewPostgres(cfg)

	tenantID := uuid.NewString()
	adminID := uuid.NewString()
	serviceID := uuid.NewString()
	now := time.Now()

	tenant := &models.Tenant{
		ID:        tenantID,
		Name:      "Demo Corp",
		Status:    "ACTIVE",
		APIKey:    "legacy-not-used",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := db.Create(tenant).Error; err != nil {
		log.Fatalf("tenant: %v", err)
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := &models.User{
		ID:           adminID,
		TenantID:     tenantID,
		Name:         "Admin User",
		Email:        "admin@demo.com",
		PasswordHash: string(hash),
		Role:         "ADMIN",
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := db.Create(admin).Error; err != nil {
		log.Fatalf("admin: %v", err)
	}

	svc := &models.Service{
		ID:          serviceID,
		TenantID:    tenantID,
		Name:        "payment-service",
		Description: "Payment processing",
		IsActive:    true,
		CreatedAt:   now,
	}
	if err := db.Create(svc).Error; err != nil {
		log.Fatalf("service: %v", err)
	}

	apiKeySvc := service.NewAPIKeyService(
		repository.NewAPIKeyRepository(db),
		repository.NewServiceRepository(db),
	)
	rawKey, err := apiKeySvc.Create(tenantID, serviceID, "default-ingest-key")
	if err != nil {
		log.Fatalf("api key: %v", err)
	}

	fmt.Println("Seed complete")
	fmt.Printf("Tenant ID: %s\n", tenantID)
	fmt.Printf("Admin: admin@demo.com / admin123\n")
	fmt.Printf("Service ID: %s\n", serviceID)
	fmt.Printf("API Key (X-API-KEY): %s\n", rawKey)
}
