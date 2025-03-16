package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config хранит все конфигурационные настройки приложения
type Config struct {
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBSource          string `mapstructure:"DB_SOURCE"`
	ServerAddress     string `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`

	// Параметры для администратора
	AdminUsername string `mapstructure:"ADMIN_USERNAME"`
	AdminEmail    string `mapstructure:"ADMIN_EMAIL"`
	AdminPassword string `mapstructure:"ADMIN_PASSWORD"`

	// Параметры для клиента
	ClientUsername string `mapstructure:"CLIENT_USERNAME"`
	ClientEmail    string `mapstructure:"CLIENT_EMAIL"`
	ClientPassword string `mapstructure:"CLIENT_PASSWORD"`

	// Параметры для провайдера услуг
	ProviderUsername string `mapstructure:"PROVIDER_USERNAME"`
	ProviderEmail    string `mapstructure:"PROVIDER_EMAIL"`
	ProviderPassword string `mapstructure:"PROVIDER_PASSWORD"`

	// Параметры для партнера
	PartnerUsername string `mapstructure:"PARTNER_USERNAME"`
	PartnerEmail    string `mapstructure:"PARTNER_EMAIL"`
	PartnerPassword string `mapstructure:"PARTNER_PASSWORD"`

	// Параметры JWT токенов
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

// LoadConfig считывает конфигурационные параметры из файла или переменных окружения
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
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
