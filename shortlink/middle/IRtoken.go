package middle

import (
	"net/http"
	"shortlink/pkg"
	"strings"

	"github.com/gin-gonic/gin"
)

const UID = "userID"

// 判定是否有合法的token
func Middletoken(jwtM *pkg.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		//取Authorization头
		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少token"})
			c.Abort() //终止中间件
			return
		}

		//必须是"Bearer <token>"格式
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token格式错误"})
			c.Abort()
			return
		}

		//解析并校验token
		claims, err := jwtM.IsTokenRight(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token无效或已过期"})
			c.Abort()
			return
		}

		//将用户id放入context
		c.Set(UID, claims.UserId)

		//放行
		c.Next()
	}
}

// 取这个用来登录
func GetUserID(c *gin.Context) (uint, bool) {
	v, exits := c.Get(UID)
	if !exits {
		return 0, false
	}
	uid, ok := v.(uint)
	return uid, ok
}
