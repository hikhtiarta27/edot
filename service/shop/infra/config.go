package infra

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	DB struct {
		Host  string `mapstructure:"host"`
		User  string `mapstructure:"user"`
		Pass  string `mapstructure:"pass"`
		Port  uint   `mapstructure:"port"`
		Name  string `mapstructure:"name"`
		Param string `mapstructure:"param"`
	} `mapstructure:"db"`
	JWT struct {
		Secret string `mapstructure:"secret"`
	} `mapstructure:"jwt"`
	App struct {
		Address string `mapstructure:"address"`
	} `mapstructure:"app"`
	Service struct {
		Warehouse string `mapstructure:"warehouse"`
	} `mapstructure:"service"`
}

var (
	configOnce sync.Once
	config     *Config
)

func LoadConfig() *Config {

	configOnce.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")

		err := viper.ReadInConfig()
		if err != nil {
			log.Fatal(err)
		}

		err = viper.Unmarshal(&config)
		if err != nil {
			return
		}

	})

	return config
}
