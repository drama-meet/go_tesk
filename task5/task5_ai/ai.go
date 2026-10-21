// 这里个人用的智谱的api_key,代码没对输出字数进行限制。
// 解析的部分比较简单，不太完善，这个是给task9任务练手的
// 运行前先设置环境变量，api_key必须设，其他随意
//
// $env:AI_API_KEY="......"
//
//export AI_BASE_URL="......"
//export AI_MODEL="......"
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// 一条对话信息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 要发给接口的请求体
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// 取出我们需要的数据
type chatResponces struct {
	Choices []struct {
		Message Message `json:"message"`
		Token   int     `json:"total_tokens"` //这个没用到
	} `json:"choices"`
}

const (
	UrlWeUse    = "https://open.bigmodel.cn/api/paas/v4"
	ModelWeUse  = "GLM-4-Flash"
	Total_token = 0
)

func ask(history []Message) (string, error) {
	apiKey := os.Getenv("Ai_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("请先设置环境变量AI_API_KEY")
	}
	baseURL := os.Getenv("AI_BASE_URL")
	if baseURL == "" {
		baseURL = UrlWeUse
	}
	model := os.Getenv("AI_MODEL")
	if model == "" {
		model = ModelWeUse
	}

	//组装请求体
	body, err := json.Marshal(ChatRequest{Model: model, Messages: history})
	if err != nil {
		return "", err
	}

	//发送Post请求
	url := strings.TrimRight(baseURL, "/") + "/chat/completions"    //将baseURl最右端的/删去
	req, err := http.NewRequest("POST", url, bytes.NewReader(body)) //将请求体转为数据流
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	client := &http.Client{Timeout: 60 * time.Second} //60秒
	resp, err := client.Do(req)                       //这里才开始发送请求
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() //返回的resp是一个结构体，且有一条TCP连接

	//读取响应并解析
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var out chatResponces
	if err := json.Unmarshal(respBytes, &out); err != nil {
		return "", fmt.Errorf("解析响应失败：%w\n原始内容：%s", err, respBytes)
	}
	return out.Choices[0].Message.Content, nil
}

func main() {
	fmt.Println("已开启对话框\nexit或者quit退出")
	var history []Message
	scanner := bufio.NewScanner(os.Stdin) //监听ing
	for {
		fmt.Print("\n至尊无敌霸王龙（不要在意细节）：")
		if !scanner.Scan() { //只有在error或者结束时才会为false
			break
		}
		question := strings.TrimSpace(scanner.Text())
		if question == "" {
			continue
		}
		if question == "exit" || question == "quit" {
			println("再见！")
			break
		}
		history = append(history, Message{Role: "user", Content: question})
		answer, err := ask(history)
		if err != nil {
			fmt.Println("报错ing...\n报错ing......\n报错ing.........")
			history = history[:len(history)-1]
			continue
		}
		history = append(history, Message{Role: "assistant", Content: answer})
		fmt.Println("Ai:" + answer)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "读取输入出错:", err) //判断是退出了还是读取出错了。
	}
}
