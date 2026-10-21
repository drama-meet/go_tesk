package main

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type Store struct {
	db *gorm.DB
}

func NewStore(dbPath string) (*Store, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	//连接数据库并设置配置选项（这里为默认）
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&Todo{}); err != nil {
		//上传结构体并自动建表
		return nil, err
	}
	return &Store{db: db}, nil
	//将连接数据传入结构体并返回
}

func (s *Store) Close() {
	if sqlDB, err := s.db.DB(); err == nil {
		//Close()需要通过DB()来拿到
		sqlDB.Close()
		//关闭连接不用上抛err,一般没人接
	}
}

// 这里没有err上抛，在handle里面
func (s *Store) List() []Todo {
	var todos []Todo
	s.db.Find(&todos)
	return todos
}

func (s *Store) Create(t *Todo) {
	s.db.Create(t)
}

func (s *Store) Delete(id int) bool {
	result := s.db.Delete(&Todo{}, id)
	return result.RowsAffected > 0
}
