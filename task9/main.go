package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	//读取配置
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
		//将错误打印到控制台然后立刻结束程序
	}

	//连接数据库
	store, err := NewStore(cfg.Database.Path)
	if err != nil {
		log.Fatal("连接数据库失败：", err)
	}
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
		//关闭gin的调试日志输出
	}

	//配置ai
	apiKey := os.Getenv("AI_API_KEY")
	var llm *LLMClient
	if apiKey == "" {
		log.Println("目前未设置AI_API_KEY,/ask接口暂不可用")
	} else {
		llm = NewLLMClient(cfg.LLM.BaseURL, cfg.LLM.Model, apiKey, time.Duration(cfg.LLM.Timeout)*time.Second)
	}

	//注册路由
	h := &Handler{store: store, llm: llm}
	r := gin.Default()
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatal(err)
	}
	r.GET("/todos", h.list)
	r.GET("/ask", h.ask)
	r.POST("/todos", h.create)
	r.DELETE("/todos/:id", h.delete)

	//启动
	log.Println("服务启动，监听", cfg.Server.Port)
	log.Println("应用名称：", cfg.AppName)
	if err := r.Run(cfg.Server.Port); err != nil {
		log.Fatal(err)
	}
}
