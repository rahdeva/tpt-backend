package config

import "github.com/tkanos/gonfig"

type DWConfiguration struct {
	DB_USERNAME string
	DB_PASSWORD string
	DB_PORT     string
	DB_HOST     string
	DB_NAME     string
}

func GetDWConfig() DWConfiguration {
	conf := DWConfiguration{}
	gonfig.GetConf("config/dw_config.json", &conf)
	return conf
}
