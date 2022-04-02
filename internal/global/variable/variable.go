package variable

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/wrpota/go-echo/pkg/config/config_interface"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	BasePath string //项目根地址

	ConfigPrefix       = "Echo:"    //缓存key前缀
	EventDestroyPrefix = "Destroy_" //  程序退出时需要销毁的事件前缀

	ConfigPath string //配置文件目录

	ZapLog     *zap.Logger //日志指针
	EchoZapLog *zap.Logger //日志指针
	// Config 配置
	Config config_interface.ConfigInterface

	GormReadMysql *gorm.DB
	GormWriteDb   *gorm.DB
)

func init() {
	// 初始化程序根目录
	if path, err := os.Getwd(); err == nil {
		if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-test") {
			BasePath = strings.Replace(path, string(filepath.Separator)+`test`, "", 1)
		} else {
			BasePath = path
		}
	} else {
		log.Fatal("获取程序允许目录失败")
	}
}
