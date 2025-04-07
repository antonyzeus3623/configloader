# ConfigLoader- 智能配置文件加载库

 [![License: MIT](../../AA%20%E8%B5%84%E6%96%99/assets/readme.assets/License-MIT-yellow-17440158049161.svg)](https://opensource.org/licenses/MIT) 基于Viper的增强型智能配置文件加载库，提供专业级的配置文件加载方法。 

## 功能亮点 

1. **实例隔离设计**：

- 使用 `ConfigLoader` 结构体封装实例状态，避免全局变量污染
- 每个实例独立维护 Viper 配置，支持多配置文件加载

2. **功能增强**：

- 自动识别文件扩展名（支持 `yml`/`toml`/`json` 等格式）
- 可配置临时文件存储路径
- 增加调试模式保留临时文件

3. **性能优化**：

- 优化 `UTF-8` 检测逻辑（`isUTF8` 快速检测）
- 完善 `BOM` 处理（支持 `UTF-16LE`/`BE`）

4. **接口清晰化**：

- `New()` 配合 `Option` 模式实现灵活配置
- 分离 `Load` 和 `Unmarshal` 阶段，支持后续动态配置更新

5. **并发安全**：

- 通过 `sync.Once` 保证初始化原子性
- 每个实例独立操作，避免多 `goroutine` 竞争

## 快速开始 

### 安装依赖 

确保在项目中安装这些依赖：

```bash
go get github.com/spf13/viper
go get golang.org/x/net/html/charset
```

- `github.com/spf13/viper`: 用于加载和解析配置文件。

- `golang.org/x/net/html/charset`: 用于字符编码转换。

### 使用方法

1. **创建配置加载器**：

```go
loader := configloader.New(
    configloader.WithTempDir("./tmp"),
    configloader.WithKeepTempFile(true),
)
```

1. **加载配置文件**：

```go
err := loader.Load("config.toml")
```

1. **解析配置到结构体**：

```go
var cfg AppConfig
err := loader.Unmarshal(&cfg)
```

### 基础用法示例

```go
package main

import (
	"antonyzeus3623/configloader"
	"fmt"
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
```

