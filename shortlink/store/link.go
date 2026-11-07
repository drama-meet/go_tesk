package store

import (
	"shortlink/model"
	"strings"
	"time"

	"gorm.io/gorm"
)

type LinkStore struct {
	db *gorm.DB
}

type PublicLink struct {
	ID        uint      `json:"id"`
	Code      string    `json:"code"`
	LongURL   string    `json:"long_url"`
	Clicks    int64     `json:"clicks"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
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
	until := since.AddDate(0, 0, 1)
	err := s.db.Model(&model.Link{}).
		Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, since, until).
		Count(&count).Error
	return count, err
}

// 列出用户的全部短链
func (s *LinkStore) ListByUser(userID uint) ([]model.Link, error) {
	var links []model.Link
	err := s.db.Where("user_id=?", userID).Order("created_at DESC").Find(&links).Error
	return links, err
}

// 搜索短链
func (s *LinkStore) SearchPublic(query string, popular bool, limit int) ([]PublicLink, error) {
	var links []PublicLink
	db := s.db.Table("links").
		Select("links.id, links.code, links.long_url, links.clicks, CASE WHEN links.comment_public = ? THEN links.comment ELSE '' END AS comment, links.created_at, users.name AS username", true).
		Joins("JOIN users ON users.id = links.user_id")

	//关键词匹配
	if query = strings.TrimSpace(query); query != "" {
		pattern := "%" + query + "%"
		db = db.Where("links.code LIKE ? OR links.long_url LIKE ? OR users.name LIKE ? OR (links.comment_public = ? AND links.comment LIKE ?)", pattern, pattern, pattern, true, pattern)
	}

	//根据最热或者最新排序
	if popular {
		db = db.Order("links.clicks DESC").Order("links.created_at DESC")
	} else {
		db = db.Order("links.created_at DESC")
	}

	//limit只截取前几个短链，Scan装进link中
	err := db.Limit(limit).Scan(&links).Error
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
