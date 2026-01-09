package util

import (
	"crypto/aes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

// RSADecrypt RSA 解密函数
// privateKeyStr: Base64 编码的 RSA 私钥（PKCS#1 格式）
// encryptedData: Base64 编码的加密数据
// 返回: 解密后的原始数据（字符串）
func RSADecrypt(privateKeyStr, encryptedData string) (string, error) {
	// 解码 Base64 私钥
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		return "", fmt.Errorf("解码私钥失败: %w", err)
	}

	// 解析 PEM 格式的私钥
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		// 如果不是 PEM 格式，尝试直接解析 DER 格式
		privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
		if err != nil {
			// 尝试 PKCS#8 格式
			key, err := x509.ParsePKCS8PrivateKey(privateKeyBytes)
			if err != nil {
				return "", fmt.Errorf("解析私钥失败: %w", err)
			}
			privateKey, ok := key.(*rsa.PrivateKey)
			if !ok {
				return "", errors.New("私钥类型错误，不是 RSA 私钥")
			}
			return decryptWithPrivateKey(privateKey, encryptedData)
		}
		return decryptWithPrivateKey(privateKey, encryptedData)
	}

	// 解析 PEM 块中的私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// 尝试 PKCS#8 格式
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return "", fmt.Errorf("解析 PEM 私钥失败: %w", err)
		}
		privateKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return "", errors.New("私钥类型错误，不是 RSA 私钥")
		}
		return decryptWithPrivateKey(privateKey, encryptedData)
	}

	return decryptWithPrivateKey(privateKey, encryptedData)
}

// decryptWithPrivateKey 使用 RSA 私钥解密数据
func decryptWithPrivateKey(privateKey *rsa.PrivateKey, encryptedData string) (string, error) {
	// 解码 Base64 加密数据
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", fmt.Errorf("解码加密数据失败: %w", err)
	}

	// RSA 解密（使用 PKCS#1 v1.5 填充）
	plaintext, err := rsa.DecryptPKCS1v15(nil, privateKey, ciphertext)
	if err != nil {
		return "", fmt.Errorf("RSA 解密失败: %w", err)
	}

	return string(plaintext), nil
}

// AESDecrypt AES 解密函数
// encryptedData: Base64 编码的加密数据
// key: AES 密钥（字符串）
// 返回: 解密后的原始数据（字符串）
func AESDecrypt(encryptedData, key string) (string, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", fmt.Errorf("AES密钥长度要求为16位、24位或32位，当前长度: %d", len(key))
	}

	// 解码 Base64 加密数据
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", fmt.Errorf("解码加密数据失败: %w", err)
	}

	// 创建 AES cipher
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("创建AES cipher失败: %w", err)
	}

	// AES 块大小必须是 16 字节
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("密文长度不足")
	}

	// 使用 ECB 模式（Java 的 SecureUtil.aes 默认使用 ECB 模式）
	// 注意：Go 标准库不直接支持 ECB，需要手动实现
	plaintext := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i += aes.BlockSize {
		block.Decrypt(plaintext[i:i+aes.BlockSize], ciphertext[i:i+aes.BlockSize])
	}

	// 去除 PKCS5/PKCS7 填充
	plaintext, err = pkcs5Unpadding(plaintext)
	if err != nil {
		return "", fmt.Errorf("去除填充失败: %w", err)
	}

	return string(plaintext), nil

}

// pkcs5Unpadding 去除 PKCS5/PKCS7 填充
func pkcs5Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("数据为空")
	}

	unpadding := int(data[length-1])
	if unpadding > length {
		return nil, errors.New("填充长度错误")
	}

	return data[:(length - unpadding)], nil
}
