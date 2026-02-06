package util

import (
	"fmt"
	"os"
)

// CreateTempFileFromString 从字符串创建临时文件
func CreateTempFileFromString(content, suffix string) (string, func(), error) {
	// 创建临时文件
	tempFile, err := os.CreateTemp("", fmt.Sprintf("rose_*%s", suffix))
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	// 写入内容
	_, err = tempFile.WriteString(content)
	if err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return "", nil, fmt.Errorf("failed to write temp file: %w", err)
	}

	// 关闭文件
	err = tempFile.Close()
	if err != nil {
		os.Remove(tempFile.Name())
		return "", nil, fmt.Errorf("failed to close temp file: %w", err)
	}

	// 返回文件路径和清理函数
	cleanup := func() {
		os.Remove(tempFile.Name())
	}

	return tempFile.Name(), cleanup, nil
}

// CreateTempFileFromBytes 从字节数据创建临时文件
func CreateTempFileFromBytes(data []byte, suffix string) (string, func(), error) {
	// 创建临时文件
	tempFile, err := os.CreateTemp("", fmt.Sprintf("rose_*%s", suffix))
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	// 写入数据
	_, err = tempFile.Write(data)
	if err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return "", nil, fmt.Errorf("failed to write temp file: %w", err)
	}

	// 关闭文件
	err = tempFile.Close()
	if err != nil {
		os.Remove(tempFile.Name())
		return "", nil, fmt.Errorf("failed to close temp file: %w", err)
	}

	// 返回文件路径和清理函数
	cleanup := func() {
		os.Remove(tempFile.Name())
	}

	return tempFile.Name(), cleanup, nil
}
