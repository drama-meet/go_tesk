// 这个测试是ai写的
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLLMClientAsk(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("期望请求方法为POST，实际为%s", r.Method)
		}
		if r.URL.Path != "/v4/chat/completions" {
			t.Errorf("请求路径错误：%s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("Authorization头错误：%s", got)
		}

		var request struct {
			Model    string `json:"model"`
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("解析请求体失败：%v", err)
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		if request.Model != "test-model" || len(request.Messages) != 1 ||
			request.Messages[0].Role != "user" || request.Messages[0].Content != "你好" {
			t.Errorf("请求体内容错误：%+v", request)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"choices":[{"message":{"content":"你好！"}}]}`)
	}))
	defer server.Close()

	client := NewLLMClient(server.URL+"/v4", "test-model", "test-key", time.Second)
	answer, err := client.Ask("你好")
	if err != nil {
		t.Fatalf("Ask返回错误：%v", err)
	}
	if answer != "你好！" {
		t.Fatalf("回答错误：%q", answer)
	}
}

func TestLLMClientAskRejectsNonSuccessStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewLLMClient(server.URL, "test-model", "invalid-key", time.Second)
	if _, err := client.Ask("你好"); err == nil {
		t.Fatal("非成功HTTP状态码应返回错误")
	}
}

func TestLLMClientAskRejectsEmptyChoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, `{"choices":[]}`)
	}))
	defer server.Close()

	client := NewLLMClient(server.URL, "test-model", "test-key", time.Second)
	if _, err := client.Ask("你好"); err == nil {
		t.Fatal("choices为空时应返回错误")
	}
}
