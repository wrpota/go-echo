package destroy

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/wrpota/go-echo/internal/global/variable"
	event_manage "github.com/wrpota/go-echo/internal/pkg/core/event"
	"go.uber.org/zap"
)

func init() {
	//  用于系统信号的监听
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM) // 监听可能的退出信号
		received := <-c                                                                           //接收信号管道中的值
		variable.ZapLog.Warn("收到信号，进程被结束", zap.String("信号值", received.String()))
		(event_manage.CreateEventManageFactory()).FuzzyCall(variable.EventDestroyPrefix) //调取注册的销毁事件
		close(c)
		os.Exit(1)
	}()

}
