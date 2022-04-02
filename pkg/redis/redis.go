package redis

import (
	"time"

	"github.com/wrpota/go-echo/configs"
	"github.com/wrpota/go-echo/internal/global/variable"
	event_manage "github.com/wrpota/go-echo/internal/pkg/event"

	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
)

var redisPool *redis.Pool

// 处于程序底层的包，init 初始化的代码段的执行会优先于上层代码，因此这里读取配置项不能使用全局配置项变量
func init() {
	redisPool = initRedisClientPool()
}

func initRedisClientPool() *redis.Pool {
	redisPool = &redis.Pool{
		MaxIdle:     configs.Get().GetInt("Redis.MaxIdle"),                        //最大空闲数
		MaxActive:   configs.Get().GetInt("Redis.MaxActive"),                      //最大活跃数
		IdleTimeout: configs.Get().GetDuration("Redis.IdleTimeout") * time.Second, //最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
		Dial: func() (redis.Conn, error) {
			//此处对应redis ip及端口号
			conn, err := redis.Dial("tcp", configs.Get().GetString("Redis.Host")+":"+configs.Get().GetString("Redis.Port"))
			if err != nil {
				variable.ZapLog.Error("初始化redis连接池失败" + err.Error())
				return nil, err
			}
			auth := configs.Get().GetString("Redis.Auth") //通过配置项设置redis密码
			if len(auth) >= 1 {
				if _, err := conn.Do("AUTH", auth); err != nil {
					_ = conn.Close()
					variable.ZapLog.Error("Redis Auth 鉴权失败，密码错误" + err.Error())
				}
			}
			_, _ = conn.Do("select", configs.Get().GetInt("Redis.IndexDb"))
			return conn, err
		},
	}
	// 将redis的关闭事件，注册在全局事件统一管理器，由程序退出时统一销毁
	event_manage.CreateEventManageFactory().Set(variable.EventDestroyPrefix+"Redis", func(args ...interface{}) {
		_ = redisPool.Close()
	})
	return redisPool
}

// GetOneRedisClient 从连接池获取一个redis连接
func GetOneRedisClient() *Client {
	maxRetryTimes := configs.Get().GetInt("Redis.ConnFailRetryTimes")
	var oneConn redis.Conn
	for i := 1; i <= maxRetryTimes; i++ {
		oneConn = redisPool.Get()
		if oneConn.Err() != nil {
			// variable.ZapLog.Error("Redis：网络中断,开始重连进行中...", zap.Error(oneConn.Err()))
			if i == maxRetryTimes {
				variable.ZapLog.Error("Redis 从连接池获取一个连接失败，超过最大重试次数", zap.Error(oneConn.Err()))
				return nil
			}
			//如果出现网络短暂的抖动，短暂休眠后，支持自动重连
			time.Sleep(time.Second * configs.Get().GetDuration("Redis.ReConnectInterval"))
		} else {
			break
		}
	}
	return &Client{oneConn}
}

// Client 定义一个redis客户端结构体
type Client struct {
	client redis.Conn
}

// Get
func (r *Client) Get(args ...interface{}) (interface{}, error) {
	return r.Execute("get", args...)
}

// Set
func (r *Client) Set(args ...interface{}) (interface{}, error) {
	return r.Execute("set", args...)
}

// Execute 为redis-go 客户端封装统一操作函数入口
func (r *Client) Execute(cmd string, args ...interface{}) (interface{}, error) {
	return r.client.Do(cmd, args...)
}

// ReleaseOneRedisClient 释放连接到连接池
func (r *Client) ReleaseOneRedisClient() {
	_ = r.client.Close()
}

// Flush 清空
func (r *Client) Flush() error {
	return r.client.Flush()
}

//  封装几个数据类型转换的函数

// Bool 类型转换
func (r *Client) Bool(reply interface{}, err error) (bool, error) {
	return redis.Bool(reply, err)
}

//string 类型转换
func (r *Client) String(reply interface{}, err error) (string, error) {
	return redis.String(reply, err)
}

// Strings strings 类型转换
func (r *Client) Strings(reply interface{}, err error) ([]string, error) {
	return redis.Strings(reply, err)
}

//Float64 类型转换
func (r *Client) Float64(reply interface{}, err error) (float64, error) {
	return redis.Float64(reply, err)
}

// Int int 类型转换
func (r *Client) Int(reply interface{}, err error) (int, error) {
	return redis.Int(reply, err)
}

// Int64 int64 类型转换
func (r *Client) Int64(reply interface{}, err error) (int64, error) {
	return redis.Int64(reply, err)
}

// Uint64 uint64 类型转换
func (r *Client) Uint64(reply interface{}, err error) (uint64, error) {
	return redis.Uint64(reply, err)
}

//Bytes 类型转换
func (r *Client) Bytes(reply interface{}, err error) ([]byte, error) {
	return redis.Bytes(reply, err)
}
