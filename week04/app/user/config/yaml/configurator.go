package yaml

import (
	"fmt"
	"go-advance/week04/app/user/config"
	"go-advance/week04/pkg/net"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	configFileNamePrefix = "config"
)

type Configurator struct {
	viper *viper.Viper
}

func New(name, env, fileDir string) (*Configurator, error) {
	viper.AutomaticEnv()
	_ = viper.GetString("Product") // 刚才设置的环境変量想要生效我们必须得重启 goland

	filePath := fileDir + "/"
	file := []string{configFileNamePrefix, name, fmt.Sprintf("%s.yaml", env)}
	filePath += strings.Join(file, "-")

	v := viper.New()
	v.SetConfigFile(filePath)

	return &Configurator{viper: v}, nil
}

func (c *Configurator) GetConfig() (*config.Server, error) {
	err := c.viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot read in config: %v", err)
	}

	var conf config.Server
	err = c.viper.Unmarshal(&conf)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal server config: %v", err)
	}

	port, err := net.GetPort() // TODO: get full addr
	if err != nil {
		zap.L().Error("cannot get port", zap.Error(err))
	}
	conf.Addr = ":" + strconv.Itoa(port)

	return &conf, nil
}

func (c *Configurator) Listen(conf *config.Server) {
	c.viper.WatchConfig()
	c.viper.OnConfigChange(func(e fsnotify.Event) {
		err := c.viper.ReadInConfig() // 读取配置数据
		if err != nil {
			zap.L().Error("cannot read in config", zap.Error(err))
			return
		}
		zap.L().Info("Config file changed", zap.String("event", e.Name))

		err = c.viper.Unmarshal(conf)
		if err != nil {
			zap.L().Error("cannot unmarshal config", zap.Error(err))
			return
		}
	})
}
