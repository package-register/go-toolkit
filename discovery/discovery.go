package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	multicastAddr = "239.0.0.1:9999"
	announceIntv  = 15 * time.Second
	timeout       = 25 * time.Second // Increased timeout to be greater than announceIntv
)

// Logger 接口定义了日志方法
type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

// StdLogger 是 Logger 接口的标准实现
type StdLogger struct{}

func (l *StdLogger) Info(msg string, args ...any) {
	log.Printf("[INFO] "+msg, args...)
}

func (l *StdLogger) Error(msg string, args ...any) {
	log.Printf("[ERROR] "+msg, args...)
}

// MessageEnvelope 定义了在网络中传输的消息结构
type MessageEnvelope struct {
	FromUUID string          `json:"fromUuid"`
	SendType string          `json:"sendType"`          // "announce" | "spec" | "response"
	SendTo   string          `json:"sendTo,omitempty"`  // 单播目标
	Command  string          `json:"command"`           // 命令
	TaskID   string          `json:"taskId"`            // 任务 ID
	Payload  json.RawMessage `json:"payload,omitempty"` // 附带参数
}

// CommandHandler 定义了处理接收到命令的函数签名
type CommandHandler func(from net.Addr, env MessageEnvelope)

// Device 结构体表示发现到的设备信息
type Device struct {
	UUID     string
	Name     string
	IP       string
	Port     int
	Version  string
	LastSeen time.Time
}

// Discovery 结构体封装了设备发现和通信的逻辑
type Discovery struct {
	uuid          string
	name          string
	ip            string
	port          int
	version       string
	logger        Logger
	ctx           context.Context
	cancel        context.CancelFunc
	mu            sync.RWMutex
	handlers      map[string]CommandHandler
	devices       map[string]*Device
	pending       map[string]chan MessageEnvelope
	multicastConns []*net.UDPConn // Changed to slice for multiple connections
	unicastConn   *net.UDPConn
}

// NewDiscovery 创建一个新的 Discovery 实例
func NewDiscovery(name, ver string, logger Logger) *Discovery {
	ip, _ := getLocalIP() // This will need to be revisited for multi-NIC
	port, _ := getAvailablePort()
	ctx, cancel := context.WithCancel(context.Background())
	return &Discovery{
		uuid:     uuid.New().String(),
		name:     name,
		ip:       ip,
		port:     port,
		version:  ver,
		logger:   logger,
		ctx:      ctx,
		cancel:   cancel,
		handlers: make(map[string]CommandHandler),
		devices:  make(map[string]*Device),
		pending:  make(map[string]chan MessageEnvelope),
	}
}

// RegisterHandler 注册命令处理器
func (d *Discovery) RegisterHandler(cmd string, handler CommandHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[cmd] = handler
}

// Start 启动设备发现服务
func (d *Discovery) Start() error {
	maddr, err := net.ResolveUDPAddr("udp", multicastAddr)
	if err != nil {
		return fmt.Errorf("resolve multicast address failed: %w", err)
	}

	// Get all active interfaces with IPv4 addresses
	interfaces, err := getActiveInterfaces()
	if err != nil {
		return fmt.Errorf("get active interfaces failed: %w", err)
	}

	// Create a multicast connection for each active interface
	for _, iface := range interfaces {
		mc, err := net.ListenMulticastUDP("udp", &iface, maddr)
		if err != nil {
			d.logger.Error("listen multicast UDP on interface %s failed: %v", iface.Name, err)
			continue // Try next interface
		}
		d.multicastConns = append(d.multicastConns, mc)
		go d.listenMulticast(mc) // Start a listener for each connection
	}

	if len(d.multicastConns) == 0 {
		return fmt.Errorf("no active multicast interfaces found to listen on")
	}

	uc, err := net.ListenUDP("udp", nil)
	if err != nil {
		return fmt.Errorf("listen unicast UDP failed: %w", err)
	}
	d.unicastConn = uc

	go d.sendMulticastAnnounce()
	go d.cleanupDevices()

	d.logger.Info("Discovery 启动: %s (%s:%d)", d.name, d.ip, d.port)
	return nil
}

// Stop 停止设备发现服务
func (d *Discovery) Stop() {
	d.cancel()
	for _, conn := range d.multicastConns {
		_ = conn.Close()
	}
	if d.unicastConn != nil {
		_ = d.unicastConn.Close()
	}
	d.logger.Info("Discovery 已停止")
}

// Send 发送消息
func (d *Discovery) Send(env MessageEnvelope) error {
	env.FromUUID = d.uuid
	data, err := json.Marshal(env)
	if err != nil {
		return fmt.Errorf("marshal message envelope failed: %w", err)
	}

	switch env.SendType {
	case "announce":
		return d.sendMulticast(data)
	case "spec", "response":
		return d.sendUnicast(env.SendTo, data)
	default:
		return fmt.Errorf("未知 sendType: %s", env.SendType)
	}
}

// RequestResponse 发送请求并等待响应
func (d *Discovery) RequestResponse(env MessageEnvelope, timeout time.Duration) (MessageEnvelope, error) {
	if env.TaskID == "" {
		env.TaskID = uuid.New().String()
	}
	ch := make(chan MessageEnvelope, 1)

	d.mu.Lock()
	d.pending[env.TaskID] = ch
	d.mu.Unlock()

	if err := d.Send(env); err != nil {
		d.mu.Lock()
		delete(d.pending, env.TaskID)
		d.mu.Unlock()
		return MessageEnvelope{}, fmt.Errorf("send request failed: %w", err)
	}

	select {
	case resp := <-ch:
		return resp, nil
	case <-time.After(timeout):
		d.mu.Lock()
		delete(d.pending, env.TaskID)
		d.mu.Unlock()
		return MessageEnvelope{}, fmt.Errorf("任务超时: %s", env.TaskID)
	}
}

// GetDevices 返回当前发现到的设备列表
func (d *Discovery) GetDevices() []*Device {
	d.mu.RLock()
	defer d.mu.RUnlock()
	devices := make([]*Device, 0, len(d.devices))
	for _, dev := range d.devices {
		devices = append(devices, dev)
	}
	return devices
}

func (d *Discovery) listenMulticast(conn *net.UDPConn) {
	buf := make([]byte, 2048)
	for {
		select {
		case <-d.ctx.Done():
			return
		default:
			// 设置读取超时，防止阻塞
			_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
			n, src, err := conn.ReadFromUDP(buf)
			if err != nil {
				// 忽略超时错误
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				d.logger.Error("读取组播消息失败: %v", err)
				continue
			}
			var env MessageEnvelope
			if err := json.Unmarshal(buf[:n], &env); err != nil {
				d.logger.Error("解析组播消息失败: %v", err)
				continue
			}
			if env.FromUUID == d.uuid { // 忽略自己的消息
				continue
			}

			d.processReceivedMessage(src, env)
		}
	}
}

func (d *Discovery) processReceivedMessage(from net.Addr, env MessageEnvelope) {
	if env.Command == "announce" {
		var info map[string]any
		if err := json.Unmarshal(env.Payload, &info); err != nil {
			d.logger.Error("解析 announce 消息 payload 失败: %v", err)
			return
		}
		d.mu.Lock()
		d.devices[env.FromUUID] = &Device{
			UUID:     env.FromUUID,
			Name:     fmt.Sprint(info["name"]),
			IP:       fmt.Sprint(info["ip"]),
			Port:     int(info["port"].(float64)),
			Version:  fmt.Sprint(info["version"]),
			LastSeen: time.Now(),
		}
		d.mu.Unlock()
		d.logger.Info("发现新设备: %s (%s:%d)", info["name"], info["ip"], int(info["port"].(float64)))
	}

	d.mu.RLock()
	ch, ok := d.pending[env.TaskID]
	d.mu.RUnlock()
	if ok {
		ch <- env
		d.mu.Lock()
		delete(d.pending, env.TaskID)
		d.mu.Unlock()
		return
	}

	d.mu.RLock()
	handler, ok := d.handlers[env.Command]
	d.mu.RUnlock()
	if ok {
		go handler(from, env)
	} else {
		d.logger.Info("未注册命令处理器: %s", env.Command)
	}
}

func (d *Discovery) sendMulticastAnnounce() {
	ticker := time.NewTicker(announceIntv)
	defer ticker.Stop()
	for {
		select {
		case <-d.ctx.Done():
			return
		case <-ticker.C:
			env := MessageEnvelope{
				SendType: "announce",
				Command:  "announce",
				TaskID:   uuid.New().String(),
				Payload:  mustJSON(map[string]any{"name": d.name, "version": d.version, "ip": d.ip, "port": d.port, "uuid": d.uuid}),
			}
			if err := d.Send(env); err != nil {
				d.logger.Error("发送 announce 消息失败: %v", err)
			}
		}
	}
}

func (d *Discovery) cleanupDevices() {
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()
	for {
		select {
		case <-d.ctx.Done():
			return
		case <-ticker.C:
			d.mu.Lock()
			now := time.Now()
			for k, dev := range d.devices {
				if now.Sub(dev.LastSeen) > timeout {
					delete(d.devices, k)
					d.logger.Info("设备过期移除: %s", dev.UUID)
				}
			}
			d.mu.Unlock()
		}
	}
}

func (d *Discovery) sendMulticast(data []byte) error {
	maddr, err := net.ResolveUDPAddr("udp", multicastAddr)
	if err != nil {
		return fmt.Errorf("resolve multicast address failed: %w", err)
	}

	var lastErr error
	// Send on all available multicast connections
	for _, conn := range d.multicastConns {
		// Set a write deadline for each send operation
		_ = conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
		_, err := conn.WriteToUDP(data, maddr)
		if err != nil {
			lastErr = fmt.Errorf("send multicast on %s failed: %w", conn.LocalAddr().String(), err)
			d.logger.Error(lastErr.Error())
		}
	}
	return lastErr // Return the last error encountered, or nil if all succeeded
}

func (d *Discovery) sendUnicast(addr string, data []byte) error {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return fmt.Errorf("resolve unicast address failed: %w", err)
	}
	_ = d.unicastConn.SetWriteDeadline(time.Now().Add(2 * time.Second))
	_, err = d.unicastConn.WriteToUDP(data, udpAddr)
	return err
}

// getActiveInterfaces returns a list of active network interfaces with at least one IPv4 address.
func getActiveInterfaces() ([]net.Interface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("get network interfaces failed: %w", err)
	}

	var activeInterfaces []net.Interface
	for _, iface := range interfaces {
		// Check if the interface is up, not a loopback, and supports multicast
		if iface.Flags&net.FlagUp != 0 &&
			iface.Flags&net.FlagLoopback == 0 &&
			iface.Flags&net.FlagMulticast != 0 {

			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}

			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil { // Only consider IPv4 addresses
						activeInterfaces = append(activeInterfaces, iface)
						break // Found an IPv4 address, move to next interface
					}
				}
			}
		}
	}
	return activeInterfaces, nil
}

func getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", fmt.Errorf("dial UDP for local IP failed: %w", err)
	}

	defer func() {
		_ = conn.Close()
	}()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func getAvailablePort() (int, error) {
	addr, _ := net.ResolveTCPAddr("tcp", "localhost:0")
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, fmt.Errorf("listen TCP for available port failed: %w", err)
	}

	defer func() {
		_ = l.Close()
	}()

	return l.Addr().(*net.TCPAddr).Port, nil
}

func mustJSON(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}
