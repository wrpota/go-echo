package gorm

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/wrpota/go-echo/configs"
	"github.com/wrpota/go-echo/internal/global/variable"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
)

func GetDbWriteDriver(sqlType string) (*gorm.DB, error) {
	var dbDialector gorm.Dialector
	if val, err := getDbDialector(sqlType, "Write"); err != nil {
		variable.ZapLog.Error("gorm dialector 初始化失败,dbType:"+sqlType, zap.Error(err))
	} else {
		dbDialector = val
	}
	gormDb, err := gorm.Open(dbDialector, &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 redefineLog(sqlType), //拦截、接管 gorm v2 自带日志
	})
	if err != nil {
		return nil, err
	}

	// 为主连接设置连接池
	if rawDb, err := gormDb.DB(); err != nil {
		return nil, err
	} else {
		rawDb.SetConnMaxIdleTime(time.Second * 30)
		rawDb.SetConnMaxLifetime(configs.Get().GetDuration("Database.Write.SetConnMaxLifetime") * time.Second)
		rawDb.SetMaxIdleConns(configs.Get().GetInt("Database.Write.SetMaxIdleConns"))
		rawDb.SetMaxOpenConns(configs.Get().GetInt("Database.Write.SetMaxOpenConns"))
		return gormDb, nil
	}
}

func GetDbReadDriver(sqlType string) (*gorm.DB, error) {
	var dbDialector gorm.Dialector
	if val, err := getDbDialector(sqlType, "Read"); err != nil {
		variable.ZapLog.Error("gorm dialector 初始化失败,dbType:"+sqlType, zap.Error(err))
	} else {
		dbDialector = val
	}
	gormDb, err := gorm.Open(dbDialector, &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 redefineLog(sqlType), //拦截、接管 gorm v2 自带日志
	})
	if err != nil {
		return nil, err
	}

	// 查询没有数据，屏蔽 gorm v2 包中会爆出的错误
	// https://github.com/go-gorm/gorm/issues/3789  此 issue 所反映的问题就是我们本次解决掉的
	_ = gormDb.Callback().Query().Before("gorm:query").Register("disable_raise_record_not_found", func(d *gorm.DB) {
		d.Statement.RaiseErrorOnNotFound = false
	})

	// 为主连接设置连接池
	if rawDb, err := gormDb.DB(); err != nil {
		return nil, err
	} else {
		rawDb.SetConnMaxIdleTime(time.Second * 30)
		rawDb.SetConnMaxLifetime(configs.Get().GetDuration("Database.Write.SetConnMaxLifetime") * time.Second)
		rawDb.SetMaxIdleConns(configs.Get().GetInt("Database.Write.SetMaxIdleConns"))
		rawDb.SetMaxOpenConns(configs.Get().GetInt("Database.Write.SetMaxOpenConns"))
		return gormDb, nil
	}
}

// 获取一个数据库方言(Dialector),通俗的说就是根据不同的连接参数，获取具体的一类数据库的连接指针
func getDbDialector(sqlType, readWrite string) (gorm.Dialector, error) {
	var dbDialector gorm.Dialector
	dsn := getDsn(sqlType, readWrite)
	switch strings.ToLower(sqlType) {
	case "mysql":
		dbDialector = mysql.Open(dsn)
	case "sqlserver", "mssql":
		dbDialector = sqlserver.Open(dsn)
	case "postgres", "postgresql", "postgre":
		dbDialector = postgres.Open(dsn)
	default:
		return nil, errors.New("数据库驱动类型不存在,目前支持的数据库类型：mysql、sqlserver、postgresql，您提交数据库类型：" + sqlType)
	}
	return dbDialector, nil
}

func getDsn(sqlType, readWrite string) string {
	Host := configs.Get().GetString("Database." + readWrite + ".Host")
	DataBase := configs.Get().GetString("Database." + readWrite + ".DataBase")
	Port := configs.Get().GetInt("Database." + readWrite + ".Port")
	User := configs.Get().GetString("Database." + readWrite + ".User")
	Pass := configs.Get().GetString("Database." + readWrite + ".Pass")
	Charset := configs.Get().GetString("Database." + readWrite + ".Charset")

	switch strings.ToLower(sqlType) {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", User, Pass, Host, Port, DataBase, Charset)
	case "sqlserver", "mssql":
		return fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s;encrypt=disable", Host, Port, DataBase, User, Pass)
	case "postgresql", "postgre", "postgres":
		return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable TimeZone=Asia/Shanghai", Host, Port, DataBase, User, Pass)
	}
	return ""
}

// 创建自定义日志模块，对 gorm 日志进行拦截、
func redefineLog(sqlType string) gormLog.Interface {
	return createCustomGormLog(sqlType,
		SetInfoStrFormat("[info] %s\n"), SetWarnStrFormat("[warn] %s\n"), SetErrStrFormat("[error] %s\n"),
		SetTraceStrFormat("[traceStr] %s [%.3fms] [rows:%v] %s\n"), SetTracWarnStrFormat("[traceWarn] %s %s [%.3fms] [rows:%v] %s\n"), SetTracErrStrFormat("[traceErr] %s %s [%.3fms] [rows:%v] %s\n"))
}
