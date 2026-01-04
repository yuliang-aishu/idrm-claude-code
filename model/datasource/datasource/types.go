package datasource

import (
	"time"

	"gorm.io/gorm"
)

// DataSource 数据源实体
type DataSource struct {
	Id          string         `gorm:"primaryKey;size:36"`          // UUID v7
	Name        string         `gorm:"size:100;not null;uniqueIndex"`
	Type        string         `gorm:"size:50;not null"`             // mysql/postgresql/redis/mongodb/sqlserver
	Host        string         `gorm:"size:200;not null"`
	Port        int            `gorm:"not null"`
	Database    string         `gorm:"size:100"`
	Username    string         `gorm:"size:100;not null"`
	Password    string         `gorm:"size:500;not null"`            // 加密存储
	Description string         `gorm:"size:500"`
	Status      string         `gorm:"size:20;not null;default:'enabled'"`
	SortOrder   int            `gorm:"not null;default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// DataSourceResp 响应结构（隐藏敏感信息）
type DataSourceResp struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Database    string `json:"database"`
	Description string `json:"description"`
	Status      string `json:"status"`
	SortOrder   int    `json:"sort_order"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// DataSourceQuery 查询参数
type DataSourceQuery struct {
	Offset    int    `form:"offset,default=1"`
	Limit     int    `form:"limit,default=10"`
	Keyword   string `form:"keyword,optional"`
	Status    string `form:"status,optional"`
	Sort      string `form:"sort,default=created_at"`
	Direction string `form:"direction,default=desc"`
}
