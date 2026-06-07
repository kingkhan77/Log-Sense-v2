package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/kingkhan77/log-sense/internal/consumer"
	"github.com/kingkhan77/log-sense/internal/controller"
	"github.com/kingkhan77/log-sense/internal/engine"
	"github.com/kingkhan77/log-sense/internal/middleware"
	"github.com/kingkhan77/log-sense/internal/repository"
	"github.com/kingkhan77/log-sense/internal/service"
	"github.com/kingkhan77/log-sense/pkg"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := pkg.LoadConfig()
	logger := pkg.NewLogger()
	middleware.InitAuth(cfg)

	db := pkg.NewPostgres(cfg)
	redisClient := pkg.NewRedis(cfg)

	kafkaProducer, err := pkg.NewKafkaProducer(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize kafka producer")
	}

	osClient, err := pkg.NewOpenSearch(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize opensearch")
	}

	userRepo := repository.NewUserRepository(db)
	serviceRepo := repository.NewServiceRepository(db)
	ruleRepo := repository.NewRuleRepository(db)
	alertRepo := repository.NewAlertRepository(db)
	oncallRepo := repository.NewOnCallRepository(db)
	apiKeyRepo := repository.NewAPIKeyRepository(db)

	authService := service.NewAuthService(userRepo, cfg)
	userService := service.NewUserService(userRepo)
	serviceService := service.NewServiceService(serviceRepo)
	ruleService := service.NewRuleService(ruleRepo)
	alertService := service.NewAlertService(alertRepo, redisClient)
	oncallService := service.NewOnCallService(oncallRepo, serviceRepo)
	apiKeyService := service.NewAPIKeyService(apiKeyRepo, serviceRepo)
	ingestionService := service.NewIngestionService(kafkaProducer, cfg)
	alertProducer := service.NewAlertProducer(kafkaProducer, cfg)
	notificationService := service.NewNotificationService(alertRepo, oncallRepo, userRepo, serviceRepo, cfg)
	dashboardService := service.NewDashboardService(alertRepo, ruleRepo, serviceRepo)

	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)
	serviceController := controller.NewServiceController(serviceService)
	ruleController := controller.NewRuleController(ruleService)
	alertController := controller.NewAlertController(alertService)
	oncallController := controller.NewOnCallController(oncallService)
	apiKeyController := controller.NewAPIKeyController(apiKeyService)
	ingestionController := controller.NewIngestionController(ingestionService)
	dashboardController := controller.NewDashboardController(dashboardService)
	healthController := controller.NewHealthController(db, redisClient, cfg)
	logController := controller.NewLogController(osClient)

	ruleEngine := engine.NewRuleEngine(ruleRepo, alertRepo, osClient, redisClient, alertProducer, cfg)
	ruleEngine.WarmAlertCache()

	logConsumer := consumer.NewLogConsumer(osClient, cfg)
	notifyConsumer := consumer.NewNotificationConsumer(notificationService, cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go ruleEngine.Start(ctx)
	go logConsumer.Start(ctx, cfg.Kafka.Brokers)
	go notifyConsumer.Start(ctx, cfg.Kafka.Brokers)

	router := gin.Default()
	router.Use(middleware.CORS(cfg.CORS.AllowedOrigins))

	router.GET("/health", healthController.Health)

	api := router.Group("/api/v1")

	api.POST("/login", authController.Login)

	logs := api.Group("/logs")
	logs.Use(middleware.APIKeyAuth(apiKeyRepo))
	logs.POST("", ingestionController.IngestLog)

	protected := api.Group("/")
	protected.Use(middleware.Auth())
	protected.PUT("/me/password", userController.ChangePassword)

	admin := protected.Group("/admin")
	admin.Use(middleware.Admin())
	{
		admin.POST("/developers", userController.CreateDeveloper)
		admin.GET("/developers", userController.ListDevelopers)
		admin.PUT("/developers/:id", userController.UpdateDeveloper)
		admin.DELETE("/developers/:id", userController.DeleteDeveloper)
		admin.POST("/oncall/schedules", oncallController.CreateSchedule)
		admin.GET("/oncall/schedules", oncallController.ListSchedules)
		admin.PUT("/oncall/schedules/:id", oncallController.UpdateSchedule)
		admin.DELETE("/oncall/schedules/:id", oncallController.DeleteSchedule)
		admin.GET("/oncall/current/:serviceId", oncallController.GetCurrentOnCall)
		admin.GET("/api-keys", apiKeyController.ListKeys)
		admin.POST("/api-keys", apiKeyController.CreateKey)
		admin.DELETE("/api-keys/:id", apiKeyController.RevokeKey)
	}

	services := protected.Group("/services")
	{
		services.GET("", serviceController.ListServices)
		services.GET("/:id", serviceController.GetService)
		services.PUT("/:id", serviceController.UpdateService)
		services.DELETE("/:id", serviceController.DeleteService)

		adminServices := services.Group("")
		adminServices.Use(middleware.Admin())
		adminServices.POST("", serviceController.CreateService)
	}

	rules := protected.Group("/rules")
	{
		rules.POST("", ruleController.CreateRule)
		rules.GET("", ruleController.ListRules)
		rules.GET("/:id", ruleController.GetRule)
		rules.PUT("/:id", ruleController.UpdateRule)
		rules.DELETE("/:id", ruleController.DeleteRule)
	}

	alerts := protected.Group("/alerts")
	{
		alerts.GET("", alertController.ListAlerts)
		alerts.GET("/:id", alertController.GetAlert)
		alerts.POST("/:id/ack", alertController.AcknowledgeAlert)
		alerts.POST("/:id/resolve", alertController.ResolveAlert)
	}

	dashboard := protected.Group("/dashboard")
	dashboard.GET("/summary", dashboardController.Summary)

	protected.GET("/logs", logController.Search)

	port := cfg.App.Port
	if port == 0 {
		port = 8080
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", strconv.Itoa(port)),
		Handler: router,
	}

	go func() {
		logger.Info().Msgf("server started on :%d", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("shutting down")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
	_ = redisClient.Close()
}
