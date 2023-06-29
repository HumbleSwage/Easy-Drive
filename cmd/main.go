package main

import (
	"easy-drive/conf"
	"easy-drive/pkg/utils"
	"easy-drive/pkg/utils/commonUtil"
	"easy-drive/repositry/cache"
	"easy-drive/repositry/dao"
	"easy-drive/router"
	"easy-drive/types"
	"time"
)

func main() {
	// 加载全局的基础配置文件
	loading()

	r := router.NewRouter()

	err := r.Run("127.0.0.1:7090")
	if err != nil {
		panic(err)
	}
}

func loading() {
	conf.InitConfig()
	dao.InitMysql()
	cache.InitRedis()
	utils.InitLog()
	commonUtil.InitNodeId(time.Now().Format("2006-01-02"), 1)
	types.InitValidate()
}
