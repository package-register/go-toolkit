package main

import (
	"fmt"
	"time"

	"github.com/package-register/go-toolkit/zerotier"
)

func main() {
	fmt.Println("🌐 ZeroTier SDK Demo - Go Toolkit")
	fmt.Println("================================")

	// 本地节点管理演示
	demonstrateLocalClient()

	// 云端管理演示 (需要 API Token)
	demonstrateCentralClient()
}

// 本地节点管理演示
func demonstrateLocalClient() {
	fmt.Println("\n📱 本地节点管理演示")
	fmt.Println("------------------")

	// 创建本地客户端
	client := zerotier.NewClient()

	// 获取节点状态
	fmt.Println("🔍 获取节点状态...")
	status, err := client.Status()
	if err != nil {
		fmt.Printf("❌ 获取状态失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 节点信息:\n")
	fmt.Printf("   地址: %s\n", status.Address)
	fmt.Printf("   在线: %v\n", status.Online)
	fmt.Printf("   版本: %s\n", status.Version)

	// 获取网络列表
	fmt.Println("\n🌐 获取网络列表...")
	networks, err := client.Networks().List()
	if err != nil {
		fmt.Printf("❌ 获取网络列表失败: %v\n", err)
		return
	}

	if len(networks) == 0 {
		fmt.Println("ℹ️  未加入任何网络")
	} else {
		fmt.Printf("✅ 已加入 %d 个网络:\n", len(networks))
		for i, network := range networks {
			fmt.Printf("   %d. %s (%s)\n", i+1, network.Name, network.ID)
		}
	}

	// 获取节点列表
	fmt.Println("\n👥 获取节点列表...")
	peers, err := client.Peers().List()
	if err != nil {
		fmt.Printf("❌ 获取节点列表失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 发现 %d 个节点:\n", len(peers))
	for i, peer := range peers {
		if i >= 5 { // 只显示前5个
			fmt.Printf("   ... (还有 %d 个节点)\n", len(peers)-5)
			break
		}
		fmt.Printf("   %d. %s - %s\n", i+1, peer.Address, peer.Role)
	}
}

// 云端管理演示
func demonstrateCentralClient() {
	fmt.Println("\n☁️  云端管理演示")
	fmt.Println("----------------")

	// 注意: 这里需要真实的 API Token
	// 在实际使用中，请从 https://my.zerotier.com 获取
	apiToken := "your_api_token_here"

	if apiToken == "your_api_token_here" {
		fmt.Println("⚠️  需要设置真实的 API Token")
		fmt.Println("   1. 访问 https://my.zerotier.com")
		fmt.Println("   2. Account → Create API Token")
		fmt.Println("   3. 替换代码中的 apiToken 变量")
		return
	}

	// 创建云端客户端
	central := zerotier.NewCentral(apiToken)

	// 获取网络列表
	fmt.Println("🌐 获取云端网络列表...")
	networks, err := central.Networks().List()
	if err != nil {
		fmt.Printf("❌ 获取网络列表失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 云端网络数量: %d\n", len(networks))
	for i, network := range networks {
		if i >= 3 { // 只显示前3个
			fmt.Printf("   ... (还有 %d 个网络)\n", len(networks)-3)
			break
		}
		fmt.Printf("   %d. %s (%s)\n", i+1, network.Config.Name, network.ID)
	}

	// 演示网络成员管理 (需要真实的网络ID)
	networkID := "your_network_id_here"
	if networkID != "your_network_id_here" {
		fmt.Println("\n👥 网络成员管理演示...")

		// 获取网络成员
		members, err := central.Networks().Members(networkID).List()
		if err != nil {
			fmt.Printf("❌ 获取成员列表失败: %v\n", err)
			return
		}

		fmt.Printf("✅ 网络成员数量: %d\n", len(members))
		for i, member := range members {
			if i >= 3 {
				fmt.Printf("   ... (还有 %d 个成员)\n", len(members)-3)
				break
			}
			// 根据最后在线时间判断是否在线
			status := "离线"
			if member.LastOnline > 0 && time.Now().Unix()-member.LastOnline < 300 {
				status = "在线"
			}
			fmt.Printf("   %d. %s - %s\n", i+1, member.NodeID, status)
		}
	}
}

// 演示网络操作
func demonstrateNetworkOperations() {
	fmt.Println("\n🔧 网络操作演示")
	fmt.Println("----------------")

	client := zerotier.NewClient()

	// 注意: 这里需要真实的网络ID
	networkID := "your_network_id_here"

	if networkID == "your_network_id_here" {
		fmt.Println("⚠️  需要设置真实的网络ID进行网络操作演示")
		return
	}

	// 加入网络
	fmt.Printf("🔗 加入网络 %s...\n", networkID)
	_, err := client.Networks().Join(networkID)
	if err != nil {
		fmt.Printf("❌ 加入网络失败: %v\n", err)
		return
	}
	fmt.Println("✅ 成功加入网络")

	// 等待一段时间
	time.Sleep(2 * time.Second)

	// 离开网络
	fmt.Printf("🔓 离开网络 %s...\n", networkID)
	err = client.Networks().Leave(networkID)
	if err != nil {
		fmt.Printf("❌ 离开网络失败: %v\n", err)
		return
	}
	fmt.Println("✅ 成功离开网络")
}
