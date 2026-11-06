package service

import (
	"errors"
	"shortlink/store"
	"time"
)

var ErrQuotaMax = errors.New("今日创建次数已达上限")

type QuotaService struct {
	links *store.LinkStore
	limit int
}

func NewQuotaService(links *store.LinkStore, limit int) *QuotaService {
	return &QuotaService{links: links, limit: limit}
}

// 检查你今天还能不能创建短链
func (s *QuotaService) Check(userID uint) error {
	today := startOfToday()

	used, err := s.links.CountToday(userID, today)
	if err != nil {
		return err
	}
	if used >= int64(s.limit) {
		return ErrQuotaMax
	}
	return nil
}

// 返回本机时区的今天的"00：00：00"
func startOfToday() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}
