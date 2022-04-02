package main

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/wrpota/go-echo/init"
	"github.com/wrpota/go-echo/internal/global/variable"
	"github.com/wrpota/go-echo/internal/router"
	"go.uber.org/zap"
)

func main() {
	webRouter := router.InitRouter()

	server := &http.Server{
		Addr:         ":" + variable.Config.GetString("HttpServer.Web.Port"),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	if err := webRouter.StartServer(server); err != nil && err != http.ErrServerClosed {
		fmt.Println("http server startup error", err.Error())
		variable.ZapLog.Fatal("http server startup error", zap.Error(err))
	}
}
