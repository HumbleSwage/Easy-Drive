package main

import (
	"context"
	"easy-drive/conf"
	"easy-drive/consts"
	"easy-drive/pkg/utils"
	"easy-drive/pkg/utils/commonUtil"
	"easy-drive/repositry/cache"
	"easy-drive/repositry/dao"
	"easy-drive/repositry/model"
	"easy-drive/router"
	"easy-drive/types"
	"errors"
	"os"
	"path/filepath"
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
	go scriptStarting()
}

func scriptStarting() {
	// 启动一些脚本

	// 脚本1：检查删除时间:TODO:删除逻辑需要优化
	go func() {
		// 获取mysql对象
		db := dao.GetMysqlClient(context.Background())

		// 创建定时器，每隔一定时间执行任务
		ticker := time.NewTicker(1 * time.Hour) // 每24小时执行一次任务

		// 在定时器的循环中执行任务
		for {
			select {
			case <-ticker.C:
				// 执行周期性任务
				// 查询并删除过期文件
				thresholdTime := time.Now().AddDate(0, 0, -1) // 当前时间减去10天
				var files []*model.File
				db.Model(&model.File{}).Where("flag IN (?,?) AND restored_at < ?", consts.DeletedFile, consts.RestoreFile, thresholdTime).Find(&files)
				for _, file := range files {
					// 获取操作目录
					workDir, _ := os.Getwd()
					filePath := filepath.Join(workDir, file.Path)

					// 检查是否存在
					_, err := os.Stat(filePath)
					if errors.Is(err, os.ErrNotExist) {
						utils.LogrusObj.Warning("检查本地文件不存在：", err)
						continue
					}

					// 执行删除文件操作
					err = os.Remove(filePath)
					if err != nil {
						utils.LogrusObj.Warning("删除本地文件出错：", err)
						continue
					}

					// 删除记录
					db.Delete(&file)
				}
			}
		}
	}()
}
