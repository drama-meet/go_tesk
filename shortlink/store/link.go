package store

import (
	"shortlink/model"
	"time"

	"gorm.io/gorm"
)

type LinkStore struct {
	db *gorm.DB
}

func NewLinkStore(db *gorm.DB) *LinkStore {
	return &LinkStore{db: db}
}

// 数据库存入短链
func (s *LinkStore) Create(l *model.Link) error {
	return s.db.Create(l).Error
}

// Exits判断短链是否已经被占用
func (s *LinkStore) Exits(code string) (bool, error) {
	var count int64
	err := s.db.Model(&model.Link{}).Where("code=?", code).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// 查找短链
func (s *LinkStore) FindByCode(code string) (*model.Link, error) {
	var l model.Link
	if err := s.db.Where("code=?", code).First(&l).Error; err != nil {
		return nil, err
	}
	return &l, nil
}

// 点击次数+1
func (s *LinkStore) IncrClicks(code string) error {
	return s.db.Model(&model.Link{}).Where("code=?", code).UpdateColumn("clicks", gorm.Expr("clicks+1")).Error
}

// 统计用户今天创建多少条短链了
func (s *LinkStore) CountToday(userID uint, since time.Time) (int64, error) {
	var count int64
	err := s.db.Model(&model.Link{}).Where("user_id=? AND created_at=?", userID, since).Count(&count).Error
	return count, err
}

// 列出用户的全部短链
func (s *LinkStore) ListByUser(userID uint) ([]model.Link, error) {
	var links []model.Link
	err := s.db.Where("user_id=?", userID).Order("created_at DESC").Find(&links).Error
	return links, err
}

// 按主键查的同时匹配所属用户
func (s *LinkStore) FindByIDAndUser(id uint, userID uint) (*model.Link, error) {
	var l model.Link
	if err := s.db.Where("id=? AND user_id=?", id, userID).First(&l).Error; err != nil {
		return nil, err
	}
	return &l, nil
}

// 删除短链
func (s *LinkStore) DeleteByIDAndUser(id uint, userID uint) (int64, error) {
	res := s.db.Where("id=? and user_id=?", id, userID).Delete(&model.Link{})
	return res.RowsAffected, res.Error
}
