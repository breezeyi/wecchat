package utils

import (
	"log"

	"github.com/spf13/viper"
)

func DispositionInit() {
	viper.SetConfigFile("./config/conf.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("读取配置文件错误！" + err.Error())
	} else {
		log.Println("配置文件加载成功！")
	}
}
