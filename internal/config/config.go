package config

import (
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/viper"
)

// Config хранит все конфигурационные настройки приложения
type Config struct {
	DBDriver  string `mapstructure:"DB_DRIVER"`
	DBSource  string `mapstructure:"DB_SOURCE"`
	DBAddress string `mapstructure:"SERVER_ADDRESS"`

	AdminUsername string `mapstructure:"ADMIN_USERNAME"`
	AdminEmail    string `mapstructure:"ADMIN_EMAIL"`
	AdminPassword string `mapstructure:"ADMIN_PASSWORD"`

	ClientUsername string `mapstructure:"CLIENT_USERNAME"`
	ClientEmail    string `mapstructure:"CLIENT_EMAIL"`
	ClientPassword string `mapstructure:"PROVIDER_PASSWORD"`

	ProviderUsername string `mapstructure:"PROVIDER_USERNAME"`
	ProviderEmail    string `mapstructure:"PROVIDER_EMAIL"`
	ProviderPassword string `mapstructure:"CLIENT_PASSWORD"`

	PartnerUsername string `mapstructure:"PARTNER_USERNAME"`
	PartnerEmail    string `mapstructure:"PARTNER_EMAIL"`
	PartnerPassword string `mapstructure:"PARTNER_PASSWORD"`

	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func GetProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	projectRoot := filepath.Join(dir, "..", "..")
	return projectRoot
}

// LoadConfig считывает конфигурационные параметры из файла или переменных окружения
func LoadConfig() (config Config, err error) {
	projectRoot := GetProjectRoot()
	configPath := filepath.Join(projectRoot, "cmd")

	viper.AddConfigPath(configPath)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
