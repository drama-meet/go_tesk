package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func (t Todo) Require() error {
	if strings.TrimSpace(t.Title) == "" {
		return errors.New("title不能为空")
	}
	return nil
}

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("todos.db"), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败：" + err.Error())
	}
	db.AutoMigrate(&Todo{})

	r := gin.Default()
	r.GET("/todos", listTodos)
	r.GET("/todos/:id", getTodo)
	r.POST("/todos", createTodo)
	r.POST("/todos/:id", changeTodo)
	r.DELETE("/todos/:id", deleteTodo)

	r.Run(":8080")

}

func listTodos(c *gin.Context) {
	var todos []Todo
	if DoneStr := c.Query("done"); DoneStr != "" {
		done, err := strconv.ParseBool(DoneStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "done只支持bool值"})
			return
		}
		db.Where("done=?", done).Find(&todos)
	} else {
		db.Find(&todos)
	}
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
	db.Create(&input)
	c.JSON(http.StatusCreated, input)
}

func deleteTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "id必须是数字"})
		return
	}
	result := db.Unscoped().Delete(&Todo{}, id)
	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "没找到这条待办"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已删除"})
}

func getTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "id必须是数字"})
		return
	}
	var outer Todo
	if err := db.First(&outer, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		}
		return
	}
	c.JSON(http.StatusOK, outer)
}

func changeTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "id必须是数字"})
		return
	}
	var changing Todo
	if err := db.First(&changing, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		}
		return
	}
	var input struct {
		Title string `json:"title"`
		Done  bool   `json:"done"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求体格式不正确"})
		return
	}
	if strings.TrimSpace(input.Title) == "" {
		c.JSON(400, gin.H{"error": "title不能为空"})
		return
	}
	changing.Title = input.Title
	changing.Done = input.Done
	if err := db.Save(&changing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "修改失败"})
		return
	}
	c.JSON(http.StatusOK, changing)
}
