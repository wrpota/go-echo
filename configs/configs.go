package configs

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/wrpota/go-echo/internal/global/variable"
	"github.com/wrpota/go-echo/pkg/env"
	"github.com/wrpota/go-echo/pkg/file"

	_ "embed"
	"log"
	"time"

	"github.com/spf13/viper"
)

var (
	//go:embed dev_configs.yml
	devConfigs []byte

	//go:embed fat_configs.yml
	fatConfigs []byte

	//go:embed pro_configs.yml
	proConfigs []byte
)

var lastChangeTime time.Time

func init() {
	lastChangeTime = time.Now()
}

var sMap sync.Map

var once sync.Once

const (
	ConfigPrefix = "config:"
)

var instance *ymlConfig

func Get() *ymlConfig {
	once.Do(func() {
		instance = createConfig().(*ymlConfig)
	})
	return instance
}

type ymlConfig struct {
	File  string
	ENV   string
	viper *viper.Viper
}

type configInterface interface {
	ConfigFileChangeListen()
	Get(keyName string) interface{}
	GetString(keyName string) string
	GetBool(keyName string) bool
	GetInt(keyName string) int
	GetInt32(keyName string) int32
	GetInt64(keyName string) int64
	GetFloat64(keyName string) float64
	GetDuration(keyName string) time.Duration
	GetStringSlice(keyName string) []string
}

var _ configInterface = (*ymlConfig)(nil)

func createConfig() configInterface {
	var r io.Reader
	switch env.Active().Value() {
	case "dev":
		r = bytes.NewReader(devConfigs)
	case "fat":
		r = bytes.NewReader(fatConfigs)
	case "pro":
		r = bytes.NewReader(proConfigs)
	default:
		r = bytes.NewReader(devConfigs)
	}

	viper := viper.New()

	viper.SetConfigType("yml")
	viper.SetConfigName(env.Active().Value() + "_configs")
	viper.AddConfigPath("./configs")
	configFile := "./configs/" + env.Active().Value() + "_configs.yml"
	_, ok := file.IsExists(configFile)
	if !ok {
		if err := os.MkdirAll(filepath.Dir(configFile), 0766); err != nil {
			panic(err)
		}

		f, err := os.Create(configFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		if b, err := io.ReadAll(r); err == nil {
			if err := os.WriteFile(configFile, b, 0766); err != nil {
				panic(err)
			}
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("初始化配置文件发生错误" + err.Error())
	}

	ymlConfig := &ymlConfig{
		File:  env.Active().Value() + "_configs.yml",
		ENV:   env.Active().Value(),
		viper: viper,
	}
	ymlConfig.ConfigFileChangeListen()

	return ymlConfig
}

// ConfigFileChangeListen 监听文件变化
func (y *ymlConfig) ConfigFileChangeListen() {
	y.viper.OnConfigChange(func(changeEvent fsnotify.Event) {
		if time.Now().Sub(lastChangeTime).Seconds() >= 1 {
			if changeEvent.Op.String() == "WRITE" {
				y.clearCache()
				lastChangeTime = time.Now()
			}
		}
	})
	y.viper.WatchConfig()
}

// 判断相关键是否已经缓存
func (y *ymlConfig) keyIsCache(keyName string) bool {
	if _, exists := y.Exists(ConfigPrefix + keyName); exists {
		return true
	} else {
		return false
	}
}

func (y *ymlConfig) Exists(key string) (interface{}, bool) {
	return sMap.Load(key)
}

// 对键值进行缓存
func (y *ymlConfig) cache(keyName string, value interface{}) (res bool) {
	var key = ConfigPrefix + keyName
	if _, exists := y.Exists(key); exists == false {
		sMap.Store(key, value)
		res = true
	} else {
		// 程序启动阶段，zaplog 未初始化，使用系统log打印启动时候发生的异常日志
		if variable.ZapLog == nil {
			log.Fatal("请解决键名重复问题,相关键：" + key)
		} else {
			// 程序启动初始化完成
			variable.ZapLog.Warn("相关键：" + key)
		}
	}
	return
}

// 通过键获取缓存的值
func (y *ymlConfig) getValueFromCache(keyName string) interface{} {
	if value, exists := y.Exists(ConfigPrefix + keyName); exists {
		return value
	}
	return nil
}

// 清空已经窜换的配置项信息
func (y *ymlConfig) clearCache() {
	sMap.Range(func(key, value interface{}) bool {
		if keyname, ok := key.(string); ok {
			if strings.HasPrefix(keyname, ConfigPrefix) {
				sMap.Delete(keyname)
			}
		}
		return true
	})
}

// Get 一个原始值
func (y *ymlConfig) Get(keyName string) interface{} {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName)
	} else {
		value := y.viper.Get(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetString
func (y *ymlConfig) GetString(keyName string) string {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(string)
	} else {
		value := y.viper.GetString(keyName)
		y.cache(keyName, value)
		return value
	}

}

// GetBool
func (y *ymlConfig) GetBool(keyName string) bool {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(bool)
	} else {
		value := y.viper.GetBool(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetInt
func (y *ymlConfig) GetInt(keyName string) int {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(int)
	} else {
		value := y.viper.GetInt(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetInt32
func (y *ymlConfig) GetInt32(keyName string) int32 {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(int32)
	} else {
		value := y.viper.GetInt32(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetInt64
func (y *ymlConfig) GetInt64(keyName string) int64 {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(int64)
	} else {
		value := y.viper.GetInt64(keyName)
		y.cache(keyName, value)
		return value
	}
}

// float64
func (y *ymlConfig) GetFloat64(keyName string) float64 {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(float64)
	} else {
		value := y.viper.GetFloat64(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetDuration
func (y *ymlConfig) GetDuration(keyName string) time.Duration {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(time.Duration)
	} else {
		value := y.viper.GetDuration(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetStringSlice
func (y *ymlConfig) GetStringSlice(keyName string) []string {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).([]string)
	} else {
		value := y.viper.GetStringSlice(keyName)
		y.cache(keyName, value)
		return value
	}
}
