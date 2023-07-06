package Utils

import (
	"fmt"

	"github.com/spf13/viper"
)

var config *viper.Viper

//func init() {
//	config = viper.New()
//	//config.SetConfigName("go_config")
//	//config.SetConfigType("properties")
//	//config.AddConfigPath(".")
//	//config.AddConfigPath("/Users/kaoweicheng/Desktop/workspace/GO/")
//	if err := config.ReadInConfig(); err != nil {
//		panic(fmt.Errorf("fatal error config file: %s", err))
//	}
//}

func InitProperties(path string, configName string, configType string) {
	config = viper.New()
	config.AddConfigPath(path) // 添加路径参数到 viper 实例中
	config.SetConfigName(configName)
	config.SetConfigType(configType)
	if err := config.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	LogUtil.Info("===properties initialize success===")
}

func GetProperties(key string) string {
	return config.GetString(key)
}
