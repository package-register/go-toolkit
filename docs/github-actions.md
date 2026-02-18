# GitHub Actions 工作流指南

## 📋 概述

本文档详细介绍了 Go Toolkit 项目的 GitHub Actions 工作流配置、优化经验和最佳实践。

## 🚀 工作流架构

### 主要工作流
- **publish.yml** - 自动发布流程
- **未来可扩展**: CI 测试、代码检查、安全扫描等

### 工作流触发方式
```yaml
on:
  push:
    tags:
      - 'v*.*.*'          # 标签推送触发
  workflow_dispatch:      # 手动触发
    inputs:
      version: string     # 版本号
      skip_build: boolean # 跳过构建开关
```

## ⚡ 性能优化

### 1. Go 模块缓存
```yaml
- name: Setup Go
  uses: actions/setup-go@v5
  with:
    go-version: '1.24'
    cache: true                    # 启用缓存
    cache-dependency-path: go.sum  # 精确缓存失效
```

**效果**:
- 首次构建: 2-3分钟下载依赖
- 后续构建: 10-30秒复用缓存
- 总体时间减少 60-80%

### 2. 跳过二进制构建
```yaml
build:
  if: false  # 默认跳过，适合库项目
```

**适用场景**:
- ✅ 库项目发布
- ✅ 快速版本更新
- ❌ 需要二进制分发的应用

### 3. 并行执行优化
```yaml
jobs:
  release:     # 创建 Release
  build:       # 构建二进制 (可选)
  update-go-mod: # 更新模块
  notify:      # 通知完成
```

## 🔧 核心组件

### Release Job
```yaml
- name: Create Release
  uses: softprops/action-gh-release@v2
  with:
    tag_name: ${{ steps.version.outputs.version }}
    generate_release_notes: true
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

**关键点**:
- 必须提供 `GITHUB_TOKEN` 避免认证错误
- 自动生成发布说明
- 支持手动和自动触发

### 构建策略
```yaml
strategy:
  matrix:
    include:
      - goos: linux, goarch: amd64
      - goos: darwin, goarch: amd64
      - goos: darwin, goarch: arm64
      - goos: windows, goarch: amd64
```

## 🛠️ 常见问题解决

### 1. 认证错误
**问题**: `Bad credentials - https://docs.github.com/rest`

**解决方案**:
```yaml
env:
  GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### 2. 语法错误
**问题**: `Unexpected symbol: ','`

**解决方案**: 移除多余的逗号
```yaml
# 错误
${{ github.repository,,}}

# 正确
${{ github.repository }}
```

### 3. 标签依赖问题
**问题**: `GitHub Releases requires a tag`

**解决方案**: 确保标签存在且格式正确
```bash
git tag v0.1.0
git push origin v0.1.0
```

## 📝 最佳实践

### 1. 版本管理
- 使用语义化版本: `v1.2.3`
- 标签与 Release 一一对应
- 自动生成 CHANGELOG

### 2. 缓存策略
- 使用 `go.sum` 作为缓存键
- 定期清理过期缓存
- 监控缓存命中率

### 3. 错误处理
```yaml
- name: Trigger Go module reindex
  run: |
    curl -sSf "https://proxy.golang.org/..." || true
```

### 4. 条件执行
```yaml
if: ${{ github.event_name != 'workflow_dispatch' || github.event.inputs.skip_build != 'true' }}
```

## 🔍 监控和调试

### 查看执行状态
1. GitHub Actions 页面
2. 工作流执行日志
3. 缓存使用情况

### 调试技巧
```yaml
- name: Debug Info
  run: |
    echo "Event: ${{ github.event_name }}"
    echo "Ref: ${{ github.ref }}"
    echo "Version: ${{ steps.version.outputs.version }}"
```

## 📚 扩展建议

### 未来可添加的工作流
- **CI**: 代码测试、覆盖率检查
- **Security**: 依赖安全扫描
- **Docs**: 文档自动部署
- **Release Notes**: 自动生成更新日志

### 性能进一步优化
- 使用自定义运行器
- 并行化测试
- 增量构建

---

## 📞 支持

如有问题，请查看：
- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [Go Toolkit Issues](https://github.com/package-register/go-toolkit/issues)
- 项目 Actions 执行日志
