package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/package-register/go-toolkit/zerotier/central"
	"github.com/package-register/go-toolkit/zerotier/client"
)

// Config represents MCP configuration
type Config struct {
	LocalToken   string `json:"local_token,omitempty"`
	CentralToken string `json:"central_token,omitempty"`
}

// LoadConfig loads configuration from environment and default locations
func LoadConfig() *Config {
	config := &Config{}

	// Try to read from environment variables first
	if token := os.Getenv("ZEROTIER_LOCAL_TOKEN"); token != "" {
		config.LocalToken = token
	} else {
		// Try to read from default location
		if token := readLocalToken(); token != "" {
			config.LocalToken = token
		}
	}

	if token := os.Getenv("ZEROTIER_CENTRAL_TOKEN"); token != "" {
		config.CentralToken = token
	}

	return config
}

// readLocalToken attempts to read the local ZeroTier token from default locations
func readLocalToken() string {
	homeDir, _ := os.UserHomeDir()

	// Try different OS-specific locations
	locations := []string{
		filepath.Join(homeDir, "Library/Application Support/ZeroTier/authtoken.secret"), // macOS
		"/var/lib/zerotier-one/authtoken.secret",                                        // Linux
		"C:\\ProgramData\\ZeroTier\\One\\authtoken.secret",                              // Windows (won't work on Unix)
	}

	for _, location := range locations {
		if data, err := os.ReadFile(location); err == nil {
			return string(data)
		}
	}

	return ""
}

// FormatNetworkStatus formats network status for better readability
func FormatNetworkStatus(networks []client.Network) string {
	if len(networks) == 0 {
		return "ℹ️ 暂未加入任何网络"
	}

	var result string
	result += fmt.Sprintf("🌐 已加入 %d 个网络:\n\n", len(networks))

	for i, network := range networks {
		status := "🔴 离线"
		if network.Status == "OK" {
			status = "🟢 在线"
		}

		result += fmt.Sprintf("%d. [%s] %s\n", i+1, network.ID, network.Name)
		result += fmt.Sprintf("   状态: %s\n", status)

		if len(network.AssignedAddresses) > 0 {
			result += fmt.Sprintf("   IP: %s\n", formatIPAddresses(network.AssignedAddresses))
		}

		result += "\n"
	}

	return result
}

// FormatMemberStatus formats member status with visual indicators
func FormatMemberStatus(members []central.Member) string {
	if len(members) == 0 {
		return "ℹ️ 暂无成员"
	}

	var result string
	result += fmt.Sprintf("👥 成员列表 (%d 个):\n\n", len(members))

	for i, member := range members {
		status := "❌ 未授权"
		if member.Config != nil && member.Config.Authorized {
			status = "✅ 已授权"
		}

		result += fmt.Sprintf("%d. %s [%s]\n", i+1, member.NodeID, status)

		if member.Name != "" {
			result += fmt.Sprintf("   🏷️ 名称: %s\n", member.Name)
		}

		if member.Config != nil && len(member.Config.IPAssignments) > 0 {
			result += fmt.Sprintf("   📍 IP: %s\n", formatIPAddresses(member.Config.IPAssignments))
		}

		// Show last online time
		if member.LastOnline > 0 {
			result += fmt.Sprintf("   🕐 最后在线: %s\n", formatTimestamp(member.LastOnline))
		}

		result += "\n"
	}

	return result
}

// formatIPAddresses formats IP addresses for display
func formatIPAddresses(addresses []string) string {
	if len(addresses) == 0 {
		return "无"
	}

	if len(addresses) == 1 {
		return addresses[0]
	}

	return fmt.Sprintf("%s (共%d个)", addresses[0], len(addresses))
}

// formatTimestamp formats Unix timestamp to readable format
func formatTimestamp(timestamp int64) string {
	return fmt.Sprintf("%d秒前", time.Now().Unix()-timestamp)
}

// ValidateNetworkID validates ZeroTier network ID format
func ValidateNetworkID(networkID string) error {
	if len(networkID) != 16 {
		return fmt.Errorf("网络ID必须是16位十六进制字符串")
	}

	// Check if all characters are valid hex digits
	for _, char := range networkID {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
			return fmt.Errorf("网络ID只能包含十六进制字符")
		}
	}

	return nil
}

// ValidateIPAddress validates IP address format
func ValidateIPAddress(ip string) error {
	// Basic IP validation - could be enhanced with net.ParseIP
	if ip == "" {
		return fmt.Errorf("IP地址不能为空")
	}

	// Simple check for IPv4 format
	parts := len(strings.Split(ip, "."))
	if parts != 4 {
		return fmt.Errorf("无效的IPv4地址格式")
	}

	return nil
}

// CreateSuccessResponse creates a standardized success response
func CreateSuccessResponse(message string, data interface{}) ToolResult {
	result := ToolResult{
		Success: true,
		Data:    message,
	}

	if data != nil {
		if jsonBytes, err := json.Marshal(data); err == nil {
			result.Data = string(jsonBytes)
		}
	}

	return result
}

// CreateErrorResponse creates a standardized error response
func CreateErrorResponse(format string, args ...interface{}) ToolResult {
	return ToolResult{
		Success: false,
		Error:   fmt.Sprintf(format, args...),
	}
}
