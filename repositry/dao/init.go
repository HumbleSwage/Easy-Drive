package dao

import (
	"context"
	"easy-drive/conf"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"strings"
	"time"
)

var _db *gorm.DB

func InitMysql() {
	mConfig := conf.Conf.Mysql["default"]
	dns := strings.Join([]string{mConfig.UserName, ":", mConfig.Password, "@tcp(", mConfig.DbHost, ":", mConfig.DbPort,
		")/", mConfig.DbName, "?charset=", mConfig.Charset, "&parseTime=True&loc=Local"}, "")

	var ormLogger = logger.Default
	if gin.Mode() == "debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dns,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameColumn:   true,
		DontSupportRenameIndex:    true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}

	sqlDb, _ := db.DB()
	sqlDb.SetConnMaxIdleTime(20)
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetConnMaxLifetime(time.Second * 30)

	_db = db
	migration()
}

func GetMysqlClient(ctx context.Context) *gorm.DB {
	return _db.WithContext(ctx)
}
