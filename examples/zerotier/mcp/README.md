# ZeroTier MCP 集成

这个演示展示了 Go Toolkit 中新增的 MCP (Model Context Protocol) 集成功能，借鉴了 Rust 版本的高级特性。

## 🚀 新增功能

### MCP 协议支持
- 🤖 **AI 助手集成**: 通过 MCP 协议让 AI 助手直接操作 ZeroTier
- 📝 **结构化工具**: 定义了完整的工具集，支持参数验证
- 🔧 **标准化接口**: 统一的工具执行和结果返回格式

### 智能配置管理
- 🔐 **自动 Token 读取**: 从系统默认位置自动读取本地 Token
- 🌍 **环境变量支持**: 支持 `ZEROTIER_LOCAL_TOKEN` 和 `ZEROTIER_CENTRAL_TOKEN`
- 🛡️ **跨平台兼容**: 支持 macOS、Linux、Windows 的 Token 路径

### 增强的网络管理
- 🎯 **IP 指定授权**: 授权成员时可以直接分配指定 IP
- 📊 **状态可视化**: 使用 ✅/❌ 状态标识和表情符号
- 🔍 **详细反馈**: 提供友好的中文错误信息和操作反馈

## 🛠️ MCP 工具列表

### 本地 API 工具
| 工具名 | 描述 | 参数 |
|--------|------|------|
| `zt_status` | 获取本地节点状态 | 无 |
| `zt_networks` | 列出已加入的网络 | 无 |
| `zt_join` | 加入网络 | `network_id` (必需) |
| `zt_leave` | 离开网络 | `network_id` (必需) |
| `zt_peers` | 列出所有节点 | 无 |

### 云端 API 工具
| 工具名 | 描述 | 参数 |
|--------|------|------|
| `zt_central_networks` | 列出云端网络 | 无 |
| `zt_central_members` | 列出网络成员 | `network_id` (必需) |
| `zt_central_authorize` | 授权成员 | `network_id`, `member_id` (必需) |
| `zt_central_authorize_with_ip` | 授权成员并指定IP | `network_id`, `member_id`, `ip_address` (必需), `name` (可选) |
| `zt_central_deauthorize` | 取消授权 | `network_id`, `member_id` (必需) |

## 🚀 快速开始

### 运行 MCP 演示

```bash
cd examples/zerotier/mcp
go run mcp_demo.go
```

### 环境配置

```bash
# 设置本地 Token (可选，会自动读取)
export ZEROTIER_LOCAL_TOKEN="your_local_token"

# 设置云端 Token (必需使用云端功能)
export ZEROTIER_CENTRAL_TOKEN="your_central_token"
```

### 获取 Token

**本地 Token**:
- macOS: `~/Library/Application Support/ZeroTier/authtoken.secret`
- Linux: `/var/lib/zerotier-one/authtoken.secret`
- Windows: `C:\ProgramData\ZeroTier\One\authtoken.secret`

**云端 Token**: 访问 [my.zerotier.com](https://my.zerotier.com) → Account → API Access Tokens

## 📝 使用示例

### 基础用法

```go
package main

import (
    "fmt"
    "github.com/package-register/go-toolkit/zerotier/mcp"
)

func main() {
    // 创建 MCP 服务器
    server := mcp.NewMcpServer()
    
    // 设置云端 Token
    server = server.WithCentralToken("your_token")
    
    // 执行工具
    result := server.ExecuteTool("zt_status", nil)
    if result.Success {
        fmt.Println(result.Data)
    } else {
        fmt.Println("错误:", result.Error)
    }
}
```

### AI 助手集成示例

```json
{
  "mcpServers": {
    "zerotier": {
      "command": "go",
      "args": ["run", "/path/to/mcp_server.go"],
      "env": {
        "ZEROTIER_CENTRAL_TOKEN": "your_token"
      }
    }
  }
}
```

### 高级功能

```go
// 带参数的工具执行
params := map[string]interface{}{
    "network_id": "8056c2e21c000001",
    "member_id": "1234567890abcdef",
    "ip_address": "10.147.20.100",
    "name": "生产服务器",
}

result := server.ExecuteTool("zt_central_authorize_with_ip", params)
```

## 🔧 配置选项

### 自动配置

```go
// 自动加载配置
config := mcp.LoadConfig()
server := mcp.NewMcpServer()

if config.CentralToken != "" {
    server = server.WithCentralToken(config.CentralToken)
}
```

### 手动配置

```go
// 手动设置 Token
server := mcp.NewMcpServer().
    WithCentralToken("your_central_token")
```

## 📊 输出格式

### 成功响应
```json
{
  "success": true,
  "data": "🌐 节点状态\n📍 地址: 1234567890\n🔗 在线: true"
}
```

### 错误响应
```json
{
  "success": false,
  "error": "❌ 未配置 Central API Token"
}
```

## 🎯 与 Rust 版本对比

| 功能 | Go Toolkit | Rust MCP | 优势 |
|------|------------|----------|------|
| **MCP 协议** | ✅ 完整支持 | ✅ 原生支持 | 相同功能 |
| **自动配置** | ✅ 智能读取 | ✅ 环境变量 | 更灵活 |
| **错误处理** | ✅ 中文提示 | ✅ 英文提示 | 更友好 |
| **AI 集成** | ✅ 标准化 | ✅ 原生集成 | 相同体验 |

## 🚀 未来计划

1. **完整 MCP 服务器**: 独立的 MCP 服务器进程
2. **插件系统**: 支持自定义工具扩展
3. **Web 界面**: 基于 MCP 的 Web 管理界面
4. **批量操作**: 支持批量网络和成员管理

---

🎉 现在Go Toolkit 具备了与 Rust 版本相同的高级功能，同时保持了 Go 的开发效率优势！
