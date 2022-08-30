/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-25
* Time: 09:59
 */

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gowebsocket/lib/redislib"
	"gowebsocket/routers"
	"gowebsocket/servers/grpcserver"
	"gowebsocket/servers/task"
	"gowebsocket/servers/websocket"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main() {
	initConfig()

	initFile()

	//实例化 redis.Client
	initRedis()

	router := gin.Default()
	// 初始化http模板views及路由
	routers.Init(router)

	//绑定websocket登录，心跳，ping函数到handlers属性-map[string]DisposeFunc
	routers.WebsocketInit()

	// 定时任务（注册定时清理连接任务）
	task.Init()

	// 实例化 server_model （服务ip:port redis配置）
	task.ServerInit()

	// 绑定/acc路由事件（当有websocket连接时，给它分别启动读、写socket协程，具体操作通过socket传送的自定义参数cmd确定，如login，heartbeat）
	// 启动ClientManager协程事件，监听注册，发送数据事件
	// 启动websocket服务
	go websocket.StartWebSocket()

	// grpc
	go grpcserver.Init()

	//访问页面
	go open()

	httpPort := viper.GetString("app.httpPort")

	//启动http服务
	http.ListenAndServe(":"+httpPort, router)

}

// 初始化日志
func initFile() {
	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor()

	// Logging to a file.
	logFile := viper.GetString("app.logFile")
	f, _ := os.Create(logFile)
	gin.DefaultWriter = io.MultiWriter(f)
}

func initConfig() {
	viper.SetConfigName("config/app")
	viper.AddConfigPath(".") // 添加搜索路径

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	fmt.Println("config app:", viper.Get("app"))
	fmt.Println("config redis:", viper.Get("redis"))

}

func initRedis() {
	redislib.ExampleNewClient()
}

func open() {

	time.Sleep(1000 * time.Millisecond)

	httpUrl := viper.GetString("app.httpUrl")
	httpUrl = "http://" + httpUrl + "/home/index"

	fmt.Println("访问页面体验:", httpUrl)

	cmd := exec.Command("open", httpUrl)
	cmd.Output()
}
