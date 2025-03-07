package service

import (
	"back-end/config"
	"back-end/entity/dto"
	"back-end/entity/pojo"
	"back-end/util"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源，生产环境应该更严格
	},
}

// WebSocketMessage 定义WebSocket消息结构
type WebSocketMessage struct {
	Type      string      `json:"type"`      // 消息类型: "chat", "system", "join", "leave"
	Content   string      `json:"content"`   // 消息内容
	SenderID  uint        `json:"senderId"`  // 发送者ID
	GroupID   uint        `json:"groupId"`   // 群组ID
	MediaType string      `json:"mediaType"` // 媒体类型，如 "text", "image" 等
	MediaURL  string      `json:"mediaUrl"`  // 媒体URL
	Timestamp time.Time   `json:"timestamp"` // 发送时间
	Nickname  string      `json:"nickname"`  // 发送者昵称
	AvatarURL string      `json:"avatarUrl"` // 发送者头像
	Data      interface{} `json:"data"`      // 额外数据
}

// ClientManager 管理所有WebSocket客户端
type ClientManager struct {
	clients    map[*Client]bool          // 所有连接的客户端
	broadcast  chan []byte               // 广播消息的通道
	register   chan *Client              // 注册客户端的通道
	unregister chan *Client              // 注销客户端的通道
	groups     map[uint]map[*Client]bool // 群组到客户端的映射
	mutex      sync.Mutex                // 互斥锁，保护并发访问
}

// Client 表示一个WebSocket客户端
type Client struct {
	id           uint              // 用户ID
	socket       *websocket.Conn   // WebSocket连接
	send         chan []byte       // 发送消息的通道
	groups       map[uint]struct{} // 用户所在的群组
	isAuthorized bool              // 是否已认证
}

// 全局的客户端管理器
var Manager = ClientManager{
	clients:    make(map[*Client]bool),
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	groups:     make(map[uint]map[*Client]bool),
}

// Start 启动WebSocket管理器
func (manager *ClientManager) Start() {
	for {
		select {
		case client := <-manager.register:
			// 注册新客户端
			manager.mutex.Lock()
			manager.clients[client] = true
			manager.mutex.Unlock()

		case client := <-manager.unregister:
			// 客户端断开连接
			if _, ok := manager.clients[client]; ok {
				manager.mutex.Lock()
				delete(manager.clients, client)
				// 从所有群组中移除该客户端
				for groupID := range client.groups {
					if _, exists := manager.groups[groupID]; exists {
						delete(manager.groups[groupID], client)
						// 如果群组没有成员了，删除该群组
						if len(manager.groups[groupID]) == 0 {
							delete(manager.groups, groupID)
						}
					}
				}
				manager.mutex.Unlock()
				close(client.send)
			}
		}
	}
}

// 添加初始化消息结构体
type GroupChatInitMessage struct {
	Token string `json:"token"`
}

// HandleGroupChat 处理群组聊天的WebSocket连接
func HandleGroupChat(c *gin.Context) {
	// 升级HTTP连接为WebSocket - 无条件接受连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	// 获取群组ID
	groupIDStr := c.Param("id")
	groupID, err := strconv.ParseUint(groupIDStr, 10, 64)
	if err != nil {
		conn.WriteJSON(dto.ErrorResponse[string](400, "无效的群组ID"))
		conn.Close()
		return
	}

	// 创建未认证的客户端
	client := &Client{
		socket:       conn,
		send:         make(chan []byte, 256),
		groups:       make(map[uint]struct{}),
		isAuthorized: false,
	}

	// 等待认证消息
	_, message, err := conn.ReadMessage()
	if err != nil {
		conn.Close()
		return
	}

	var initMsg GroupChatInitMessage
	if err := json.Unmarshal(message, &initMsg); err != nil {
		conn.WriteJSON(dto.ErrorResponse[string](500, "无法解析认证消息"))
		conn.Close()
		return
	}

	// 验证token并获取用户ID
	claims, err := util.ParseToken(initMsg.Token)
	if err != nil {
		conn.WriteJSON(dto.ErrorResponse[string](500, "token验证失败"))
		conn.Close()
		return
	}
	userID := claims.UserID

	// 验证用户是否为群组成员
	var member pojo.GroupMember
	err = config.MysqlDataBase.Where("group_id = ? AND user_id = ? AND status != 3", groupID, userID).First(&member).Error
	if err != nil {
		conn.WriteJSON(dto.ErrorResponse[string](403, "您不是该群组成员或群组不存在"))
		conn.Close()
		return
	}

	// 设置客户端信息
	client.id = uint(userID)
	client.isAuthorized = true
	client.groups[uint(groupID)] = struct{}{}

	// 向管理器注册客户端
	Manager.register <- client

	// 向群组添加客户端
	Manager.mutex.Lock()
	if _, exists := Manager.groups[uint(groupID)]; !exists {
		Manager.groups[uint(groupID)] = make(map[*Client]bool)
	}
	Manager.groups[uint(groupID)][client] = true
	Manager.mutex.Unlock()

	// 发送认证成功消息
	err = conn.WriteJSON(dto.SuccessResponse("认证成功"))
	if err != nil {
		conn.Close()
		return
	}

	// 发送用户加入的系统消息
	var user pojo.User
	config.MysqlDataBase.Select("id, username, avatar").Where("id = ?", userID).First(&user)

	joinMessage := WebSocketMessage{
		Type:      "join",
		SenderID:  uint(userID),
		GroupID:   uint(groupID),
		Timestamp: time.Now(),
		Nickname:  user.Username,
		AvatarURL: user.Avatar,
		Content:   fmt.Sprintf("%s 加入了群聊", user.Username),
	}

	// 广播加入消息
	broadcastToGroup(joinMessage, uint(groupID))

	// 启动读取和写入的goroutine
	go client.readPump(uint(groupID))
	go client.writePump()

	// 使用通道阻塞，直到连接断开
	done := make(chan struct{})

	// 当连接断开时关闭通道
	go func() {
		for {
			// 尝试读取消息，如果连接断开，ReadMessage 会返回错误
			_, _, err := conn.ReadMessage()
			if err != nil {
				// 连接断开，关闭通道
				close(done)
				return
			}
			// 不做任何处理，因为 readPump 已经在处理消息
		}
	}()

	// 等待连接断开
	<-done

	// conn会在 Client.readPump 或 Client.writePump 结束时关闭
	fmt.Printf("WebSocket连接已关闭: 用户 %d, 群组 %d\n", userID, groupID)
}

// readPump 从WebSocket读取消息
func (c *Client) readPump(groupID uint) {
	defer func() {
		Manager.unregister <- c
		c.socket.Close()
	}()

	c.socket.SetReadLimit(512000)
	c.socket.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.socket.SetPongHandler(func(string) error {
		c.socket.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("WebSocket错误: %v\n", err)
			}
			break
		}

		// 检查认证状态
		if !c.isAuthorized {
			c.socket.WriteJSON(dto.ErrorResponse[string](401, "未认证的连接"))
			continue
		}

		// 解析消息
		var wsMessage WebSocketMessage
		if err := json.Unmarshal(message, &wsMessage); err != nil {
			continue
		}

		// 处理ping消息
		if wsMessage.Type == "ping" {
			// 回复pong消息
			pongMessage := WebSocketMessage{
				Type:      "pong",
				Timestamp: time.Now(),
			}

			if data, err := json.Marshal(pongMessage); err == nil {
				c.send <- data
			}

			// 更新读取截止时间
			c.socket.SetReadDeadline(time.Now().Add(60 * time.Second))
			continue
		}

		// 设置消息属性
		wsMessage.SenderID = c.id
		wsMessage.GroupID = groupID
		wsMessage.Timestamp = time.Now()

		// 查询发送者信息
		var user pojo.User
		if err := config.MysqlDataBase.Select("username, avatar").Where("id = ?", c.id).First(&user).Error; err == nil {
			wsMessage.Nickname = user.Username
			wsMessage.AvatarURL = user.Avatar
		}

		// 保存消息到数据库
		saveMessage(wsMessage)

		// 广播消息到群组
		broadcastToGroup(wsMessage, groupID)
	}
}

// writePump 向WebSocket写入消息
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.socket.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// 通道已关闭
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.socket.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 添加队列中的消息
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			// 发送ping保持连接
			c.socket.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// broadcastToGroup 向群组广播消息
func broadcastToGroup(message WebSocketMessage, groupID uint) {
	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	Manager.mutex.Lock()
	defer Manager.mutex.Unlock()

	if clients, exists := Manager.groups[groupID]; exists {
		for client := range clients {
			// 只向已认证的客户端发送消息
			if client.isAuthorized {
				select {
				case client.send <- data:
				default:
					close(client.send)
					delete(Manager.groups[groupID], client)
					delete(Manager.clients, client)
				}
			}
		}
	}
}

// saveMessage 将消息保存到数据库
func saveMessage(message WebSocketMessage) {
	// 创建消息记录
	dbMessage := pojo.Message{
		Type:      pojo.GroupMessage, // 群聊消息类型
		Content:   message.Content,
		SenderID:  message.SenderID,
		GroupID:   message.GroupID,
		IsRead:    false,
		MediaType: message.MediaType,
		MediaURL:  message.MediaURL,
	}

	// 保存到数据库
	config.MysqlDataBase.Create(&dbMessage)
}

// HandleStreamChat 处理流式聊天的WebSocket连接
func HandleStreamChat(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	// 读取初始请求
	_, message, err := ws.ReadMessage()
	if err != nil {
		return
	}

	var chatRequest util.ChatRequest
	if err := json.Unmarshal(message, &chatRequest); err != nil {
		ws.WriteJSON(map[string]interface{}{
			"error": "请求JSON有误" + err.Error(),
		})
		return
	}

	// 创建上下文，支持取消
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 获取流式响应
	streamChan, err := util.StreamChatCompletion(ctx, chatRequest)
	if err != nil {
		ws.WriteJSON(map[string]interface{}{
			"error": "流失请求失败" + err.Error(),
		})
		return
	}

	// 发送流式响应
	for response := range streamChan {
		if err := ws.WriteJSON(response); err != nil {
			return
		}
	}
}

// 初始化函数，启动WebSocket管理器
func init() {
	go Manager.Start()
}
