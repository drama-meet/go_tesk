package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func (t Todo) Require() error {
	if strings.TrimSpace(t.Title) == "" {
		return errors.New("title不能为空")
	}
	return nil
}

var (
	todos  []Todo
	nextId = 1
	mu     sync.Mutex
)

func main() {
	r := gin.Default()
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatal(err)
	}
	r.GET("/todos", listTodos)
	r.GET("/todos/:id", getId)
	r.POST("/todos", createTodo)
	r.POST("/todos/:id", changeTodo)
	r.DELETE("/todos/:id", deleteTodo)

	r.Run(":11111")

}

func listTodos(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	c.JSON(http.StatusOK, todos)
}

func createTodo(c *gin.Context) {
	var input Todo
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := input.Require(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mu.Lock()
	defer mu.Unlock()
	input.Id = nextId
	nextId++
	todos = append(todos, input)
	c.JSON(http.StatusCreated, input)
}

func deleteTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id必须是数字"})
		return
	}
	mu.Lock()
	defer mu.Unlock()
	for i, t := range todos {
		if t.Id == id {
			todos = append(todos[:i], todos[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "已删除"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "没找到这条待办"})
}

func getId(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id必须是数字"})
		return
	}
	mu.Lock()
	defer mu.Unlock()
	for _, t := range todos {
		if t.Id == id {
			c.JSON(http.StatusOK, t)
			return
		}
	}
	c.JSON(404, gin.H{"error": "Not Found"})
}

func changeTodo(c *gin.Context) {
	var input Todo
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := input.Require(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//验证title是否为空
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id必须是数字"})
		return
	}
	//验证id
	mu.Lock()
	defer mu.Unlock()
	for i, t := range todos {
		if t.Id == id {
			todos[i].Done = input.Done
			todos[i].Title = input.Title
			c.JSON(http.StatusOK, gin.H{"message": "修改成功"})
			return
		}
	}
	c.JSON(404, gin.H{"error": "Not Found"})
}
