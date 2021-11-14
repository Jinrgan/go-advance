package nacos

import "github.com/nacos-group/nacos-sdk-go/common/constant"

func CreateCltOpts(scs []*constant.ServerConfig, cc *constant.ClientConfig) map[string]interface{} {
	var sc []constant.ServerConfig
	for _, c := range scs {
		sc = append(sc, *c)
	}
	return map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  *cc,
	}
}
