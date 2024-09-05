package storage

import (
	"context"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TorStorage struct {
	db *gorm.DB
}

type TorExitNode struct {
	IP          string    `gorm:"primaryKey"`
	LastSeen    time.Time `gorm:"primaryKey"`
	FetchedFrom string    `gorm:"primaryKey"`
}

type AllowListIP struct {
	IP        string `gorm:"primaryKey"`
	Note      string
	TimeAdded time.Time
}

func New(sqlitePath string) (*TorStorage, error) {
	db, err := gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&TorExitNode{})
	db.AutoMigrate(&AllowListIP{})
	return &TorStorage{db: db}, nil
}

func (s *TorStorage) InsertBatch(ctx context.Context, nodes []TorExitNode) {
	s.db.WithContext(ctx).Create(nodes)
}

func (s *TorStorage) InsertAllowIP(entry AllowListIP) {
	s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "ip"}},
		DoUpdates: clause.AssignmentColumns([]string{"note", "time_added"}),
	}).Create(&entry)
}

func (s *TorStorage) DeleteAllowIP(IP string) {
	s.db.Delete(&AllowListIP{IP: IP})
}

func (s *TorStorage) GetAllowListIPs() []AllowListIP {
	list := []AllowListIP{}
	s.db.Model(&list).Find(&list)
	return list
}

func (s *TorStorage) GetTorNodes() []TorExitNode {
	sub := s.db.Model(AllowListIP{}).Select("ip")
	res := []TorExitNode{}
	s.db.Model(TorExitNode{}).Where("ip NOT IN (?)", sub).Find(&res)
	return res
}
