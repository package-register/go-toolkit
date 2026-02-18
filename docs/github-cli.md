# GitHub CLI (gh) 使用指南

## 📋 概述

GitHub CLI (`gh`) 是 GitHub 的官方命令行工具，让你可以直接在终端中与 GitHub 交互。本文档总结了 Go Toolkit 项目开发中的实用经验和最佳实践。

## 🚀 安装和配置

### 安装
```bash
# macOS
brew install gh

# 其他平台
# 访问 https://github.com/cli/cli#installation
```

### 认证
```bash
# 登录 GitHub
gh auth login

# 查看认证状态
gh auth status

# 登出
gh auth logout
```

### 配置
```bash
# 设置默认编辑器
gh config set editor code

# 查看配置
gh config list

# 设置协议
gh config set git_protocol https
```

## 🛠️ 核心功能

### 1. 仓库管理
```bash
# 创建仓库
gh repo create package-register/go-toolkit --public --description "统一的 Go 开发工具包"

# 克隆仓库
gh repo clone package-register/go-toolkit

# 查看仓库信息
gh repo view

# 编辑仓库
gh repo edit --description "新描述" --homepage "https://example.com"
```

### 2. 协作者管理
```bash
# 添加协作者
gh api repos/package-register/go-toolkit/collaborators/username -X PUT

# 查看协作者
gh api repos/package-register/go-toolkit/collaborators

# 移除协作者
gh api repos/package-register/go-toolkit/collaborators/username -X DELETE
```

### 3. Issues 和 PR
```bash
# 创建 Issue
gh issue create --title "Bug: 功能异常" --body "详细描述..."

# 查看 Issues
gh issue list

# 创建 PR
gh pr create --title "新功能" --body "变更描述..."

# 查看 PR
gh pr list
```

### 4. Release 管理
```bash
# 创建 Release
gh release create v1.0.0 --title "版本 1.0.0" --notes "发布说明"

# 查看 Releases
gh release list

# 下载 Release 资源
gh release download v1.0.0
```

### 5. 工作流管理
```bash
# 查看工作流
gh workflow list

# 运行工作流
gh workflow run publish.yml --field version=v0.1.4

# 查看工作流执行
gh run list

# 查看执行详情
gh run view 123456
```

## 📝 Go Toolkit 项目经验

### 项目初始化流程
```bash
# 1. 创建仓库
gh repo create package-register/go-toolkit --public --description "统一的 Go 开发工具包"

# 2. 初始化本地 Git
git init
git remote add origin git@github.com:package-register/go-toolkit.git

# 3. 添加协作者
gh api repos/package-register/go-toolkit/collaborators/Fromsko -X PUT
```

### 发布流程
```bash
# 1. 提交代码
git add .
git commit -m "feat: 新功能发布"

# 2. 推送代码
git push origin main

# 3. 创建标签
git tag v0.1.0

# 4. 推送标签
git push origin v0.1.0

# 5. 触发工作流 (可选)
gh workflow run publish.yml --field version=v0.1.4 --field skip_build=false
```

### 权限问题解决
```bash
# 问题: Permission denied
# 解决: 切换到 SSH 协议
git remote set-url origin git@github.com:package-register/go-toolkit.git

# 或添加协作者权限
gh api repos/package-register/go-toolkit/collaborators/username -X PUT
```

## 🔧 高级用法

### 1. API 调用
```bash
# 直接调用 GitHub API
gh api repos/package-register/go-toolkit
gh api repos/package-register/go-toolkit/issues
gh api user

# POST 请求
gh api repos/package-register/go-toolkit/issues -f title="Bug" -f body="Description"

# PUT 请求
gh api repos/package-register/go-toolkit/collaborators/username -X PUT
```

### 2. 批量操作
```bash
# 批量关闭 Issues
for issue in $(gh issue list --json number --jq '.[].number'); do
  gh issue close $issue --comment "批量关闭"
done

# 批量创建标签
for version in v1.0.0 v1.1.0 v1.2.0; do
  gh release create $version --title "Release $version" --notes "自动发布"
done
```

### 3. 脚本集成
```bash
#!/bin/bash
# 自动发布脚本

VERSION=$1
if [ -z "$VERSION" ]; then
  echo "Usage: $0 <version>"
  exit 1
fi

# 创建标签和推送
git tag $VERSION
git push origin $VERSION

# 等待 Actions 完成
echo "等待 GitHub Actions 完成..."
sleep 30

# 查看执行状态
gh run list --workflow=publish.yml
```

## 🎯 最佳实践

### 1. 安全性
- 使用 SSH 密钥而非 HTTPS
- 定期更新 gh CLI
- 使用个人访问令牌 (PAT) 进行自动化

### 2. 效率提升
- 配置默认编辑器
- 使用别名简化命令
- 批量操作减少重复工作

### 3. 工作流集成
- 结合 Makefile 使用
- 自动化发布流程
- 监控执行状态

## 📚 常用命令速查

```bash
# 仓库操作
gh repo create <name>           # 创建仓库
gh repo clone <name>            # 克隆仓库
gh repo view                    # 查看仓库

# 认证管理
gh auth login                   # 登录
gh auth status                  # 查看状态
gh auth logout                  # 登出

# 发布管理
gh release create <tag>         # 创建 Release
gh release list                 # 查看 Releases
gh release download <tag>       # 下载资源

# 工作流
gh workflow list                # 查看工作流
gh workflow run <name>          # 运行工作流
gh run list                     # 查看执行记录

# API 调用
gh api <endpoint>               # API 请求
gh api <endpoint> -X POST      # POST 请求
gh api <endpoint> -f key=value  # 表单数据
```

## 🔗 相关资源

- [GitHub CLI 官方文档](https://cli.github.com/manual/)
- [gh 项目主页](https://github.com/cli/cli)
- [Go Toolkit 项目](https://github.com/package-register/go-toolkit)
- [GitHub API 文档](https://docs.github.com/en/rest)

---

💡 **提示**: 将常用命令添加到 shell 别名中可以大幅提升效率！
