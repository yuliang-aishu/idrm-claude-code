// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"database/sql"

	"github.com/yuliang-aishu/idrm-claude-code/spec-cc-0104/api/internal/config"
	"github.com/yuliang-aishu/idrm-claude-code/spec-cc-0104/model/datasource/datasource"
)

type ServiceContext struct {
	Config             config.Config
	DataSourceModel    datasource.DataSourceModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// TODO: 从配置中获取数据库连接
	var db *sql.DB

	return &ServiceContext{
		Config:             c,
		DataSourceModel:    datasource.NewDataSourceModel(db, "gorm"),
	}
}
