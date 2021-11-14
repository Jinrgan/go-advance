package nacos

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type UnmarshalFn func([]byte, interface{}) error

type Configurator struct {
	DataID string
	Group  string
	Client config_client.IConfigClient
}

func (c *Configurator) GetConfig(unmFn func([]byte, UnmarshalFn) error) error {
	s, err := c.Client.GetConfig(vo.ConfigParam{
		DataId: c.DataID,
		Group:  c.Group,
	})
	if err != nil {
		return fmt.Errorf("cannot get config: %v", err)
	}

	err = unmFn([]byte(s), json.Unmarshal)
	if err != nil {
		return fmt.Errorf("cannot unmarshal config: %v", err)
	}

	return nil
}

func (c *Configurator) Listen(unmFn func([]byte, UnmarshalFn) error) error {
	err := c.Client.ListenConfig(vo.ConfigParam{
		DataId: "web",
		Group:  "dev",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("config changed")
			err := unmFn([]byte(data), json.Unmarshal)
			if err != nil {
				zap.L().Error("cannot unmarshal config", zap.Error(err))
			}
			fmt.Printf("group: %s, data id: %s, data: %s", group, dataId, data)
		},
	})

	if err != nil {
		return fmt.Errorf("cannot listen config: %v", err)
	}

	return nil
}
