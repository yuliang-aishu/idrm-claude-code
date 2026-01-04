package datasource

import (
	"context"
	"database/sql"
	"fmt"
)

// NewDataSourceModel 创建数据源模型实例
// ormType: "gorm" 或 "sqlx"，默认为 gorm
func NewDataSourceModel(db *sql.DB, ormType string) DataSourceModel {
	switch ormType {
	case "sqlx":
		// TODO: 实现 SQLx 版本
		panic("SQLx version not implemented yet")
	default:
		return newGormDataSourceModel(db)
	}
}

// newGormDataSourceModel 创建 GORM 版本的数据源模型
func newGormDataSourceModel(db *sql.DB) *gormDataSourceModel {
	return &gormDataSourceModel{
		db: db,
	}
}

// gormDataSourceModel GORM 实现
type gormDataSourceModel struct {
	db *sql.DB
}

// ExecTx 执行事务
func (m *gormDataSourceModel) ExecTx(ctx context.Context, fn func(context.Context, DataSourceModel) error) error {
	// TODO: 实现事务逻辑
	return fn(ctx, m)
}

// WithTx 创建事务模型
func (m *gormDataSourceModel) WithTx(tx interface{}) DataSourceModel {
	// TODO: 实现事务模型
	return m
}

// Trans 执行事务
func (m *gormDataSourceModel) Trans(ctx context.Context, fn func(ctx context.Context, model DataSourceModel) error) error {
	// TODO: 实现事务逻辑
	return fn(ctx, m)
}

// ErrorWrap 错误包装
func (m *gormDataSourceModel) ErrorWrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	// TODO: 使用正确的错误包装
	// return errorx.Wrapf(err, "datasource model: %s", msg)
	return fmt.Errorf("datasource model: %s: %v", msg, err)
}
