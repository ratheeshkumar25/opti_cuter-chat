package config

import "github.com/spf13/viper"

type Config struct {
	DBurl      string `mapstructure:"DBURL"`
	DBName     string `mapstructure:"DBNAME"`
	GrpcPort   string `mapstructure:"GRPCPORT"`
	OpenApiKey string `mapstructure:"OPEN_APIKEY"`
}

func LoadConfig() *Config {
	var config Config
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	viper.Unmarshal(&config)
	return &config
}
