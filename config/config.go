package config

import (
	"github.com/spf13/viper"
)

type (
	config struct {
		BaseURL        string `mapstructure:"BASE_URL"`
		NumRequests    int    `mapstructure:"NUM_REQUESTS"`
		Concurrency    int    `mapstructure:"CONCURRENCY"`
		TokenListFile  string `mapstructure:"TOKEN_LIST_FILE"`
		UrlListFile    string `mapstructure:"URL_LIST_FILE"`
		RequestTimeout int    `mapstructure:"REQUEST_TIMEOUT"`
	}
)

var (
	Config *config
)

func init() {
	// === Set the config file name and type
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")

	// === Read the config file
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// === Load .env file into global variables
	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(err)
	}
}
