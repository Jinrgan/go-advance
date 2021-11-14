package config

type Server struct {
	Name        string  `mapstructure:"name"`
	Addr        string  `mapstructure:"addr"`
	ServiceName string  `mapstructure:"service-name"`
	Sms         *Sms    `mapstructure:"sms"`
	Redis       *Redis  `mapstructure:"redis"`
	Mysql       *Mysql  `mapstructure:"mysql" json:"mysql"`
	Consul      *Consul `mapstructure:"consul" json:"consul"`
}
