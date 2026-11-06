package handler

import (
	"net/http"
	"shortlink/model"
	"shortlink/pkg"
	"shortlink/store"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	user *store.UserStore
	jwt  *pkg.JWTManager
}

func NewUserHandler(user *store.UserStore, jwt *pkg.JWTManager) *UserHandler {
	return &UserHandler{user: user, jwt: jwt}
}

// binding标签可以让gin触发自动校验
type UPlength struct {
	Username string `json:"username" binding:"required,max=32"`
	Password string `json:"password" binding:"required,max=255"`
}

// 创建用户
func (h UserHandler) UserRH(c *gin.Context) {
	var req UPlength
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名最多32位，密码最多255位"})
		return
	}

	//确保用户不存在
	exits, err := h.user.Exits(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误(1)"})
		return
	}
	if exits {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名被占用"})
		return
	}

	//密码加密
	hash, err := pkg.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误(2)"})
		return
	}

	//将数据存入库中
	u := &model.User{
		Name:     req.Username,
		Password: hash,
	}
	if err := h.user.CreateNewUser(u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误(3)"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "注册成功"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req UPlength
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名或密码超出服务范围"})
		return
	}

	//确保用户不存在
	u, err := h.user.FindUser(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名密码错误"})
		return
	}

	//验证密码
	if !pkg.CheckPassword(u.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	//签发token
	token, err := h.jwt.GiveToken(u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误(4)"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
