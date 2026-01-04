package datasource

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ============================================
// 基础 CRUD 操作
// ============================================

// Insert 插入数据源
func (m *gormDataSourceModel) Insert(ctx context.Context, data *DataSource) (*DataSource, error) {
	if data == nil {
		return nil, ErrParamInvalid.Errorf("数据源信息不能为空")
	}

	// 生成 UUID v7 主键（如果未设置）
	if data.Id == "" {
		data.Id = GenerateUUIDv7()
	}

	// 设置创建和更新时间
	now := time.Now()
	data.CreatedAt = now
	data.UpdatedAt = now

	// TODO: 使用 GORM 实现插入到数据库
	// 示例代码：
	// db := m.db.WithContext(ctx)
	// result := db.Create(data)
	// if result.Error != nil {
	//     return nil, m.ErrorWrap(result.Error, "插入数据源失败")
	// }

	return data, nil
}

// FindOne 根据 ID 查询数据源
func (m *gormDataSourceModel) FindOne(ctx context.Context, id string) (*DataSource, error) {
	// TODO: 使用 GORM 实现查询
	// - 根据 ID 查询
	// - 软删除过滤
	return nil, fmt.Errorf("not implemented")
}

// Update 更新数据源
func (m *gormDataSourceModel) Update(ctx context.Context, data *DataSource) error {
	// TODO: 使用 GORM 实现更新
	// - 更新记录
	// - 软删除过滤
	return fmt.Errorf("not implemented")
}

// Delete 软删除数据源
func (m *gormDataSourceModel) Delete(ctx context.Context, id string) error {
	// TODO: 使用 GORM 实现软删除
	// - 设置 deleted_at 字段
	return fmt.Errorf("not implemented")
}

// ============================================
// 列表查询
// ============================================

// FindList 列表查询（支持分页、搜索、筛选）
func (m *gormDataSourceModel) FindList(ctx context.Context, query *DataSourceQuery) ([]*DataSource, int64, error) {
	// TODO: 使用 GORM 实现列表查询
	// - 支持分页（offset, limit）
	// - 支持关键字搜索（name, description 模糊匹配）
	// - 支持状态筛选
	// - 支持排序（sort, direction）
	// - 软删除过滤
	// - 返回总数和分页数据
	return nil, 0, fmt.Errorf("not implemented")
}

// CheckNameExists 检查名称是否重复
func (m *gormDataSourceModel) CheckNameExists(ctx context.Context, name string, excludeId ...string) (bool, error) {
	// TODO: 使用 GORM 实现名称重复检查
	// - 查询相同名称的记录
	// - 如果有 excludeId，排除该记录
	// - 返回是否存在
	return false, fmt.Errorf("not implemented")
}

// TestConnection 连接测试
func (m *gormDataSourceModel) TestConnection(ctx context.Context, config *DataSource) error {
	// TODO: 实现连接测试
	// - 根据类型测试连接
	// - 支持 MySQL/PostgreSQL/Redis/MongoDB/SQLServer
	// - 返回连接结果
	return fmt.Errorf("not implemented")
}

// ============================================
// 私有辅助方法
// ============================================

// buildQuery 构建查询条件
// func (m *gormDataSourceModel) buildQuery(db *gorm.DB, query *DataSourceQuery) *gorm.DB {
	// TODO: 实现查询条件构建
	// 分页
	// if query.Limit > 0 {
	// 	db = db.Limit(query.Limit)
	// }
	// if query.Offset > 0 {
	// 	db = db.Offset((query.Offset - 1) * query.Limit)
	// }

	// 状态筛选
	// if query.Status != "" {
	// 	db = db.Where("status = ?", query.Status)
	// }

	// 关键字搜索
	// if query.Keyword != "" {
	// 	keyword := "%" + query.Keyword + "%"
	// 	db = db.Where("name LIKE ? OR description LIKE ?", keyword, keyword)
	// }

	// 排序
	// sortField := query.Sort
	// if sortField == "" {
	// 	sortField = DefaultSortField
	// }

	// direction := query.Direction
	// if direction == "" {
	// 	direction = DefaultSortDirection
	// }

	// 验证排序字段
	// validSortFields := map[string]bool{
	// 	"created_at": true,
	// 	"updated_at": true,
	// 	"name":       true,
	// 	"sort_order": true,
	// }

	// if !validSortFields[sortField] {
	// 	sortField = DefaultSortField
	// }

	// if strings.ToUpper(direction) == "ASC" {
	// 	db = db.Order(sortField + " ASC")
	// } else {
	// 	db = db.Order(sortField + " DESC")
	// }

	// 软删除过滤
	// db = db.Where("deleted_at IS NULL") // TODO: 启用软删除

	// return db
// }

// toResp 转换为响应结构
func (m *gormDataSourceModel) toResp(data *DataSource) *DataSourceResp {
	if data == nil {
		return nil
	}

	return &DataSourceResp{
		Id:          data.Id,
		Name:        data.Name,
		Type:        data.Type,
		Host:        data.Host,
		Port:        data.Port,
		Database:    data.Database,
		Description: data.Description,
		Status:      data.Status,
		SortOrder:   data.SortOrder,
		CreatedAt:   data.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   data.UpdatedAt.Format(time.RFC3339),
	}
}

// GenerateUUIDv7 生成 UUID v7
func GenerateUUIDv7() string {
	id, err := uuid.NewV7()
	if err != nil {
		// 如果 UUID v7 不可用，使用 v4
		id, err = NewV4()
		if err != nil {
			// 如果都失败，使用时间戳
			return fmt.Sprintf("%016d", time.Now().UnixNano())
		}
	}
	return id.String()
}

// NewV4 生成 UUID v4
func NewV4() (uuid.UUID, error) {
	return uuid.NewUUID()
}
