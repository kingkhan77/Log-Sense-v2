package pkg

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name string `mapstructure:"name"`
		Env  string `mapstructure:"env"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"app"`

	Database struct {
		Host         string `mapstructure:"host"`
		Port         int    `mapstructure:"port"`
		User         string `mapstructure:"user"`
		Password     string `mapstructure:"password"`
		Name         string `mapstructure:"name"`
		SSLMode      string `mapstructure:"sslmode"`
		MaxOpenConns int    `mapstructure:"max_open_conns"`
		MaxIdleConns int    `mapstructure:"max_idle_conns"`
	} `mapstructure:"database"`

	JWT struct {
		Secret                 string `mapstructure:"secret"`
		AccessTokenExpiryHours int    `mapstructure:"access_token_expiry_hours"`
	} `mapstructure:"jwt"`

	Redis struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`

	Kafka struct {
		Brokers []string `mapstructure:"brokers"`
		Topics  struct {
			Logs   string `mapstructure:"logs"`
			Alerts string `mapstructure:"alerts"`
		} `mapstructure:"topics"`
		ConsumerGroups struct {
			Logs   string `mapstructure:"logs"`
			Alerts string `mapstructure:"alerts"`
		} `mapstructure:"consumer_groups"`
	} `mapstructure:"kafka"`

	OpenSearch struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"opensearch"`

	SMTP struct {
		Host      string `mapstructure:"host"`
		Port      int    `mapstructure:"port"`
		Username  string `mapstructure:"username"`
		Password  string `mapstructure:"password"`
		FromEmail string `mapstructure:"from_email"`
	} `mapstructure:"smtp"`

	Alerting struct {
		EvaluationIntervalSeconds int `mapstructure:"evaluation_interval_seconds"`
		DefaultWindowMinutes      int `mapstructure:"default_window_minutes"`
		WorkerCount               int `mapstructure:"worker_count"`
	} `mapstructure:"alerting"`

	Logging struct {
		Level  string `mapstructure:"level"`
		Format string `mapstructure:"format"`
	} `mapstructure:"logging"`

	CORS struct {
		AllowedOrigins []string `mapstructure:"allowed_origins"`
	} `mapstructure:"cors"`
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal(err)
	}

	if cfg.JWT.AccessTokenExpiryHours == 0 {
		cfg.JWT.AccessTokenExpiryHours = 24
	}
	if cfg.Alerting.EvaluationIntervalSeconds == 0 {
		cfg.Alerting.EvaluationIntervalSeconds = 60
	}
	if cfg.Alerting.WorkerCount == 0 {
		cfg.Alerting.WorkerCount = 10
	}
	if cfg.Kafka.Topics.Logs == "" {
		cfg.Kafka.Topics.Logs = "logs"
	}
	if cfg.Kafka.Topics.Alerts == "" {
		cfg.Kafka.Topics.Alerts = "alerts"
	}
	if cfg.Kafka.ConsumerGroups.Logs == "" {
		cfg.Kafka.ConsumerGroups.Logs = "log-processors"
	}
	if cfg.Kafka.ConsumerGroups.Alerts == "" {
		cfg.Kafka.ConsumerGroups.Alerts = "alert-processors"
	}

	return &cfg
}
