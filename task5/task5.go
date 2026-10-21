package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Httpsnote struct {
	Test    string `json:"hitokoto"`
	From    string `json:"from"`
	Id      int    `json:"id"`
	Creator string `json:"creator"`
}

func Get(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("请求失败:", err)
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}
	var note Httpsnote
	err = json.Unmarshal(body, &note)
	if err != nil {
		return "", fmt.Errorf("解析JSON失败: %v", err)
	}
	vvv := "                "
	return fmt.Sprintf("%d.\n[%s]\n%s--%s\n%s(%s)", note.Id, note.Test, vvv, note.From, vvv, note.Creator), nil
}

func main() {
	url := "https://v1.hitokoto.cn/"
	test, err := Get(url)
	if err != nil {
		fmt.Println("获取数据失败:", err)
		return
	}
	fmt.Print(test)
}
