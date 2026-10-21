// 写model.go里有点乱，所以改到这里面了
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// 将大模型相关的信息进行封装然后存接口中
type LLMClient struct {
	baseURL string
	model   string
	apiKey  string
	Client  *http.Client //可设置超时时间
}

func NewLLMClient(baseURL, model, apiKey string, timeout time.Duration) *LLMClient {
	return &LLMClient{
		baseURL: baseURL,
		model:   model,
		apiKey:  apiKey,
		Client:  &http.Client{Timeout: timeout},
	}
}
func (l *LLMClient) Ask(question string) (string, error) {
	chat, err := json.Marshal(chatRequest{
		Model:    l.model,
		Messages: []chatMessage{{Role: "user", Content: question}},
	})
	if err != nil {
		return "", fmt.Errorf("构造请求体失败：%w", err)
	}
	url := strings.TrimRight(l.baseURL, "/") + "/chat/completions"  //将baseURl最右端的/删去
	req, err := http.NewRequest("POST", url, bytes.NewReader(chat)) //将请求体转为数据流
	if err != nil {
		return "", fmt.Errorf("构造大模型请求失败：%w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+l.apiKey) //通用的api_key的请求

	//发送构造好的请求体
	resp, err := l.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求大模型失败：%w", err)
	}
	defer resp.Body.Close() //存有TCP连接

	//读取响应
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("解析响应失败:%w", err)
	}
	//下面3行ai写的
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("大模型返回非成功状态码：%s", resp.Status)
	}

	var out chatResponse
	if err := json.Unmarshal(respBytes, &out); err != nil {
		return "", fmt.Errorf("解析大模型响应失败：%w", err)
	}
	if len(out.Choices) == 0 {
		return "", fmt.Errorf("大模型响应中没有可用的回答")
	}
	return out.Choices[0].Message.Content, nil
}
