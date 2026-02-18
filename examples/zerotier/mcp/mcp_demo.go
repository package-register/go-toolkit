package main

import (
	"fmt"
	"log"

	"github.com/package-register/go-toolkit/zerotier/mcp"
)

func main() {
	fmt.Println("🌐 ZeroTier MCP Server Demo")
	fmt.Println("============================")

	// Load configuration
	config := mcp.LoadConfig()
	
	// Create MCP server
	server := mcp.NewMcpServer()
	
	// Set central token if available
	if config.CentralToken != "" {
		server = server.WithCentralToken(config.CentralToken)
		fmt.Println("✅ 已配置 Central API Token")
	} else {
		fmt.Println("⚠️ 未配置 Central API Token，云端功能将不可用")
	}

	if config.LocalToken != "" {
		fmt.Println("✅ 已读取本地 Token")
	} else {
		fmt.Println("⚠️ 未找到本地 Token，请确保 ZeroTier 服务正在运行")
	}

	// Demo local tools
	demonstrateLocalTools(server)

	// Demo central tools (if token is available)
	if config.CentralToken != "" {
		demonstrateCentralTools(server)
	}

	// Demo MCP tool execution
	demonstrateMcpExecution(server)
}

func demonstrateLocalTools(server *mcp.McpServer) {
	fmt.Println("\n📱 本地工具演示")
	fmt.Println("----------------")

	tools := server.GetTools()
	localTools := []string{"zt_status", "zt_networks", "zt_peers"}

	for _, toolName := range localTools {
		for _, tool := range tools {
			if tool.Name == toolName {
				fmt.Printf("\n🔧 执行工具: %s\n", tool.Description)
				
				var result mcp.ToolResult
				switch toolName {
				case "zt_status":
					result = server.ExecuteTool("zt_status", nil)
				case "zt_networks":
					result = server.ExecuteTool("zt_networks", nil)
				case "zt_peers":
					result = server.ExecuteTool("zt_peers", nil)
				}

				if result.Success {
					fmt.Printf("✅ 成功:\n%s\n", result.Data)
				} else {
					fmt.Printf("❌ 失败: %s\n", result.Error)
				}
				break
			}
		}
	}
}

func demonstrateCentralTools(server *mcp.McpServer) {
	fmt.Println("\n☁️ 云端工具演示")
	fmt.Println("----------------")

	tools := server.GetTools()
	centralTools := []string{"zt_central_networks"}

	for _, toolName := range centralTools {
		for _, tool := range tools {
			if tool.Name == toolName {
				fmt.Printf("\n🔧 执行工具: %s\n", tool.Description)
				
				result := server.ExecuteTool(toolName, nil)
				
				if result.Success {
					fmt.Printf("✅ 成功:\n%s\n", result.Data)
				} else {
					fmt.Printf("❌ 失败: %s\n", result.Error)
				}
				break
			}
		}
	}
}

func demonstrateMcpExecution(server *mcp.McpServer) {
	fmt.Println("\n🤖 MCP 执行演示")
	fmt.Println("----------------")

	// Demo joining a network (this will fail without a real network ID, but shows the interface)
	fmt.Println("\n🔧 演示网络加入 (使用测试网络ID):")
	params := map[string]interface{}{
		"network_id": "1234567890abcdef", // Test network ID
	}
	
	result := server.ExecuteTool("zt_join", params)
	if result.Success {
		fmt.Printf("✅ 成功: %s\n", result.Data)
	} else {
		fmt.Printf("❌ 失败 (预期): %s\n", result.Error)
	}

	// Demo member authorization with IP
	fmt.Println("\n🔧 演示成员授权并指定IP:")
	params = map[string]interface{}{
		"network_id": "1234567890abcdef",
		"member_id": "abcdef1234567890",
		"ip_address": "10.147.20.100",
		"name": "测试服务器",
	}
	
	result = server.ExecuteTool("zt_central_authorize_with_ip", params)
	if result.Success {
		fmt.Printf("✅ 成功: %s\n", result.Data)
	} else {
		fmt.Printf("❌ 失败 (预期): %s\n", result.Error)
	}

	// Show all available tools
	fmt.Println("\n📋 所有可用工具:")
	tools := server.GetTools()
	for i, tool := range tools {
		fmt.Printf("%d. %s - %s\n", i+1, tool.Name, tool.Description)
		if len(tool.Parameters) > 0 {
			fmt.Printf("   参数: ")
			for j, param := range tool.Parameters {
				if j > 0 {
					fmt.Printf(", ")
				}
				required := ""
				if param.Required {
					required = "*"
				}
				fmt.Printf("%s%s", param.Name, required)
			}
			fmt.Println()
		}
	}
}

func init() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
