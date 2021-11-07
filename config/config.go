package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/sirupsen/logrus"
)

type Configuration struct {
	Bot struct {
		Name  string `yaml:"name" env:"INSTA_TG_BOT_NAME" env-upd`
		Token string `yaml:"token" env:"INSTA_TG_BOT_TOKEN" env-required:"true" env-upd`
		Debug bool   `yaml:"debug" env:"INSTA_TG_BOT_DEBUG" env-default:"false" env-upd`
	} `yaml:"bot"`
	Instagram struct {
		AuthLocation string `yaml:"auth-location" env:"INSTAGRAM_AUTH_LOCATION" env-required:"true"`
	} `yaml:"instagram"`
}

var Cfg Configuration

func init() {
	err := cleanenv.ReadConfig("config.yaml", &Cfg)
	if err != nil {
		log.Errorf("error during configuration loading: %v", err)
		panic(err)
	}
}

func updateConfig() {
	err := cleanenv.UpdateEnv(&Cfg)
	if err != nil {
		log.Errorf("configuration wasn't updated: %v", err)
	}
}
