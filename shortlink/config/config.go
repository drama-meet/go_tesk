package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Quota    QuotaConfig
	Auth     AuthConfig
}
type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Dsn string
}

type QuotaConfig struct {
	DailyCreateLimit int `mapstructure:"daily_create_limit"`
}

type AuthConfig struct {
	Secret         string `mapstructure:"jwt_secret"`
	TokenTotalTime int    `mapstructure:"token_ttl_hour"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("config")
	v.AutomaticEnv() //允许环境变量覆盖

	//设置非须保密的默认值
	v.SetDefault("quota.daily_create_limit", 10)
	v.SetDefault("server.port", 8080)
	v.SetDefault("auth.token_ttl_hour", 24)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败：%w", err)
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败：%w", err)
	}
	return &cfg, nil
}

func (c *Config) TokenTTL() time.Duration {
	return time.Duration(c.Auth.TokenTotalTime) * time.Hour
}
