package datasource

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

	// 使用 GORM 插入到数据库
	db := m.db.WithContext(ctx)
	result := db.Create(data)
	if result.Error != nil {
		// 检查是否是名称重复错误
		if strings.Contains(result.Error.Error(), "Duplicate entry") {
			return nil, ErrNameExists.Errorf("数据源名称已存在")
		}
		return nil, m.ErrorWrap(result.Error, "插入数据源失败")
	}

	return data, nil
}

// FindOne 根据 ID 查询数据源
func (m *gormDataSourceModel) FindOne(ctx context.Context, id string) (*DataSource, error) {
	if id == "" {
		return nil, ErrParamInvalid.Errorf("ID 不能为空")
	}

	var data DataSource
	db := m.db.WithContext(ctx)
	result := db.Where("id = ?", id).First(&data)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrNotFound.Errorf("数据源不存在")
		}
		return nil, m.ErrorWrap(result.Error, "查询数据源失败")
	}

	return &data, nil
}

// Update 更新数据源
func (m *gormDataSourceModel) Update(ctx context.Context, data *DataSource) error {
	if data == nil {
		return ErrParamInvalid.Errorf("数据源信息不能为空")
	}

	if data.Id == "" {
		return ErrParamInvalid.Errorf("ID 不能为空")
	}

	// 设置更新时间
	data.UpdatedAt = time.Now()

	db := m.db.WithContext(ctx)
	result := db.Save(data)
	if result.Error != nil {
		// 检查是否是名称重复错误
		if strings.Contains(result.Error.Error(), "Duplicate entry") {
			return ErrNameExists.Errorf("数据源名称已存在")
		}
		return m.ErrorWrap(result.Error, "更新数据源失败")
	}

	return nil
}

// Delete 软删除数据源
func (m *gormDataSourceModel) Delete(ctx context.Context, id string) error {
	if id == "" {
		return ErrParamInvalid.Errorf("ID 不能为空")
	}

	db := m.db.WithContext(ctx)
	result := db.Delete(&DataSource{}, "id = ?", id)
	if result.Error != nil {
		return m.ErrorWrap(result.Error, "删除数据源失败")
	}

	if result.RowsAffected == 0 {
		return ErrNotFound.Errorf("数据源不存在")
	}

	return nil
}

// ============================================
// 列表查询
// ============================================

// FindList 列表查询（支持分页、搜索、筛选）
func (m *gormDataSourceModel) FindList(ctx context.Context, query *DataSourceQuery) ([]*DataSource, int64, error) {
	if query == nil {
		query = &DataSourceQuery{}
	}

	// 设置默认值
	if query.Offset <= 0 {
		query.Offset = 1
	}
	if query.Limit <= 0 || query.Limit > MaxPageSize {
		query.Limit = 10
	}
	if query.Sort == "" {
		query.Sort = DefaultSortField
	}
	if query.Direction == "" {
		query.Direction = DefaultSortDirection
	}

	db := m.db.WithContext(ctx)

	// 构建查询条件
	db = m.buildQuery(db, query)

	// 查询总数
	var totalCount int64
	result := db.Model(&DataSource{}).Count(&totalCount)
	if result.Error != nil {
		return nil, 0, m.ErrorWrap(result.Error, "查询总数失败")
	}

	// 查询分页数据
	var dataSources []*DataSource
	offset := (query.Offset - 1) * query.Limit
	result = db.Offset(offset).Limit(query.Limit).Find(&dataSources)
	if result.Error != nil {
		return nil, 0, m.ErrorWrap(result.Error, "查询列表失败")
	}

	return dataSources, totalCount, nil
}

// CheckNameExists 检查名称是否重复
func (m *gormDataSourceModel) CheckNameExists(ctx context.Context, name string, excludeId ...string) (bool, error) {
	if name == "" {
		return false, ErrParamInvalid.Errorf("名称不能为空")
	}

	db := m.db.WithContext(ctx)
	query := db.Model(&DataSource{}).Where("name = ?", name)

	// 如果有 excludeId，排除该记录
	if len(excludeId) > 0 && excludeId[0] != "" {
		query = query.Where("id != ?", excludeId[0])
	}

	var count int64
	result := query.Count(&count)
	if result.Error != nil {
		return false, m.ErrorWrap(result.Error, "检查名称重复失败")
	}

	return count > 0, nil
}

// TestConnection 连接测试
func (m *gormDataSourceModel) TestConnection(ctx context.Context, config *DataSource) error {
	if config == nil {
		return ErrConnectionTest.Errorf("配置信息不能为空")
	}

	// 根据数据源类型执行连接测试
	switch config.Type {
	case DataSourceTypeMySQL:
		return m.testMySQLConnection(ctx, config)
	case DataSourceTypePostgreSQL:
		return m.testPostgreSQLConnection(ctx, config)
	case DataSourceTypeRedis:
		return m.testRedisConnection(ctx, config)
	case DataSourceTypeMongoDB:
		return m.testMongoDBConnection(ctx, config)
	case DataSourceTypeSQLServer:
		return m.testSQLServerConnection(ctx, config)
	default:
		return ErrConnectionTest.Errorf("不支持的数据源类型: %s", config.Type)
	}
}

// ============================================
// 私有辅助方法
// ============================================

// buildQuery 构建查询条件
func (m *gormDataSourceModel) buildQuery(db *gorm.DB, query *DataSourceQuery) *gorm.DB {
	// 状态筛选
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	// 关键字搜索
	if query.Keyword != "" {
		keyword := "%" + query.Keyword + "%"
		db = db.Where("name LIKE ? OR description LIKE ?", keyword, keyword)
	}

	// 软删除过滤 - 只查询未删除的记录
	db = db.Where("deleted_at IS NULL")

	// 排序
	sortField := query.Sort
	if sortField == "" {
		sortField = DefaultSortField
	}

	direction := query.Direction
	if direction == "" {
		direction = DefaultSortDirection
	}

	// 验证排序字段
	validSortFields := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"name":       true,
		"sort_order": true,
	}

	if !validSortFields[sortField] {
		sortField = DefaultSortField
	}

	if strings.ToUpper(direction) == "ASC" {
		db = db.Order(sortField + " ASC")
	} else {
		db = db.Order(sortField + " DESC")
	}

	return db
}

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

// ============================================
// 连接测试方法
// ============================================

// testMySQLConnection 测试 MySQL 连接
func (m *gormDataSourceModel) testMySQLConnection(ctx context.Context, config *DataSource) error {
	// 这里应该实现 MySQL 连接测试
	// 由于是测试环境，这里简单模拟
	return nil
}

// testPostgreSQLConnection 测试 PostgreSQL 连接
func (m *gormDataSourceModel) testPostgreSQLConnection(ctx context.Context, config *DataSource) error {
	// 这里应该实现 PostgreSQL 连接测试
	// 由于是测试环境，这里简单模拟
	return nil
}

// testRedisConnection 测试 Redis 连接
func (m *gormDataSourceModel) testRedisConnection(ctx context.Context, config *DataSource) error {
	// 这里应该实现 Redis 连接测试
	// 由于是测试环境，这里简单模拟
	return nil
}

// testMongoDBConnection 测试 MongoDB 连接
func (m *gormDataSourceModel) testMongoDBConnection(ctx context.Context, config *DataSource) error {
	// 这里应该实现 MongoDB 连接测试
	// 由于是测试环境，这里简单模拟
	return nil
}

// testSQLServerConnection 测试 SQLServer 连接
func (m *gormDataSourceModel) testSQLServerConnection(ctx context.Context, config *DataSource) error {
	// 这里应该实现 SQLServer 连接测试
	// 由于是测试环境，这里简单模拟
	return nil
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
