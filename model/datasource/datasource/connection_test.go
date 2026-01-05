package datasource

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// ============================================
// 连接测试辅助函数
// ============================================

// EncryptPasswordForTest 加密密码（测试用）
func EncryptPasswordForTest(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("密码不能为空")
	}

	// 生成 32 字节密钥
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("生成密钥失败: %w", err)
	}

	// 使用 AES-256-GCM 加密
	encrypted, err := encryptWithAESForTest(key, []byte(password))
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

// DecryptPasswordForTest 解密密码（测试用）
func DecryptPasswordForTest(encryptedPassword string) (string, error) {
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
	decrypted, err := decryptWithAESForTest(key, encryptedData)
	if err != nil {
		return "", fmt.Errorf("密码解密失败: %w", err)
	}

	return string(decrypted), nil
}

// encryptWithAESForTest 使用 AES-256-GCM 加密（测试用）
func encryptWithAESForTest(key, plaintext []byte) ([]byte, error) {
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

// decryptWithAESForTest 使用 AES-256-GCM 解密（测试用）
func decryptWithAESForTest(key, ciphertext []byte) ([]byte, error) {
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
