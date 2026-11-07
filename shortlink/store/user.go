package store

import (
	"shortlink/model"

	"gorm.io/gorm"
)

// 负责针对user数据库的所有操作
type UserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{db: db}
}

// 创建新用户
func (s *UserStore) CreateNewUser(u *model.User) error {
	return s.db.Create(u).Error
}

// 查询用户
func (s *UserStore) FindUser(username string) (*model.User, error) {
	var u model.User
	if err := s.db.Where("name=?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// 判断用户名称是否已经存在
func (s *UserStore) Exits(username string) (bool, error) {
	var count int64 //gorm.Count()只支持int64
	err := s.db.Model(&model.User{}).Where("name=?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// 判断用户邮箱是否已经注册
func (s *UserStore) EmailExits(email string) (bool, error) {
	var count int64 //gorm.Count()只支持int64
	err := s.db.Model(&model.User{}).Where("email=?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
