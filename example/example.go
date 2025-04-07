package main

import (
	"fmt"
	"github.com/antonyzeus3623/configloader"
	"log"
	"sync"
)

var (
	conf *Config
	once sync.Once // 使用 sync.Once 保证只初始化一次
)

type Config struct {
	ProcessName string `mapstructure:"ProcessName"`
}

func LoadConfig(configPath string) (*Config, error) {
	var loadErr error
	once.Do(func() { // 使用 sync.Once 保证只初始化一次
		// 创建加载器（可配置选项）
		loader := configloader.New(
			configloader.WithTempDir("./tmp"), // 自定义临时目录
		)

		if err := loader.Load(configPath); err != nil {
			loadErr = err
			return
		}

		if err := loader.Unmarshal(&conf); err != nil {
			loadErr = fmt.Errorf("unmarshal config failed: %w", err)
			return
		}
	})

	return conf, loadErr
}

func main() {
	conf, err := LoadConfig("config.toml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Loaded config: %+v", conf)
}
