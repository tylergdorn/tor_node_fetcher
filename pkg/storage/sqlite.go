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

func (s *TorStorage) InsertBatch(ctx context.Context, nodes []TorExitNode) error {
	err := s.db.WithContext(ctx).Create(nodes).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *TorStorage) InsertAllowIP(ctx context.Context, entry AllowListIP) error {
	err := s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "ip"}},
		DoUpdates: clause.AssignmentColumns([]string{"note", "time_added"}),
	}).Create(&entry).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *TorStorage) DeleteAllowIP(ctx context.Context, IP string) error {
	if err := s.db.WithContext(ctx).Delete(&AllowListIP{IP: IP}).Error; err != nil {
		return err
	}
	return nil
}

func (s *TorStorage) GetAllowListIPs(ctx context.Context) ([]AllowListIP, error) {
	list := []AllowListIP{}
	if err := s.db.WithContext(ctx).Model(&list).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *TorStorage) GetTorNodesPaginated(ctx context.Context, limit int, offset int) ([]TorExitNode, error) {
	sub := s.db.WithContext(ctx).Model(AllowListIP{}).Select("ip")
	if sub.Error != nil {
		return nil, sub.Error
	}
	res := []TorExitNode{}
	if err := s.db.Model(TorExitNode{}).Where("ip NOT IN (?)", sub).Limit(limit).Offset(offset).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (s *TorStorage) GetTorNodes(ctx context.Context) ([]TorExitNode, error) {
	return s.GetTorNodesPaginated(ctx, -1, -1)
}
