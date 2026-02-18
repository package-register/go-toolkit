# Go Toolkit Makefile
# 精简、清晰、跨平台的版本管理工具

# Variables
BINARY_NAME := go-toolkit
REMOTE_REPO ?= origin
VERSION_FILE := VERSION
GO_VERSION := 1.24

# Colors for terminal output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
PURPLE := \033[0;35m
CYAN := \033[0;36m
WHITE := \033[0;37m
BOLD := \033[1m
RESET := \033[0m

# Default target
.DEFAULT_GOAL := help

.PHONY: help version bump-version release test clean build install

# Help target - 显示所有可用命令
help:
	@echo "$(BOLD)$(CYAN)🛠️  Go Toolkit - 版本管理工具$(RESET)"
	@echo ""
	@echo "$(BOLD)📋 可用命令:$(RESET)"
	@echo "  $(GREEN)help$(RESET)           显示此帮助信息"
	@echo "  $(GREEN)version$(RESET)        显示当前版本"
	@echo "  $(GREEN)bump-version$(RESET)   自动升级版本"
	@echo "  $(GREEN)release$(RESET)        完整发布流程"
	@echo "  $(GREEN)test$(RESET)           运行所有测试"
	@echo "  $(GREEN)build$(RESET)           构建二进制文件"
	@echo "  $(GREEN)clean$(RESET)           清理构建产物"
	@echo "  $(GREEN)install$(RESET)         本地安装"
	@echo ""
	@echo "$(BOLD)🚀 快速开始:$(RESET)"
	@echo "  make release                    # 完整发布流程"
	@echo "  make bump-version patch         # 升级补丁版本"
	@echo "  make bump-version minor         # 升级次版本"
	@echo "  make bump-version major         # 升级主版本"

# Version target - 显示当前版本
version:
	@if [ -f $(VERSION_FILE) ]; then \
		echo "$(BOLD)$(BLUE)📦 当前版本:$(RESET) $(GREEN)$$(cat $(VERSION_FILE))$(RESET)"; \
	else \
		echo "$(YELLOW)⚠️  未找到版本文件，初始化为 v0.1.0$(RESET)"; \
		echo "v0.1.0" > $(VERSION_FILE); \
	fi

# Bump version target - 自动升级版本
bump-version:
	@echo "$(BOLD)$(YELLOW)🔄 版本升级工具$(RESET)"
	@echo ""
	@if [ ! -f $(VERSION_FILE) ]; then \
		echo "v0.1.0" > $(VERSION_FILE); \
	fi
	@CURRENT_VERSION=$$(cat $(VERSION_FILE)); \
	echo "$(CYAN)当前版本:$(RESET) $$CURRENT_VERSION"; \
	echo ""; \
	echo "$(BOLD)请选择升级类型:$(RESET)"; \
	echo "1) patch  - 补丁版本 (0.1.$$(echo $$CURRENT_VERSION | cut -d. -f3 | sed 's/v//') → 0.1.$$(($$(echo $$CURRENT_VERSION | cut -d. -f3 | sed 's/v//') + 1)))"; \
	echo "2) minor  - 次版本 (0.1 → 0.2)"; \
	echo "3) major  - 主版本 (0 → 1)"; \
	echo ""; \
	read -p "$(BOLD)请输入选择 [1-3]:$(RESET) " choice; \
	case $$choice in \
		1) \
			NEW_VERSION=$$(echo $$CURRENT_VERSION | awk -F. '{printf "v%d.%d.%d", $$1, $$2, $$3+1}' | sed 's/v//g' | sed 's/^/v/'); \
			echo "$(GREEN)✅ 补丁版本升级: $$CURRENT_VERSION → $$NEW_VERSION$(RESET)"; \
			;; \
		2) \
			NEW_VERSION=$$(echo $$CURRENT_VERSION | awk -F. '{printf "v%d.%d.0", $$1, $$2+1}' | sed 's/v//g' | sed 's/^/v/'); \
			echo "$(GREEN)✅ 次版本升级: $$CURRENT_VERSION → $$NEW_VERSION$(RESET)"; \
			;; \
		3) \
			NEW_VERSION=$$(echo $$CURRENT_VERSION | awk -F. '{printf "v%d.0.0", $$1+1}' | sed 's/v//g' | sed 's/^/v/'); \
			echo "$(GREEN)✅ 主版本升级: $$CURRENT_VERSION → $$NEW_VERSION$(RESET)"; \
			;; \
		*) \
			echo "$(RED)❌ 无效选择$(RESET)"; \
			exit 1; \
			;; \
	esac; \
	echo $$NEW_VERSION > $(VERSION_FILE); \
	echo "$(BOLD)$(CYAN)📝 更新日志:$(RESET)"; \
	read -p "$(BOLD)请输入此版本的变更描述:$(RESET) " changelog; \
	echo "$$NEW_VERSION: $$changelog" >> CHANGELOG.md; \
	echo "$(GREEN)✅ 版本文件已更新$(RESET)"

# Release target - 完整发布流程
release: check-git test
	@echo "$(BOLD)$(PURPLE)🚀 开始发布流程$(RESET)"
	@echo ""
	@$(MAKE) version
	@echo ""
	@echo "$(BOLD)📋 发布前检查:$(RESET)"
	@echo "✓ 测试已通过"
	@echo "✓ 工作区状态检查"
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "$(YELLOW)⚠️  发现未提交的变更$(RESET)"; \
		read -p "$(BOLD)是否提交这些变更? [y/N]:$(RESET) " commit; \
		if [ "$$commit" = "y" ] || [ "$$commit" = "Y" ]; then \
			git add .; \
			git commit -m "🔖 Release preparation for $$(cat $(VERSION_FILE))"; \
			echo "$(GREEN)✅ 变更已提交$(RESET)"; \
		fi; \
	fi
	@echo ""
	@echo "$(BOLD)🏷️  创建版本标签$(RESET)"
	@VERSION=$$(cat $(VERSION_FILE)); \
	git tag -a "$$VERSION" -m "Release $$VERSION"; \
	echo "$(GREEN)✅ 标签 $$VERSION 已创建$(RESET)"
	@echo ""
	@echo "$(BOLD)📡 推送到远程仓库$(RESET)"
	@git push $(REMOTE_REPO) main --follow-tags; \
	echo "$(GREEN)✅ 代码和标签已推送$(RESET)"
	@echo ""
	@echo "$(BOLD)🎉 发布完成!$(RESET)"
	@echo "$(CYAN)🔗 GitHub Actions 将自动处理发布流程$(RESET)"
	@echo "$(CYAN)🔗 查看进度: https://github.com/package-register/go-toolkit/actions$(RESET)"
	@echo "$(CYAN)🔗 查看发布: https://github.com/package-register/go-toolkit/releases$(RESET)"

# Test target - 运行测试
test:
	@echo "$(BOLD)$(BLUE)🧪 运行测试$(RESET)"
	@go test -v ./...
	@echo "$(GREEN)✅ 所有测试通过$(RESET)"

# Build target - 构建二进制文件
build:
	@echo "$(BOLD)$(BLUE)🔨 构建二进制文件$(RESET)"
	@VERSION=$$(cat $(VERSION_FILE) 2>/dev/null || echo "dev"); \
	mkdir -p dist; \
	for os in linux darwin windows; do \
		for arch in amd64 arm64; do \
			if [ "$$os" = "windows" ] && [ "$$arch" = "arm64" ]; then continue; fi; \
			echo "$(CYAN)构建 $$os/$$arch...$(RESET)"; \
			GOOS=$$os GOARCH=$$arch go build \
				-ldflags "-s -w -X main.version=$$VERSION" \
				-o "dist/$(BINARY_NAME)-$$VERSION-$$os-$$arch$$([ "$$os" = "windows" ] && echo .exe)" \
				./cmd/main.go 2>/dev/null || echo "$(YELLOW)⚠️  跳过库项目构建$(RESET)"; \
		done; \
	done
	@echo "$(GREEN)✅ 构建完成$(RESET)"
	@ls -la dist/

# Clean target - 清理构建产物
clean:
	@echo "$(BOLD)$(YELLOW)🧹 清理构建产物$(RESET)"
	@rm -rf dist/
	@echo "$(GREEN)✅ 清理完成$(RESET)"

# Install target - 本地安装
install: build
	@echo "$(BOLD)$(BLUE)📦 本地安装$(RESET)"
	@if [ -f "dist/$(BINARY_NAME)-$$(cat $(VERSION_FILE) 2>/dev/null || echo dev)-$$(go env GOOS)-$$(go env GOARCH)" ]; then \
		cp "dist/$(BINARY_NAME)-$$(cat $(VERSION_FILE) 2>/dev/null || echo dev)-$$(go env GOOS)-$$(go env GOARCH)" "$$(go env GOPATH)/bin/$(BINARY_NAME)"; \
		echo "$(GREEN)✅ 安装完成: $$(go env GOPATH)/bin/$(BINARY_NAME)$(RESET)"; \
	else \
		echo "$(YELLOW)⚠️  这是一个库项目，无需安装二进制文件$(RESET)"; \
	fi

# Check git status
check-git:
	@if ! git rev-parse --git-dir > /dev/null 2>&1; then \
		echo "$(RED)❌ 错误: 当前目录不是 Git 仓库$(RESET)"; \
		exit 1; \
	fi
	@if ! git remote | grep -q "$(REMOTE_REPO)"; then \
		echo "$(RED)❌ 错误: 未配置远程仓库 $(REMOTE_REPO)$(RESET)"; \
		echo "$(CYAN)请运行: git remote add origin <repository-url>$(RESET)"; \
		exit 1; \
	fi

# Quick version bump shortcuts
patch:
	@echo "1" | $(MAKE) bump-version

minor:
	@echo "2" | $(MAKE) bump-version

major:
	@echo "3" | $(MAKE) bump-version
