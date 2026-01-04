package errorx

// ============================================
// 错误码定义 (30400-30499)
// ============================================

const (
	// 数据源管理模块错误码 (30400-30499)
	ErrCodeDataSourceParamInvalid   = 30400 // 参数错误
	ErrCodeDataSourceNameExists     = 30409 // 名称重复
	ErrCodeDataSourceNotFound       = 30411 // 资源不存在
	ErrCodeDataSourceConnectionTest = 30413 // 连接测试失败
)

// ============================================
// 错误消息定义
// ============================================

const (
	// 数据源管理模块错误消息
	ErrMsgDataSourceParamInvalid   = "参数错误：%s"
	ErrMsgDataSourceNameExists     = "数据源名称已存在"
	ErrMsgDataSourceNotFound       = "数据源不存在"
	ErrMsgDataSourceConnectionTest = "数据源连接测试失败：%s"
)

// ============================================
// 错误创建函数
// ============================================

// NewDataSourceParamInvalid 创建参数错误
func NewDataSourceParamInvalid(format string, args ...interface{}) error {
	return NewWithCode(ErrCodeDataSourceParamInvalid, format, args...)
}

// NewDataSourceNameExists 创建名称重复错误
func NewDataSourceNameExists() error {
	return NewWithCode(ErrCodeDataSourceNameExists, ErrMsgDataSourceNameExists)
}

// NewDataSourceNotFound 创建资源不存在错误
func NewDataSourceNotFound() error {
	return NewWithCode(ErrCodeDataSourceNotFound, ErrMsgDataSourceNotFound)
}

// NewDataSourceConnectionTest 创建连接测试失败错误
func NewDataSourceConnectionTest(format string, args ...interface{}) error {
	return NewWithCode(ErrCodeDataSourceConnectionTest, format, args...)
}
