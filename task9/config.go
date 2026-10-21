package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string
	}
	Database struct {
		Path string
	}
	Debug   bool
	AppName string `mapstructure:"app_name"`
	LLM     struct {
		BaseURL string `mapstructure:"base_url"`
		Model   string
		Timeout int
	}
}

func loadConfig() (*Config, error) {
	viper.SetConfigName("config")
	//文件名config(不带后缀)
	viper.SetConfigType("yaml")
	//文件后缀
	viper.AddConfigPath(".")
	//在当前目录查找
	viper.AutomaticEnv()
	//允许环境变量覆盖

	//设置默认值，使config.yaml没有时也能启动
	viper.SetDefault("sever.port", ":8080")
	viper.SetDefault("database.path", "todos.db")
	viper.SetDefault("debug", true)
	viper.SetDefault("app_name", "task_9")
	viper.SetDefault("llm.base_url", "https://open.bigmodel.cn/api/paas/v4")
	viper.SetDefault("llm.model", "GLM-4-Flash")
	viper.SetDefault("llm.timeout", 60)

	if err := viper.ReadInConfig(); err != nil {
		//只容忍文件不存在的情况
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("读取配置文件出错：%w", err)
		}
		fmt.Println("提示：未找到config.yaml,已使用默认配置。")
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		//get所有配置到结构体里面
		return nil, fmt.Errorf("解析配置出错：%w", err)
	}
	return &cfg, nil
}
