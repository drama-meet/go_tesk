// 自动化测试
// 测试文件名必须要_test.go结尾
package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// 测试list接口
func TestHandlerList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	//gin一共有3种模式，其中TestMode的输出较为干净

	store, err := NewStore(filepath.Join(t.TempDir(), "test.db"))
	//创建临时数据库
	if err != nil {
		t.Fatal("创建Store失败：", err)
	}
	defer store.Close()

	store.Create(&Todo{Title: "测试任务"})
	router := gin.New()
	//注册路由
	router.GET("/todos", (&Handler{store: store}).list)

	request := httptest.NewRequest(http.MethodGet, "/todos", nil)
	//这里请求造假，然后创建响应记录器
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("期望状态码%d，实际%d", http.StatusOK, recorder.Code)
	}

	//先将json解析后转存进todos中，再进行测试
	var todos []Todo
	if err := json.Unmarshal(recorder.Body.Bytes(), &todos); err != nil {
		t.Fatalf("解析列表响应失败：%v", err)
	}
	if len(todos) != 1 || todos[0].Title != "测试任务" {
		t.Fatalf("列表响应不符合预期：%+v", todos)
	}
	//出错就直接停下，正常的话输出pass
}

// 测试create接口
func TestHandlerCreateInvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)

	store, err := NewStore(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal("创建Store失败：", err)
	}
	defer store.Close()

	router := gin.New()
	//创建路由
	router.POST("/todos", (&Handler{store: store}).create)

	request := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(`{"title":`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("期望状态码%d，实际%d", http.StatusBadRequest, recorder.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("解析错误响应失败：%v", err)
	}
	if response["error"] == "" {
		t.Fatal("错误响应应包含错误信息")
	}
}
