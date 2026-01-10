package util

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// SseConnection 表示一个 SSE 连接
type SseConnection struct {
	writer  io.Writer
	flusher http.Flusher
	ctx     context.Context
	cancel  context.CancelFunc
}

// SseEmitterManager 管理所有 SSE 连接
// 结构：userId -> token -> SseConnection
type SseEmitterManager struct {
	// userId -> token -> SseConnection 的映射
	connections map[int64]map[string]*SseConnection
	mu          sync.RWMutex
	// 心跳检测停止通道
	stopChan chan struct{}
	// 心跳检测完成通道
	doneChan chan struct{}
}

var globalSseManager *SseEmitterManager
var sseOnce sync.Once

// GetSseEmitterManager 获取全局 SSE 连接管理器（单例模式）
func GetSseEmitterManager() *SseEmitterManager {
	sseOnce.Do(func() {
		globalSseManager = &SseEmitterManager{
			connections: make(map[int64]map[string]*SseConnection),
			stopChan:    make(chan struct{}),
			doneChan:    make(chan struct{}),
		}
		// 启动心跳检测
		go globalSseManager.startHeartbeat()
	})
	return globalSseManager
}

// Connect 建立与指定用户的 SSE 连接
// userId: 用户的唯一标识符，用于区分不同用户的连接
// token: 用户的唯一令牌，用于识别具体的连接
// w: http.ResponseWriter
// flusher: http.Flusher
func (m *SseEmitterManager) Connect(userId int64, token string, w io.Writer, flusher http.Flusher) *SseConnection {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 获取或创建当前用户的连接映射表
	userConnections, exists := m.connections[userId]
	if !exists {
		userConnections = make(map[string]*SseConnection)
		m.connections[userId] = userConnections
	}

	// 关闭已存在的连接，防止超过最大连接数
	if oldConn, exists := userConnections[token]; exists {
		if oldConn.cancel != nil {
			oldConn.cancel()
		}
	}

	// 创建新的 context 和 cancel
	ctx, cancel := context.WithCancel(context.Background())

	// 创建新的 SSE 连接
	conn := &SseConnection{
		writer:  w,
		flusher: flusher,
		ctx:     ctx,
		cancel:  cancel,
	}

	// 设置响应头
	if writer, ok := w.(http.ResponseWriter); ok {
		writer.Header().Set("Content-Type", "text/event-stream")
		writer.Header().Set("Cache-Control", "no-cache")
		writer.Header().Set("Connection", "keep-alive")
		writer.Header().Set("X-Accel-Buffering", "no")
	}

	// 发送连接成功消息
	_ = conn.SendComment("connected")

	// 注册到映射表
	userConnections[token] = conn

	// 监听连接关闭
	go m.monitorConnection(userId, token, conn)

	return conn
}

// Disconnect 断开指定用户的 SSE 连接
func (m *SseEmitterManager) Disconnect(userId int64, token string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	userConnections, exists := m.connections[userId]
	if !exists {
		return
	}

	conn, exists := userConnections[token]
	if !exists {
		return
	}

	// 发送断开连接消息
	_ = conn.SendComment("disconnected")

	// 取消 context
	if conn.cancel != nil {
		conn.cancel()
	}

	// 从映射表中移除
	delete(userConnections, token)
	if len(userConnections) == 0 {
		delete(m.connections, userId)
	}
}

// SendMessage 向指定用户发送消息
func (m *SseEmitterManager) SendMessage(userId int64, message string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	userConnections, exists := m.connections[userId]
	if !exists {
		return
	}

	var toRemove []string
	for token, conn := range userConnections {
		if err := conn.Send(message); err != nil {
			logx.Errorf("发送 SSE 消息失败 userId=%d token=%s: %v", userId, token, err)
			toRemove = append(toRemove, token)
		}
	}

	// 移除失败的连接
	if len(toRemove) > 0 {
		m.mu.RUnlock()
		m.mu.Lock()
		for _, token := range toRemove {
			if conn, exists := userConnections[token]; exists {
				if conn.cancel != nil {
					conn.cancel()
				}
				delete(userConnections, token)
			}
		}
		if len(userConnections) == 0 {
			delete(m.connections, userId)
		}
		m.mu.Unlock()
		m.mu.RLock()
	}
}

// SendMessageToAll 向所有用户发送消息
func (m *SseEmitterManager) SendMessageToAll(message string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for userId := range m.connections {
		m.SendMessage(userId, message)
	}
}

// SendComment 发送注释消息（用于心跳和连接状态）
func (c *SseConnection) SendComment(comment string) error {
	select {
	case <-c.ctx.Done():
		return c.ctx.Err()
	default:
	}

	_, err := fmt.Fprintf(c.writer, ": %s\n\n", comment)
	if err != nil {
		return err
	}
	if c.flusher != nil {
		c.flusher.Flush()
	}
	return nil
}

// Send 发送消息
func (c *SseConnection) Send(data string) error {
	select {
	case <-c.ctx.Done():
		return c.ctx.Err()
	default:
	}

	_, err := fmt.Fprintf(c.writer, "event: message\ndata: %s\n\n", data)
	if err != nil {
		return err
	}
	if c.flusher != nil {
		c.flusher.Flush()
	}
	return nil
}

// Done 返回连接关闭通道
func (c *SseConnection) Done() <-chan struct{} {
	return c.ctx.Done()
}

// monitorConnection 监听连接状态
func (m *SseEmitterManager) monitorConnection(userId int64, token string, conn *SseConnection) {
	<-conn.Done()
	// 连接已关闭，从映射表中移除
	m.mu.Lock()
	defer m.mu.Unlock()

	userConnections, exists := m.connections[userId]
	if !exists {
		return
	}

	if _, exists := userConnections[token]; exists {
		delete(userConnections, token)
		if len(userConnections) == 0 {
			delete(m.connections, userId)
		}
	}
}

// startHeartbeat 启动心跳检测
func (m *SseEmitterManager) startHeartbeat() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	defer close(m.doneChan)

	for {
		select {
		case <-ticker.C:
			m.sendHeartbeat()
		case <-m.stopChan:
			return
		}
	}
}

// sendHeartbeat 发送心跳
func (m *SseEmitterManager) sendHeartbeat() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var toRemove []struct {
		userId int64
		token  string
	}

	for userId, userConnections := range m.connections {
		for token, conn := range userConnections {
			if err := conn.SendComment("heartbeat"); err != nil {
				// 发送失败，标记为需要移除
				toRemove = append(toRemove, struct {
					userId int64
					token  string
				}{userId: userId, token: token})
			}
		}
	}

	// 移除失败的连接
	if len(toRemove) > 0 {
		m.mu.RUnlock()
		m.mu.Lock()
		for _, item := range toRemove {
			m.Disconnect(item.userId, item.token)
		}
		m.mu.Unlock()
		m.mu.RLock()
	}
}

// Stop 停止 SSE 管理器
func (m *SseEmitterManager) Stop() {
	close(m.stopChan)
	<-m.doneChan

	m.mu.Lock()
	defer m.mu.Unlock()

	// 关闭所有连接
	for userId, userConnections := range m.connections {
		for _, conn := range userConnections {
			if conn.cancel != nil {
				conn.cancel()
			}
		}
		delete(m.connections, userId)
	}
}
