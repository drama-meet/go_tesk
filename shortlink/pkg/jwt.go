package pkg

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 这个会塞进token中
type Claims struct {
	UserId               uint `json:"user_id"`
	jwt.RegisteredClaims      //jwt包中自带的结构体
}

// 负责签发和验收token
type JWTManager struct {
	secret []byte        //密钥
	ttl    time.Duration //有效期
}

func NewJWTManager(secret string, ttl time.Duration) *JWTManager {
	return &JWTManager{secret: []byte(secret), ttl: ttl}
}

// 给用户签发token
func (m *JWTManager) GiveToken(userId uint) (string, error) {
	claims := Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //选择jwt的算法进行盖章
	return token.SignedString(m.secret)
}

// 验证token的合法性并取出数据
func (m *JWTManager) IsTokenRight(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return m.secret, nil //
		}, //拆开token之后再回调查看密钥，能看懂功能，但是看不懂函数
		jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("无效token") //jwt自带的判定方式
	}
	return claims, nil
}
