package config

type Service struct {
	Name  string `mapstructure:"name"`
	Addr  string `mapstructure:"addr"`
	Mysql Mysql  `json:"mysql"`
}
