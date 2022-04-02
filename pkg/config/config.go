package config

import (
	"flag"

	"github.com/fsnotify/fsnotify"
	"github.com/wrpota/go-echo/internal/global/variable"
	"github.com/wrpota/go-echo/internal/pkg/core/container"
	"github.com/wrpota/go-echo/pkg/config/config_interface"
	"github.com/wrpota/go-echo/pkg/env"

	"log"
	"time"

	"github.com/spf13/viper"
)

const (
	Dev  = "dev"
	Prod = "prod"
	Test = "test"
)

type ymlConfig struct {
	File  string
	ENV   string
	viper *viper.Viper
}

var lastChangeTime time.Time

func init() {
	lastChangeTime = time.Now()
}

func CreateYamlFactory() config_interface.ConfigInterface {

	var configPath string
	// 读取输入目录
	flag.StringVar(&configPath, "conf", "./configs", "配置文件目录")
	flag.Parse()
	// 配置文件所在目录
	if configPath == "" {
		configPath = variable.BasePath + "/configs"
	}

	viper := viper.New()

	viper.AddConfigPath(configPath)

	viper.SetConfigName(env.Active().Value() + "_config")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("初始化配置文件发生错误" + err.Error())
	}

	ymlConfig := &ymlConfig{
		File:  configPath + "/config.yml",
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
	if _, exists := container.CreateContainersFactory().Exists(variable.ConfigPrefix + keyName); exists {
		return true
	} else {
		return false
	}
}

// 对键值进行缓存
func (y *ymlConfig) cache(keyName string, value interface{}) bool {
	return container.CreateContainersFactory().Set(variable.ConfigPrefix+keyName, value)
}

// 通过键获取缓存的值
func (y *ymlConfig) getValueFromCache(keyName string) interface{} {
	return container.CreateContainersFactory().Get(variable.ConfigPrefix + keyName)
}

// 清空已经窜换的配置项信息
func (y *ymlConfig) clearCache() {
	container.CreateContainersFactory().FuzzyDelete(variable.ConfigPrefix)
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
