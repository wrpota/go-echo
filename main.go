package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/wrpota/go-echo/configs"
	_ "github.com/wrpota/go-echo/init"
	"github.com/wrpota/go-echo/internal/global/variable"
	"github.com/wrpota/go-echo/internal/router"
	"go.uber.org/zap"
)

func main() {
	s := router.NewHttpServer()
	server := &http.Server{
		Addr:         ":" + configs.Get().GetString("HttpServer.Web.Port"),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	println(os.Getwd())

	if err := s.StartServer(server); err != nil && err != http.ErrServerClosed {
		fmt.Println("http server startup error", err.Error())
		variable.ZapLog.Fatal("http server startup error", zap.Error(err))
	}
}
