// pkg一般用来存放通用、与业务无关的工具
package pkg

import "golang.org/x/crypto/bcrypt"

// 将明文密码转成哈希值
func HashPassword(pwd string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost) //生成hash,返回[]byte;cost表示的是计算强度
	if err != nil {
		return "", err
	}
	return string(b), nil
}

//对比明文密码与数据库中存储的hash值是否匹配
func CheckPassword(hash string, pwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd)) == nil
}
