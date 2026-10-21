package main

import (
	"context"
	"fmt"
	"time"
)

type Job struct {
	Name string
	Cost time.Duration
}

var result = make(chan string, 1)

func main() {
	jobs := []Job{
		{Name: "读取配置", Cost: 1 * time.Second},
		{Name: "调用外部服务", Cost: 2 * time.Second},
		{Name: "生成统计报告", Cost: 2 * time.Second},
		{Name: "发送邮件", Cost: 2 * time.Second},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	for _, job := range jobs {
		go work(ctx, job, result)
	}
	for range jobs {
		fmt.Println(<-result)
	}
	fmt.Println("所有任务完成")
}

func work(ctx context.Context, job Job, result chan<- string) {
	select {
	case <-ctx.Done():
		result <- fmt.Sprintf("%s 被取消：%v", job.Name, ctx.Err())
	case <-time.After(job.Cost):
		result <- fmt.Sprintf("%s 任务完成，用时 %v", job.Name, job.Cost)
	}
}
