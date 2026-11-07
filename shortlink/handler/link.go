package handler

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"shortlink/middle"
	"shortlink/model"
	"shortlink/pkg"
	"shortlink/service"
	"shortlink/store"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LinkHandler struct {
	links *store.LinkStore
	quota *service.QuotaService
}

func NewLinkHandler(links *store.LinkStore, quota *service.QuotaService) *LinkHandler {
	return &LinkHandler{links: links, quota: quota}
}

type creasteReq struct {
	URL           string `json:"url" binding:"required"`
	Comment       string `json:"comment" binding:"max=500"`
	CommentPublic bool   `json:"comment_public"`
}

func (h *LinkHandler) Create(c *gin.Context) {
	var req creasteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供url，注释不能超过500字"})
		return
	}

	//	只允许http/https，且必须有域名
	u, err := url.ParseRequestURI(req.URL)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url必须是合法的http/https地址"})
		return
	}

	//取出当前用户id
	uid, ok := middle.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	//检查配额
	if err := h.quota.Check(uid); err != nil {
		if errors.Is(err, service.ErrQuotaMax) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}

	//生成短链ing
	code, err := h.CreateUniqueCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "短码生成失败，请重新尝试"})
		return
	}

	//入库
	l := &model.Link{
		UserID:        uid,
		Code:          code,
		LongUrl:       req.URL,
		Comment:       req.Comment,
		CommentPublic: req.CommentPublic,
	}
	if err := h.links.Create(l); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}

	//拼短地址
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	shortURL := scheme + "://" + c.Request.Host + "/s/" + code
	c.JSON(http.StatusCreated, gin.H{"code": code, "short_url": shortURL})
}

// 最多尝试创建11次
func (h *LinkHandler) CreateUniqueCode() (string, error) {
	for i := 0; i <= 10; i++ {
		code, err := pkg.CreateShortCode(6)
		if err != nil {
			return "", err
		}
		exits, err := h.links.Exits(code)
		if err != nil {
			return "", err
		}
		if !exits {
			return code, nil
		}
	}
	return "", errors.New("短码连续冲突")
}

// 重定向
func (h *LinkHandler) Redirect(c *gin.Context) {
	code := c.Param("code")

	//查询短链
	l, err := h.links.FindByCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "短链不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}

	//302临时重定向
	c.Redirect(http.StatusFound, l.LongUrl)

	//跳转失败不计入次数
	if err := h.links.IncrClicks(code); err != nil {
		log.Printf("点击计数失败 code=%s:%v", code, err)
	}
}

// 列出当前用户的所有短链
func (h *LinkHandler) List(c *gin.Context) {
	uid, ok := middle.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	links, err := h.links.ListByUser(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": links, "total": len(links)})
}

// 搜索和热门短链的接口
func (h *LinkHandler) PublicList(c *gin.Context) {
	limit := 100
	if c.Query("sort") == "popular" {
		limit = 10
	}
	links, err := h.links.SearchPublic(c.Query("q"), c.Query("sort") == "popular", limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": links, "total": len(links)})
}

// 查看自己的单条短链
func (h *LinkHandler) Detail(c *gin.Context) {
	uid, ok := middle.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id不合法"})
		return
	}
	l, err := h.links.FindByIDAndUser(uint(id), uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "短链不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}
	c.JSON(http.StatusOK, l)
}

// Delete删除单链表
func (h *LinkHandler) Delete(c *gin.Context) {
	uid, ok := middle.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id不合法"})
		return
	}

	rows, err := h.links.DeleteByIDAndUser(uint(id), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "短链不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": "删除成功"})
}
