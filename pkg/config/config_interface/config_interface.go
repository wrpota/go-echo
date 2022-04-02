// Package config_interface 因config 与 variable相互依赖 添加接口层解决 import cycle not allowed错误
package config_interface

import "time"

type ConfigInterface interface {
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
