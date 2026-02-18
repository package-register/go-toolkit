# 🛠️ Go Toolkit

![Go Version](https://img.shields.io/badge/go-%3E%3D1.24-blue)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![GitHub release](https://img.shields.io/github/release/package-register/go-toolkit.svg)](https://github.com/package-register/go-toolkit/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/package-register/go-toolkit)](https://goreportcard.com/report/github.com/package-register/go-toolkit)

🚀 **统一的 Go 开发工具包** - 整合基础工具链与专用功能库，提供完整的开发解决方案

## ✨ 特性

### 🔧 基础工具链
- **🏗️ 构建工具** (`build/`) - 跨平台编译、测试覆盖率、版本管理
- **💾 缓存组件** (`cache/`) - 内存缓存、并发安全、失效策略
- **🐳 Docker 工具** (`docker/`) - 镜像构建、容器管理、健康检查
- **🔄 GitOps** (`gitops/`) - 版本标签、自动化发布、GitHub Actions
- **🤖 浏览器自动化** (`rod/`) - 基于 Rod 的网页操作工具

### 🎯 专用功能
- **🌐 翻译服务** (`trans/`) - 讯飞 API 集成、多语言支持
- **🖼️ 图像处理** (`image/`) - 时间表生成、图片处理工具
- **🔍 服务发现** (`discovery/`) - 网络服务发现机制
- **🌐 ZeroTier SDK** (`zerotier/`) - ZeroTier 网络管理 API

## 🚀 快速开始

### 安装

```bash
go get github.com/package-register/go-toolkit
```

### 基础使用

```go
package main

import (
    "time"
    "fmt"
    "github.com/package-register/go-toolkit/cache"
    "github.com/package-register/go-toolkit/build"
    "github.com/package-register/go-toolkit/trans"
    "github.com/package-register/go-toolkit/zerotier"
)

func main() {
    // 缓存使用
    c := cache.New()
    c.Set("key", "value", 10*time.Minute)

    // 翻译使用
    t := trans.New(
        trans.WithAppID("your-app-id"),
        trans.WithSecret("your-secret"),
        trans.WithFromLang("cn"),
        trans.WithToLang("en"),
    )

    result, err := t.Translate("你好世界")
    if err == nil {
        fmt.Println(result)
    }

    // ZeroTier 使用
    local := zerotier.NewClient()
    status, _ := local.Status()
    fmt.Printf("ZeroTier 节点状态: %v\n", status.Online)
}
```

## 📦 模块详情

### 🏗️ 构建工具

```go
import "github.com/package-register/go-toolkit/build"

// 跨平台构建
builder := build.New(
    build.WithPlatforms(build.Platform{OS: "linux", Arch: "amd64"}),
    build.WithPlatforms(build.Platform{OS: "darwin", Arch: "arm64"}),
)

err := builder.Build("./cmd/main.go")
```

### 💾 缓存组件

```go
import "github.com/package-register/go-toolkit/cache"

// 创建缓存实例
c := cache.New()

// 设置缓存
c.Set("user:123", userData, 30*time.Minute)

// 获取缓存
if value, found := c.Get("user:123"); found {
    fmt.Println(value)
}

// 删除缓存
c.Delete("user:123")
```

### 🌐 翻译服务

```go
import "github.com/package-register/go-toolkit/trans"

// 创建翻译器
translator := trans.New(
    trans.WithAppID("your-app-id"),
    trans.WithSecret("your-secret"),
    trans.WithAPIKey("your-api-key"),
    trans.WithFromLang("cn"),
    trans.WithToLang("en"),
)

// 翻译文本
result, err := translator.Translate("今天天气很好")
if err == nil {
    fmt.Println(result)
}

// 获取结构化结果
res, err := translator.TranslateWithResult("你好世界")
if err == nil {
    fmt.Printf("原文: %s\n译文: %s\n", res.Source, res.Target)
}
```

### � Docker 工具

```go
import "github.com/package-register/go-toolkit/docker"

// 创建 Docker 客户端
client, err := docker.NewClient()
if err != nil {
    log.Fatal(err)
}

// 构建镜像
err := client.BuildImage(docker.BuildOptions{
    Context:    "./",
    Dockerfile: "Dockerfile",
    Tag:        "myapp:latest",
})
```

### 🌐 ZeroTier SDK

```go
import "github.com/package-register/go-toolkit/zerotier"

// 本地节点管理
local := zerotier.NewClient()
status, _ := local.Status()
networks, _ := local.Networks().List()

// 云端管理
cloud := zerotier.NewCentral("your_api_token")
networks, _ := cloud.Networks().List()
cloud.Networks().Members("nwid").Authorize("member_id")
```

## 🛠️ 开发工具

### Make 命令

```bash
# 查看帮助
make help

# 版本管理
make version          # 查看当前版本
make bump-version     # 交互式升级版本
make patch           # 升级补丁版本
make minor           # 升级次版本
make major           # 升级主版本

# 发布流程
make release          # 完整发布流程
make test            # 运行测试
make build           # 构建二进制
make clean           # 清理产物
```

### 自动化发布

1. **版本升级**: `make bump-version`
2. **自动发布**: `make release`
3. **GitHub Actions**: 自动构建多平台二进制文件
4. **Release**: 自动创建 GitHub Release

## 🧪 测试

```bash
# 运行所有测试
make test

# 运行特定模块测试
go test ./build/...
go test ./cache/...
go test ./trans/...
```

## 📋 系统要求

- **Go**: >= 1.24
- **Git**: >= 2.30
- **Docker**: >= 20.10 (可选)

## 🗂️ 项目结构

```
go-toolkit/
├── build/          # 构建工具
├── cache/          # 缓存组件
├── docker/         # Docker 工具
├── gitops/         # GitOps 工具
├── rod/            # 浏览器自动化
├── trans/          # 翻译服务
├── image/          # 图像处理
├── discovery/      # 服务发现
├── zerotier/       # ZeroTier SDK
├── docs/           # 开发文档
│   ├── README.md   # 文档索引
│   ├── github-actions.md  # GitHub Actions 指南
│   └── github-cli.md      # GitHub CLI 指南
├── examples/       # 使用示例
├── .github/        # GitHub Actions
├── Makefile        # 构建脚本
├── go.mod          # Go 模块
└── README.md       # 项目文档
```

## 🤝 贡献

我们欢迎所有形式的贡献！

### 开发流程

1. **Fork** 项目
2. **创建** 特性分支: `git checkout -b feature/amazing-feature`
3. **提交** 变更: `git commit -m 'Add amazing feature'`
4. **推送** 分支: `git push origin feature/amazing-feature`
5. **创建** Pull Request

### 代码规范

- 遵循 Go 官方代码规范
- 添加适当的测试用例
- 更新相关文档
- 确保 `make test` 通过

## � 许可证

本项目采用 [MIT 许可证](LICENSE)

## 🔗 相关链接

- **GitHub**: https://github.com/package-register/go-toolkit
- **Releases**: https://github.com/package-register/go-toolkit/releases
- **Issues**: https://github.com/package-register/go-toolkit/issues
- **文档**: https://pkg.go.dev/github.com/package-register/go-toolkit

## � 版本历史

查看 [CHANGELOG.md](CHANGELOG.md) 了解详细的版本变更记录。

---

🦄 **Made with ❤️ by oAo Team** | 📧 hnkong666@gmail.com
