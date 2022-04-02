package zap_log

import (
	"log"
	"path/filepath"
	"time"

	"github.com/wrpota/go-echo/internal/global/variable"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func CreateZapLogFactory(entry func(zapcore.Entry) error, logFileName string) *zap.Logger {

	// 获取程序所处的模式：  开发调试 、 生产
	appDebug := variable.Config.GetBool("AppDebug")

	// 判断程序当前所处的模式，调试模式直接返回一个便捷的zap日志管理器地址，所有的日志打印到控制台即可
	if appDebug == true {
		if logger, err := zap.NewDevelopment(zap.Hooks(entry)); err == nil {
			return logger
		} else {
			log.Fatal("创建zap日志包失败，详情：" + err.Error())
		}
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	timePrecision := variable.Config.GetString("Logs.TimePrecision")
	var recordTimeFormat string
	switch timePrecision {
	case "second":
		recordTimeFormat = "2006-01-02 15:04:05"
	case "millisecond":
		recordTimeFormat = "2006-01-02 15:04:05.000"
	default:
		recordTimeFormat = "2006-01-02 15:04:05"

	}
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(recordTimeFormat))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.TimeKey = "created_at"

	var encoder zapcore.Encoder
	switch variable.Config.GetString("Logs.TextFormat") {
	case "console":
		encoder = zapcore.NewConsoleEncoder(encoderConfig) // 普通模式
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig) // json格式
	default:
		encoder = zapcore.NewConsoleEncoder(encoderConfig) // 普通模式
	}
	LogPath := variable.Config.GetString("Logs.LogPath")
	if LogPath == "" {
		LogPath = variable.BasePath + string(filepath.Separator) + filepath.Join("logs")
	}
	fileName := LogPath + string(filepath.Separator) + logFileName

	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,                                  //日志文件的位置
		MaxSize:    variable.Config.GetInt("Logs.MaxSize"),    //在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: variable.Config.GetInt("Logs.MaxBackups"), //保留旧文件的最大个数
		MaxAge:     variable.Config.GetInt("Logs.MaxAge"),     //保留旧文件的最大天数
		Compress:   variable.Config.GetBool("Logs.Compress"),  //是否压缩/归档旧文件
	}
	writer := zapcore.AddSync(lumberJackLogger)

	zapCore := zapcore.NewCore(encoder, writer, zap.InfoLevel)

	return zap.New(zapCore, zap.AddCaller(), zap.Hooks(entry), zap.AddStacktrace(zap.WarnLevel))
}
