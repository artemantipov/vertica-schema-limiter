package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

type SchemaCheck struct {
	Limit float64
	Role  string
}

type VerticaStruct struct {
	Host          string
	Port          string
	User          string
	DB            string
	Pass          string
	CheckInterval int
}

type LimiterConfig struct {
	Vertica VerticaStruct
	Schemas map[string]SchemaCheck
}

var config = ReadConfig()

//ReadConfig read app config file
func ReadConfig() (limiterConf LimiterConfig) {
	conf := viper.New()
	conf.SetEnvPrefix("limiter")
	conf.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	conf.AutomaticEnv()
	conf.SetConfigName("config")
	conf.SetConfigType("yaml")
	conf.AddConfigPath(".")
	conf.AddConfigPath(os.Getenv("CONFIG_PATH"))
	conf.WatchConfig()
	conf.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
		os.Exit(0)
	})
	err := conf.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	err = conf.Unmarshal(&limiterConf)
	if err != nil {
		log.Fatalf("Unable to decode datamon config, %v", err)
	}
	return
}

//GetConfig return config struct
func GetConfig() (limiterConf LimiterConfig) {
	return config
}
