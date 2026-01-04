package datasource

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
)

// ============================================
// 连接测试实现
// ============================================

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

// testMySQLConnection 测试 MySQL 连接
func (m *gormDataSourceModel) testMySQLConnection(ctx context.Context, config *DataSource) error {
	// TODO: 实现 MySQL 连接测试
	// - 使用 sql.Open 连接 MySQL
	// - 执行 SELECT 1 测试
	// - 返回连接结果
	return fmt.Errorf("MySQL 连接测试未实现")
}

// testPostgreSQLConnection 测试 PostgreSQL 连接
func (m *gormDataSourceModel) testPostgreSQLConnection(ctx context.Context, config *DataSource) error {
	// TODO: 实现 PostgreSQL 连接测试
	return fmt.Errorf("PostgreSQL 连接测试未实现")
}

// testRedisConnection 测试 Redis 连接
func (m *gormDataSourceModel) testRedisConnection(ctx context.Context, config *DataSource) error {
	// TODO: 实现 Redis 连接测试
	return fmt.Errorf("Redis 连接测试未实现")
}

// testMongoDBConnection 测试 MongoDB 连接
func (m *gormDataSourceModel) testMongoDBConnection(ctx context.Context, config *DataSource) error {
	// TODO: 实现 MongoDB 连接测试
	return fmt.Errorf("MongoDB 连接测试未实现")
}

// testSQLServerConnection 测试 SQLServer 连接
func (m *gormDataSourceModel) testSQLServerConnection(ctx context.Context, config *DataSource) error {
	// TODO: 实现 SQLServer 连接测试
	return fmt.Errorf("SQLServer 连接测试未实现")
}

// ============================================
// 加密解密实现
// ============================================

// EncryptPassword 加密密码
func (m *gormDataSourceModel) EncryptPassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("密码不能为空")
	}

	// 生成 32 字节密钥
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("生成密钥失败: %w", err)
	}

	// 使用 AES-256-GCM 加密
	encrypted, err := m.encryptWithAES(key, []byte(password))
	if err != nil {
		return "", fmt.Errorf("密码加密失败: %w", err)
	}

	// 组合密钥和加密结果
	result := make([]byte, 0, len(key)+len(encrypted))
	result = append(result, key...)
	result = append(result, encrypted...)

	// Base64 编码
	return base64.StdEncoding.EncodeToString(result), nil
}

// DecryptPassword 解密密码
func (m *gormDataSourceModel) DecryptPassword(encryptedPassword string) (string, error) {
	if encryptedPassword == "" {
		return "", fmt.Errorf("加密密码不能为空")
	}

	// Base64 解码
	data, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return "", fmt.Errorf("Base64 解码失败: %w", err)
	}

	if len(data) < 32 {
		return "", fmt.Errorf("加密数据格式错误")
	}

	// 提取密钥
	key := data[:32]
	encryptedData := data[32:]

	// AES 解密
	decrypted, err := m.decryptWithAES(key, encryptedData)
	if err != nil {
		return "", fmt.Errorf("密码解密失败: %w", err)
	}

	return string(decrypted), nil
}

// encryptWithAES 使用 AES-256-GCM 加密
func (m *gormDataSourceModel) encryptWithAES(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// decryptWithAES 使用 AES-256-GCM 解密
func (m *gormDataSourceModel) decryptWithAES(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("密文太短")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// GenerateUUIDv7 生成 UUID v7
func (m *gormDataSourceModel) GenerateUUIDv7() string {
	id, err := uuid.NewV7()
	if err != nil {
		// 如果 UUID v7 不可用，使用 v4
		id, err = uuid.NewV4()
		if err != nil {
			// 如果都失败，使用时间戳
			return fmt.Sprintf("%016d", time.Now().UnixNano())
		}
	}
	return id.String()
}
