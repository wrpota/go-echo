package init

import (
	"log"

	"github.com/wrpota/go-echo/internal/global/variable"
	_ "github.com/wrpota/go-echo/internal/pkg/core/destroy" //退出信号监听
	"github.com/wrpota/go-echo/internal/pkg/log_hook"
	"github.com/wrpota/go-echo/pkg/config"
	cgorm "github.com/wrpota/go-echo/pkg/grom"
	"github.com/wrpota/go-echo/pkg/zap_log"
)

func init() {
	variable.Config = config.CreateYamlFactory()

	// 初始化全局日志句柄，并载入日志钩子处理函数
	variable.ZapLog = zap_log.CreateZapLogFactory(log_hook.ZapLogHandler, variable.Config.GetString("Logs.GoLogName"))
	variable.EchoZapLog = zap_log.CreateZapLogFactory(log_hook.ZapLogHandler, variable.Config.GetString("Logs.EchoLogName"))
	//初始化数据库连接
	if dbRead, err := cgorm.GetDbReadDriver(variable.Config.GetString("Database.UseDbType")); err != nil {
		log.Fatal("Gorm 数据库驱动、连接初始化失败" + err.Error())
	} else {
		variable.GormReadMysql = dbRead
	}
	if dbWrite, err := cgorm.GetDbWriteDriver(variable.Config.GetString("Database.UseDbType")); err != nil {
		log.Fatal("Gorm 数据库驱动、连接初始化失败" + err.Error())
	} else {
		variable.GormWriteDb = dbWrite
	}
}
