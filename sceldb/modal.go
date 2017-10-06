package sceldb

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Cache struct {
	gorm.Model

	Key    string `gorm:"not null;unique"`
	Text   string `gorm:"txt"`
	Binary []byte `gorm:"bin"`
	Size   int    `gorm:"size"`
}

type Dict struct {
	gorm.Model

	Name           string    `gorm:"name"`
	Count          int       `gorm:"count"`
	Creator        string    `gorm:"creator"`
	Size           int       `gorm:"size"`           // 网页中显示的大小
	DictUpdateTime time.Time `gorm:"not null;index"` // 网页中显示的更新时间
	Version        int       `gorm:"version"`        // 网页中显示的版本号
	Summary        string    `gorm:"summary"`
	DownloadCount  int       `gorm:"download_num"`
	DownloadUrl    string    `gorm:"download_url"`
	CheckAt        time.Time `gorm:"not null;index"` // 检测时间
}
