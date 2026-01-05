package datasource

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// NewDataSourceModel 创建数据源模型实例
// ormType: "gorm" 或 "sqlx"，默认为 gorm
func NewDataSourceModel(db *gorm.DB, ormType string) DataSourceModel {
	switch ormType {
	case "sqlx":
		// TODO: 实现 SQLx 版本
		panic("SQLx version not implemented yet")
	default:
		return newGormDataSourceModel(db)
	}
}

// newGormDataSourceModel 创建 GORM 版本的数据源模型
func newGormDataSourceModel(db *gorm.DB) *gormDataSourceModel {
	return &gormDataSourceModel{
		db: db,
	}
}

// gormDataSourceModel GORM 实现
type gormDataSourceModel struct {
	db *gorm.DB
}

// ExecTx 执行事务
func (m *gormDataSourceModel) ExecTx(ctx context.Context, fn func(context.Context, DataSourceModel) error) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		model := &gormDataSourceModel{db: tx}
		return fn(ctx, model)
	})
}

// WithTx 创建事务模型
func (m *gormDataSourceModel) WithTx(tx interface{}) DataSourceModel {
	if dbTx, ok := tx.(*gorm.DB); ok {
		return &gormDataSourceModel{db: dbTx}
	}
	return m
}

// Trans 执行事务
func (m *gormDataSourceModel) Trans(ctx context.Context, fn func(ctx context.Context, model DataSourceModel) error) error {
	return m.ExecTx(ctx, fn)
}

// ErrorWrap 错误包装
func (m *gormDataSourceModel) ErrorWrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("datasource model: %s: %v", msg, err)
}
