package container

import (
	"log"
	"strings"
	"sync"

	"github.com/wrpota/go-echo/internal/global/variable"
)

var sMap sync.Map

// 创建容器工厂
func CreateContainersFactory() *containers {
	return &containers{}
}

// 容器结构体
type containers struct {
}

// 以键值对的形式将代码注册到容器
func (c *containers) Set(key string, value interface{}) (res bool) {

	if _, exists := c.Exists(key); exists == false {
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

func (c *containers) Get(key string) interface{} {
	if value, exists := c.Exists(key); exists {
		return value
	}
	return nil
}

// 是否已注册
func (c *containers) Exists(key string) (interface{}, bool) {
	return sMap.Load(key)
}

//模糊删除删除key
func (c *containers) FuzzyDelete(keyPrefix string) {
	sMap.Range(func(key, value interface{}) bool {
		if keyname, ok := key.(string); ok {
			if strings.HasPrefix(keyname, keyPrefix) {
				sMap.Delete(keyname)
			}
		}
		return true
	})
}
