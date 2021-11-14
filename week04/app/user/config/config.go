package config

import "go-advance/week04/app/user/config/nacos"

type Configurator interface {
	GetConfig(unmFn func([]byte, nacos.UnmarshalFn) error) error
	Listen(unmFn func([]byte, nacos.UnmarshalFn) error) error
}

type Sms struct {
	Expire int `mapstructure:"expire"`
}

type Redis struct {
	Addr string `mapstructure:"addr"`
}

type Mysql struct {
	Addr     string `mapstructure:"addr"`
	DB       string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type Consul struct {
	Addr string `mapstructure:"addr"`
}
