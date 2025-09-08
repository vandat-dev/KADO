package initialize

import (
	_ "base_go_be/docs"
	"base_go_be/global"
	"fmt"

	"go.uber.org/zap"
)

func handleErr(err error) {
	if err != nil {
		global.Logger.Error("Run app fail!", zap.Error(err))
		panic(err)
	}
}

func Run() {
	LoadConfig()
	InitLogger()
	//global.Logger.Info("check logger", zap.String("key", "value"))
	//Mysql()
	Postgres()
	Redis()
	InitWebSocketManager()

	r := InitRouter()
	port := fmt.Sprintf(":%d", global.Config.Server.Port)
	global.Logger.Info("Server starting on port", zap.String("port", port))
	err := r.Run(port)
	handleErr(err)
}
