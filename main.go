package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/m/v2/controller"
	"example.com/m/v2/dao/mysql"
	"example.com/m/v2/dao/redis"
	"example.com/m/v2/logger"
	"example.com/m/v2/pkg/snowflake"
	"example.com/m/v2/routes"
	"example.com/m/v2/settings"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// 1、加载配置
	if err := settings.Init(); err != nil {
		log.Fatalf("load config failed,err:%v\n", err)
		return
	}
	// 2、初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		log.Fatalf("load config failed,err:%v\n", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Debug("logger init success...")
	// 3、初始化MySQL连接
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		log.Fatalf("load config failed,err:%v\n", err)
		return
	}
	defer mysql.Close()
	// 4、初始化Redis连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		log.Fatal("load config failed,err:", zap.Error(err))
		return
	}
	defer redis.Close()
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		log.Fatal("load config failed,err:", zap.Error(err))
		return
	}
	// 初始化gin框架内置的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		zap.L().Fatal("init validator trans failed", zap.Error(err))
		return
	}
	// 5、注册路由
	r := routes.Setup()
	r.SetTrustedProxies([]string{"127.0.0.1"})
	err := r.Run(fmt.Sprintf(":%s", settings.Conf.Port))
	if err != nil {
		log.Fatal("run server failed,err:", zap.Error(err))
		return
	}
	// 6、启动服务(优雅关机)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动自动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen:", zap.Error(err))
		}
	}()
	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用地ctrl+c就是出发系统SINGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时地context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown:", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
