package datasource

import "fmt"

// ErrorWrapper 错误包装器
type ErrorWrapper struct {
	code int
	msg  string
}

// NewError 创建新错误
func NewError(code int, msg string) *ErrorWrapper {
	return &ErrorWrapper{
		code: code,
		msg:  msg,
	}
}

// Error 实现 error 接口
func (e *ErrorWrapper) Error() string {
	return fmt.Sprintf("%d: %s", e.code, e.msg)
}

// Errorf 格式化错误消息
func (e *ErrorWrapper) Errorf(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return NewError(e.code, msg)
}

// ============================================
// 常量定义
// ============================================

// DataSourceType 数据源类型常量
const (
	DataSourceTypeMySQL      = "mysql"
	DataSourceTypePostgreSQL = "postgresql"
	DataSourceTypeRedis      = "redis"
	DataSourceTypeMongoDB    = "mongodb"
	DataSourceTypeSQLServer  = "sqlserver"
)

// DataSourceStatus 数据源状态常量
const (
	DataSourceStatusEnabled  = "enabled"
	DataSourceStatusDisabled = "disabled"
)

// ============================================
// 错误信息定义
// ============================================

var (
	// 错误码常量 (30400-30499)
	ErrCodeParamInvalid    = 30400 // 参数错误
	ErrCodeNameExists      = 30409 // 名称重复
	ErrCodeNotFound        = 30411 // 资源不存在
	ErrCodeConnectionTest  = 30413 // 连接测试失败

	// 错误消息
	ErrMsgParamInvalid     = "参数错误"
	ErrMsgNameExists       = "数据源名称已存在"
	ErrMsgNotFound         = "数据源不存在"
	ErrMsgConnectionTest   = "数据源连接测试失败"
)

// ============================================
// 错误定义
// ============================================

var (
	// 参数错误
	ErrParamInvalid = NewError(ErrCodeParamInvalid, ErrMsgParamInvalid)

	// 名称重复错误
	ErrNameExists = NewError(ErrCodeNameExists, ErrMsgNameExists)

	// 资源不存在错误
	ErrNotFound = NewError(ErrCodeNotFound, ErrMsgNotFound)

	// 连接测试失败错误
	ErrConnectionTest = NewError(ErrCodeConnectionTest, ErrMsgConnectionTest)
)

// ============================================
// 默认值
// ============================================

const (
	// 默认排序字段
	DefaultSortField = "created_at"

	// 默认排序方向
	DefaultSortDirection = "desc"

	// 最大分页大小
	MaxPageSize = 2000
)
