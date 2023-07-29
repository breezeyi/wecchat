package main

import (
	"log"
	"wexin/models"
	"wexin/routers"
	"wexin/utils"

	"github.com/spf13/viper"
)

func main() {
	//加载配置文件
	utils.DispositionInit()

	//加载MySQL配置
	err := models.InitializeMYSQL()
	if err != nil {
		log.Fatalf("MySQL start err :" + err.Error())
	}
	log.Println("mysql start OK!")

	//初始化服务器
	r := routers.Init()
	if err := r.Run(":" + viper.GetString("port")); err != nil {
		log.Fatalln("服务器端口启动错误！")
	}
}
