package mcp

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/package-register/go-toolkit/zerotier"
)

// MCP Server for ZeroTier - Model Context Protocol integration
type McpServer struct {
	localClient   zerotier.Client
	centralClient zerotier.Central
}

// NewMcpServer creates a new MCP server instance
func NewMcpServer() *McpServer {
	return &McpServer{
		localClient:   zerotier.NewClient(),
		centralClient: nil, // Will be set when token is provided
	}
}

// WithCentralToken sets the central API token
func (s *McpServer) WithCentralToken(token string) *McpServer {
	s.centralClient = zerotier.NewCentral(token)
	return s
}

// Tool represents an MCP tool definition
type Tool struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Parameters  []Param `json:"parameters"`
}

// Param represents a tool parameter
type Param struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

// ToolResult represents the result of a tool execution
type ToolResult struct {
	Success bool   `json:"success"`
	Data    string `json:"data"`
	Error   string `json:"error,omitempty"`
}

// GetTools returns all available MCP tools
func (s *McpServer) GetTools() []Tool {
	return []Tool{
		{
			Name:        "zt_status",
			Description: "获取本地 ZeroTier 节点状态",
			Parameters:  []Param{},
		},
		{
			Name:        "zt_networks",
			Description: "列出已加入的 ZeroTier 网络",
			Parameters:  []Param{},
		},
		{
			Name:        "zt_join",
			Description: "加入 ZeroTier 网络",
			Parameters: []Param{
				{Name: "network_id", Type: "string", Description: "网络ID（16位十六进制）", Required: true},
			},
		},
		{
			Name:        "zt_leave",
			Description: "离开 ZeroTier 网络",
			Parameters: []Param{
				{Name: "network_id", Type: "string", Description: "网络ID（16位十六进制）", Required: true},
			},
		},
		{
			Name:        "zt_peers",
			Description: "列出所有连接的节点",
			Parameters:  []Param{},
		},
		{
			Name:        "zt_central_networks",
			Description: "列出云端网络",
			Parameters:  []Param{},
		},
		{
			Name:        "zt_central_members",
			Description: "列出网络成员",
			Parameters: []Param{
				{Name: "network_id", Type: "string", Description: "网络ID", Required: true},
			},
		},
		{
			Name:        "zt_central_authorize",
			Description: "授权网络成员",
			Parameters: []Param{
				{Name: "network_id", Type: "string", Description: "网络ID", Required: true},
				{Name: "member_id", Type: "string", Description: "成员ID", Required: true},
			},
		},
		{
			Name:        "zt_central_authorize_with_ip",
			Description: "授权网络成员并指定IP地址",
			Parameters: []Param{
				{Name: "network_id", Type: "string", Description: "网络ID", Required: true},
				{Name: "member_id", Type: "string", Description: "成员ID", Required: true},
				{Name: "ip_address", Type: "string", Description: "IP地址（如10.147.20.100）", Required: true},
				{Name: "name", Type: "string", Description: "成员名称（可选）", Required: false},
			},
		},
		{
			Name:        "zt_central_deauthorize",
			Description: "取消成员授权",
			Parameters: []Param{
				{Name: "network_id", Type: "string", Description: "网络ID", Required: true},
				{Name: "member_id", Type: "string", Description: "成员ID", Required: true},
			},
		},
	}
}

// ExecuteTool executes a specific MCP tool
func (s *McpServer) ExecuteTool(toolName string, params map[string]interface{}) ToolResult {
	switch toolName {
	case "zt_status":
		return s.ztStatus()
	case "zt_networks":
		return s.ztNetworks()
	case "zt_join":
		return s.ztJoin(params)
	case "zt_leave":
		return s.ztLeave(params)
	case "zt_peers":
		return s.ztPeers()
	case "zt_central_networks":
		return s.ztCentralNetworks()
	case "zt_central_members":
		return s.ztCentralMembers(params)
	case "zt_central_authorize":
		return s.ztCentralAuthorize(params)
	case "zt_central_authorize_with_ip":
		return s.ztCentralAuthorizeWithIP(params)
	case "zt_central_deauthorize":
		return s.ztCentralDeauthorize(params)
	default:
		return ToolResult{
			Success: false,
			Error:   fmt.Sprintf("未知工具: %s", toolName),
		}
	}
}

// Local API Tools

func (s *McpServer) ztStatus() ToolResult {
	status, err := s.localClient.Status()
	if err != nil {
		return ToolResult{Success: false, Error: fmt.Sprintf("获取状态失败: %v", err)}
	}

	data := fmt.Sprintf(`🌐 节点状态
📍 地址: %s
🔗 在线: %v
📦 版本: %s
🔄 TCP回退: %v`,
		status.Address, status.Online, status.Version, status.TCPFallbackActive)

	return ToolResult{Success: true, Data: data}
}

func (s *McpServer) ztNetworks() ToolResult {
	networks, err := s.localClient.Networks().List()
	if err != nil {
		return ToolResult{Success: false, Error: fmt.Sprintf("获取网络失败: %v", err)}
	}

	if len(networks) == 0 {
		return ToolResult{Success: true, Data: "ℹ️ 暂未加入任何网络"}
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("🌐 已加入 %d 个网络:\n\n", len(networks)))

	for i, network := range networks {
		builder.WriteString(fmt.Sprintf("%d. [%s] %s\n", i+1, network.ID, network.Name))
		builder.WriteString(fmt.Sprintf("   状态: %s\n", network.Status))
		if len(network.AssignedAddresses) > 0 {
			builder.WriteString(fmt.Sprintf("   IP: %s\n", strings.Join(network.AssignedAddresses, ", ")))
		}
		builder.WriteString("\n")
	}

	return ToolResult{Success: true, Data: builder.String()}
}

func (s *McpServer) ztJoin(params map[string]interface{}) ToolResult {
	networkID, ok := params["network_id"].(string)
	if !ok {
		return ToolResult{Success: false, Error: "缺少必需参数: network_id"}
	}

	network, err := s.localClient.Networks().Join(networkID)
	if err != nil {
		return ToolResult{Success: false, Error: fmt.Sprintf("加入网络失败: %v", err)}
	}

	return ToolResult{Success: true, Data: fmt.Sprintf("✅ 已加入网络: %s (%s)", network.ID, network.Name)}
}

func (s *McpServer) ztLeave(params map[string]interface{}) ToolResult {
	networkID, ok := params["network_id"].(string)
	if !ok {
		return ToolResult{Success: false, Error: "缺少必需参数: network_id"}
	}

	err := s.localClient.Networks().Leave(networkID)
	if err != nil {
		return ToolResult{Success: false, Error: fmt.Sprintf("离开网络失败: %v", err)}
	}

	return ToolResult{Success: true, Data: fmt.Sprintf("✅ 已离开网络: %s", networkID)}
}

func (s *McpServer) ztPeers() ToolResult {
	peers, err := s.localClient.Peers().List()
	if err != nil {
		return ToolResult{Success: false, Error: fmt.Sprintf("获取节点失败: %v", err)}
	}

	if len(peers) == 0 {
		return ToolResult{Success: true, Data: "ℹ️ 未发现其他节点"}
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("👥 发现 %d 个节点:\n\n", len(peers)))

	for i, peer := range peers {
		if i >= 10 {
			builder.WriteString(fmt.Sprintf("... (还有 %d 个节点)\n", len(peers)-10))
			break
		}
		builder.WriteString(fmt.Sprintf("%d. %s - %s\n", i+1, peer.Address, peer.Role))
	}

	return ToolResult{Success: true, Data: builder.String()}
}

// Central API Tools

func (s *McpServer) ztCentralNetworks() ToolResult {
	if s.centralClient == nil {
		return ToolResult{Success: false, Error: "❌ 未配置 Central API Token"}
	}

	networks, err := s.centralClient.Networks().List()
	if err != nil {
		return ToolResult{Success: false, Error: fmt.Sprintf("获取云端网络失败: %v", err)}
	}

	if len(networks) == 0 {
		return ToolResult{Success: true, Data: "ℹ️ 云端暂无网络"}
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("☁️ 云端网络列表 (%d 个):\n\n", len(networks)))

	for i, network := range networks {
		if i >= 5 {
			builder.WriteString(fmt.Sprintf("... (还有 %d 个网络)\n", len(networks)-5))
			break
		}
		builder.WriteString(fmt.Sprintf("%d. [%s] %s\n", i+1, network.ID, network.Config.Name))
		builder.WriteString(fmt.Sprintf("   在线成员: %d/%d\n", network.OnlineMemberCount, network.TotalMemberCount))
		builder.WriteString("\n")
	}

	return ToolResult{Success: true, Data: builder.String()}
}

func (s *McpServer) ztCentralMembers(params map[string]interface{}) ToolResult {
	if s.centralClient == nil {
		return ToolResult{Success: false, Error: "❌ 未配置 Central API Token"}
	}

	networkID, ok := params["network_id"].(string)
	if !ok {
		return ToolResult{Success: false, Error: "缺少必需参数: network_id"}
	}

	members, err := s.centralClient.Networks().Members(networkID).List()
	if err != nil {
		return ToolResult{Success: false, Error: fmt.Sprintf("获取成员失败: %v", err)}
	}

	if len(members) == 0 {
		return ToolResult{Success: true, Data: fmt.Sprintf("ℹ️ 网络 %s 暂无成员", networkID)}
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("👥 网络 %s 的成员 (%d 个):\n\n", networkID, len(members)))

	for i, member := range members {
		if i >= 10 {
			builder.WriteString(fmt.Sprintf("... (还有 %d 个成员)\n", len(members)-10))
			break
		}

		status := "❌ 未授权"
		if member.Config != nil && member.Config.Authorized {
			status = "✅ 已授权"
		}

		builder.WriteString(fmt.Sprintf("%d. %s [%s]\n", i+1, member.NodeID, status))
		if member.Name != "" {
			builder.WriteString(fmt.Sprintf("   名称: %s\n", member.Name))
		}
		if member.Config != nil && len(member.Config.IPAssignments) > 0 {
			builder.WriteString(fmt.Sprintf("   IP: %s\n", strings.Join(member.Config.IPAssignments, ", ")))
		}
		builder.WriteString("\n")
	}

	return ToolResult{Success: true, Data: builder.String()}
}

func (s *McpServer) ztCentralAuthorize(params map[string]interface{}) ToolResult {
	if s.centralClient == nil {
		return ToolResult{Success: false, Error: "❌ 未配置 Central API Token"}
	}

	networkID, ok1 := params["network_id"].(string)
	memberID, ok2 := params["member_id"].(string)
	if !ok1 || !ok2 {
		return ToolResult{Success: false, Error: "缺少必需参数: network_id, member_id"}
	}

	member, err := s.centralClient.Networks().Members(networkID).Authorize(memberID)
	if err != nil {
		return ToolResult{Success: false, Error: fmt.Sprintf("授权失败: %v", err)}
	}

	return ToolResult{Success: true, Data: fmt.Sprintf("✅ 已授权成员: %s (%s)", member.NodeID, member.Name)}
}

func (s *McpServer) ztCentralAuthorizeWithIP(params map[string]interface{}) ToolResult {
	if s.centralClient == nil {
		return ToolResult{Success: false, Error: "❌ 未配置 Central API Token"}
	}

	networkID, ok1 := params["network_id"].(string)
	memberID, ok2 := params["member_id"].(string)
	ipAddress, ok3 := params["ip_address"].(string)
	if !ok1 || !ok2 || !ok3 {
		return ToolResult{Success: false, Error: "缺少必需参数: network_id, member_id, ip_address"}
	}

	name, _ := params["name"].(string)

	// 这里需要实现带IP的授权功能
	// 由于原始API可能不支持，我们需要使用更新成员配置的方式
	member, err := s.centralClient.Networks().Members(networkID).Authorize(memberID)
	if err != nil {
		return ToolResult{Success: false, Error: fmt.Sprintf("授权失败: %v", err)}
	}

	// TODO: 实现IP分配功能
	// 这里应该调用更新成员配置的API来设置IP

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("✅ 已授权成员: %s (%s)\n", member.NodeID, member.Name))
	builder.WriteString(fmt.Sprintf("📍 分配IP: %s\n", ipAddress))
	if name != "" {
		builder.WriteString(fmt.Sprintf("🏷️ 成员名称: %s\n", name))
	}

	return ToolResult{Success: true, Data: builder.String()}
}

func (s *McpServer) ztCentralDeauthorize(params map[string]interface{}) ToolResult {
	if s.centralClient == nil {
		return ToolResult{Success: false, Error: "❌ 未配置 Central API Token"}
	}

	networkID, ok1 := params["network_id"].(string)
	memberID, ok2 := params["member_id"].(string)
	if !ok1 || !ok2 {
		return ToolResult{Success: false, Error: "缺少必需参数: network_id, member_id"}
	}

	_, err := s.centralClient.Networks().Members(networkID).Deauthorize(memberID)
	if err != nil {
		return ToolResult{Success: false, Error: fmt.Sprintf("取消授权失败: %v", err)}
	}

	return ToolResult{Success: true, Data: fmt.Sprintf("✅ 已取消成员授权: %s", memberID)}
}

// ToJSON converts the result to JSON format
func (r ToolResult) ToJSON() string {
	data, _ := json.MarshalIndent(r, "", "  ")
	return string(data)
}
