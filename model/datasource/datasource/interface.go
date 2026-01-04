package datasource

import "context"

// DataSourceModel 数据源模型接口
type DataSourceModel interface {
	// 插入数据源
	Insert(ctx context.Context, data *DataSource) (*DataSource, error)

	// 根据 ID 查询数据源
	FindOne(ctx context.Context, id string) (*DataSource, error)

	// 更新数据源
	Update(ctx context.Context, data *DataSource) error

	// 软删除数据源
	Delete(ctx context.Context, id string) error

	// 列表查询（支持分页、搜索、筛选）
	FindList(ctx context.Context, query *DataSourceQuery) ([]*DataSource, int64, error)

	// 检查名称是否重复
	CheckNameExists(ctx context.Context, name string, excludeId ...string) (bool, error)

	// 连接测试
	TestConnection(ctx context.Context, config *DataSource) error

	// 事务支持
	WithTx(tx interface{}) DataSourceModel
	Trans(ctx context.Context, fn func(ctx context.Context, model DataSourceModel) error) error
}
