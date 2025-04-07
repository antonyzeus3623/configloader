package configloader

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"golang.org/x/net/html/charset"
)

type ConfigLoader struct {
	vip      *viper.Viper
	tmpDir   string // 临时文件目录
	keepTemp bool   // 是否保留临时文件（调试用）
}

type Option func(*ConfigLoader)

// New 创建配置加载器实例
func New(options ...Option) *ConfigLoader {
	loader := &ConfigLoader{
		vip:      viper.New(),
		tmpDir:   os.TempDir(),
		keepTemp: false,
	}
	for _, opt := range options {
		opt(loader)
	}
	return loader
}

// WithTempDir 设置临时文件目录
func WithTempDir(dir string) Option {
	return func(cl *ConfigLoader) {
		cl.tmpDir = dir
	}
}

// WithKeepTempFile 保留临时文件（默认不保留）
func WithKeepTempFile(keep bool) Option {
	return func(cl *ConfigLoader) {
		cl.keepTemp = keep
	}
}

// Load 加载并解析配置文件
func (cl *ConfigLoader) Load(filePath string) error {
	rawData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read file failed: %w", err)
	}

	cleanData := removeBOM(rawData)
	utf8Data, err := convertToUTF8(cleanData)
	if err != nil {
		return fmt.Errorf("encoding conversion failed: %w", err)
	}

	// 创建临时文件
	ext := filepath.Ext(filePath)
	tmpFile, err := os.CreateTemp(cl.tmpDir, "config-*"+ext)
	if err != nil {
		return fmt.Errorf("create temp file failed: %w", err)
	}

	if !cl.keepTemp {
		defer os.Remove(tmpFile.Name())
	}

	if _, err := tmpFile.Write(utf8Data); err != nil {
		return fmt.Errorf("write temp file failed: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("close temp file failed: %w", err)
	}

	cl.vip.SetConfigFile(tmpFile.Name())
	if err := cl.vip.ReadInConfig(); err != nil {
		return fmt.Errorf("read config failed: %w", err)
	}
	return nil
}

// Unmarshal 将配置解析到结构体
func (cl *ConfigLoader) Unmarshal(cfg interface{}) error {
	return cl.vip.Unmarshal(cfg)
}

// convertToUTF8 编码转换核心逻辑
func convertToUTF8(input []byte) ([]byte, error) {
	if isUTF8(input) {
		return input, nil
	}

	utf8Reader, err := charset.NewReader(bytes.NewReader(input), "")
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, utf8Reader); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// isUTF8 高效检测UTF-8编码
func isUTF8(input []byte) bool {
	_, err := charset.NewReader(bytes.NewReader(input), "UTF-8")
	return err == nil
}

// removeBOM 移除字节顺序标记
func removeBOM(data []byte) []byte {
	if len(data) < 3 {
		return data
	}

	switch {
	case data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF: // UTF-8
		return data[3:]
	case data[0] == 0xFE && data[1] == 0xFF: // UTF-16BE
		return data[2:]
	case data[0] == 0xFF && data[1] == 0xFE: // UTF-16LE
		return data[2:]
	}
	return data
}
