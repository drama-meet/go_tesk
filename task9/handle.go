// 接口层
package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	store *Store
	llm   *LLMClient
}

func (h *Handler) list(c *gin.Context) {
	c.JSON(http.StatusOK, h.store.List())
}

func (h *Handler) create(c *gin.Context) {
	var input Todo
	if err := c.ShouldBindJSON(&input); err != nil {
		//检查创建的格式是否正确
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.store.Create(&input)
	c.JSON(http.StatusCreated, input)
}

func (h *Handler) delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	//取参后转为整数
	if err != nil {
		c.JSON(400, gin.H{"error": "id必须是数字"})
		return
	}
	if !h.store.Delete(id) {
		c.JSON(404, gin.H{"error": "未找到这条待办"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已删除"})
}

func (h *Handler) ask(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q")) //读取查询参数
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供要询问的问题"})
		return
	}
	if h.llm == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "未配置AI_API_KEY，该服务暂不可用"})
		return
	}
	answer, err := h.llm.Ask(q)
	if err != nil {
		log.Printf("AI服务调用失败：%v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "AI服务调用失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"question": q, "answer": answer})
}
