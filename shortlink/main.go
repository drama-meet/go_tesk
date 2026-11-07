package main

import (
	"fmt"
	"log"
	"net/http"
	"shortlink/config"
	"shortlink/handler"
	"shortlink/middle"
	"shortlink/model"
	"shortlink/pkg"
	"shortlink/service"
	"shortlink/store"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	//读取配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("读取配置失败：%v", err)
	}

	//连接数据库
	db, err := gorm.Open(mysql.Open(cfg.Database.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败：%v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取底层连接失败：%v", err)
	}
	sqlDB.SetMaxOpenConns(20)           //同时最多打开20个连接
	sqlDB.SetMaxIdleConns(10)           //空闲时保留10个连接开启
	sqlDB.SetConnMaxLifetime(time.Hour) //单条连接最多活一个小时
	log.Println("数据库连接成功")

	//自动迁移数据库表结构
	if err := db.AutoMigrate(&model.User{}, &model.Link{}); err != nil {
		log.Fatalf("自动迁移数据库表结构失败：%v", err)
	}

	//注册路由
	r := gin.Default()
	//下方代码由ai创建前端时产生
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/", func(c *gin.Context) {
		c.File("log.html")
	})
	r.GET("/log.html", func(c *gin.Context) {
		c.File("log.html")
	})
	r.GET("/links.html", func(c *gin.Context) {
		c.File("links.html")
	})
	r.GET("/public-links.html", func(c *gin.Context) {
		c.File("public-links.html")
	})

	//创建用户&登录&创建短链
	jwtManager := pkg.NewJWTManager(cfg.Auth.Secret, cfg.TokenTTL())
	userStore := store.NewUserStore(db)
	linkStore := store.NewLinkStore(db)
	quotaService := service.NewQuotaService(linkStore, cfg.Quota.DailyCreateLimit)
	userHandler := handler.NewUserHandler(userStore, jwtManager)
	linkHandler := handler.NewLinkHandler(linkStore, quotaService)
	r.POST("/register", userHandler.UserRH)
	r.POST("/login", userHandler.Login)

	r.GET("/s/:code", linkHandler.Redirect)
	r.GET("/public/links", linkHandler.PublicList)

	//将需要登录后才能访问的接口放在一个分组中
	group := r.Group("/")
	group.Use(middle.Middletoken(jwtManager))
	{
		group.GET("/me", func(c *gin.Context) {
			uid, _ := middle.GetUserID(c)
			c.JSON(http.StatusOK, gin.H{"user_id": uid})
		})
		group.POST("/shorten", linkHandler.Create)
		group.GET("/links", linkHandler.List)
		group.GET("/links/:id", linkHandler.Detail)
		group.DELETE("/links/:id", linkHandler.Delete)
	}

	//启动服务
	if err := r.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil { //Run会阻塞当前goroutine，直到服务关闭,不直接输入端口号是因为只接受字符串类型的参数
		log.Fatalf("启动服务失败：%v", err)
	} else {
		log.Printf("服务启动成功，监听端口：%d", cfg.Server.Port)
	}
}
