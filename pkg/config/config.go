package config

import "github.com/spf13/viper"

type Config struct {
	Port           int    `mapstructure:"PORT"`
	AuthServiceUrl string `mapstructure:"AUTH_SERVICE_URL"`
	ChatServiceUrl string `mapstructure:"CHAT_SERVICE_URL"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./pkg/config/envs")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
