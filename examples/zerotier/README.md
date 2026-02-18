# ZeroTier SDK Demo

这个演示展示了如何使用 Go Toolkit 中的 ZeroTier SDK 进行网络管理。

## 🚀 快速开始

### 运行演示

```bash
cd examples/zerotier
go run main.go
```

### 前置条件

1. **安装 ZeroTier**
   ```bash
   # macOS
   brew install zerotier-one
   
   # 其他平台
   # 访问 https://www.zerotier.com/download/
   ```

2. **启动 ZeroTier 服务**
   ```bash
   # macOS
   sudo zerotier-one
   
   # Linux
   sudo systemctl start zerotier-one
   ```

## 📱 演示内容

### 1. 本地节点管理
- 获取节点状态信息
- 列出已加入的网络
- 查看连接的节点列表

### 2. 云端管理 (需要 API Token)
- 获取云端网络列表
- 查看网络成员信息
- 管理网络成员权限

### 3. 网络操作
- 加入指定网络
- 离开网络
- 网络状态监控

## 🔧 配置说明

### 获取 API Token

1. 访问 [ZeroTier Central](https://my.zerotier.com)
2. 登录账户
3. 进入 **Account** 页面
4. 点击 **Create API Token**
5. 复制生成的 Token

### 获取网络 ID

1. 在 ZeroTier Central 中创建网络
2. 或使用现有网络
3. 复制 16 位网络 ID (如: `8056c2e21c000001`)

## 📝 代码示例

### 基础用法

```go
package main

import (
    "fmt"
    "github.com/package-register/go-toolkit/zerotier"
)

func main() {
    // 本地客户端
    client := zerotier.NewClient()
    
    // 获取节点状态
    status, err := client.Status()
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("节点地址: %s\n", status.Address)
    fmt.Printf("在线状态: %v\n", status.Online)
    
    // 云端客户端
    central := zerotier.NewCentral("your_api_token")
    
    // 获取网络列表
    networks, err := central.Networks().List()
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("网络数量: %d\n", len(networks))
}
```

### 网络操作

```go
// 加入网络
network, err := client.Networks().Join("8056c2e21c000001")
if err != nil {
    panic(err)
}

fmt.Printf("已加入网络: %s\n", network.Name)

// 离开网络
err = client.Networks().Leave("8056c2e21c000001")
if err != nil {
    panic(err)
}

fmt.Println("已离开网络")
```

### 成员管理

```go
// 获取网络成员
members, err := central.Networks().Members("8056c2e21c000001").List()
if err != nil {
    panic(err)
}

// 授权成员
err = central.Networks().Members("8056c2e21c000001").Authorize("1234567890")
if err != nil {
    panic(err)
}

// 取消授权
err = central.Networks().Members("8056c2e21c000001").Deauthorize("1234567890")
if err != nil {
    panic(err)
}
```

## 🔍 故障排除

### 常见问题

1. **连接失败**
   - 检查 ZeroTier 服务是否运行
   - 确认本地 API 端口 (9993) 可访问

2. **API 认证失败**
   - 验证 API Token 是否正确
   - 检查 Token 权限范围

3. **网络操作失败**
   - 确认网络 ID 格式正确 (16 位十六进制)
   - 检查网络是否存在且可访问

### 调试技巧

```go
// 启用详细日志
client := zerotier.NewClient(
    zerotier.WithHost("localhost"),
    zerotier.WithPort(9993),
)

// 检查连接状态
status, err := client.Status()
if err != nil {
    fmt.Printf("连接错误: %v\n", err)
    return
}

fmt.Printf("服务版本: %s\n", status.Version)
```

## 📚 更多资源

- [ZeroTier 官方文档](https://www.zerotier.com/documentation/)
- [ZeroTier API 参考](https://docs.zerotier.com/rest/v1)
- [Go Toolkit 项目](https://github.com/package-register/go-toolkit)

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来改进这个演示！

---

💡 **提示**: 在生产环境中使用时，请妥善保管 API Token，避免硬编码在代码中。
