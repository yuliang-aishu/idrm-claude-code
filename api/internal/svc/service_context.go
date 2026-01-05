// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"fmt"

	"github.com/yuliang-aishu/idrm-claude-code/spec-cc-0104/api/internal/config"
	"github.com/yuliang-aishu/idrm-claude-code/spec-cc-0104/model/datasource/datasource"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ServiceContext struct {
	Config             config.Config
	DataSourceModel    datasource.DataSourceModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 GORM 数据库连接
	// 这里需要从环境变量或配置文件中获取数据库连接信息
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local",
		c.DB.Default.Username, c.DB.Default.Password, c.DB.Default.Host, c.DB.Default.Port, c.DB.Default.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 关闭 GORM 日志
		NamingStrategy: nil, // 使用默认命名策略
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("连接数据库失败: " + err.Error())
	}

	return &ServiceContext{
		Config:             c,
		DataSourceModel:    datasource.NewDataSourceModel(db, "gorm"),
	}
}
