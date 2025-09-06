package initialize

import (
	"base_go_be/global"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ConnectionManager struct {
	mu sync.RWMutex
	// map user_id
	connections map[string]map[*websocket.Conn]struct{}
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]map[*websocket.Conn]struct{}),
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	//Adjust CORS depending on the environment
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (cm *ConnectionManager) Connect(w http.ResponseWriter, r *http.Request, userID string) (*websocket.Conn, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	cm.mu.Lock()
	if _, ok := cm.connections[userID]; !ok {
		cm.connections[userID] = make(map[*websocket.Conn]struct{})
	}
	cm.connections[userID][ws] = struct{}{}
	cm.mu.Unlock()

	log.Printf("User %s connected.", userID)
	return ws, nil
}

func (cm *ConnectionManager) Disconnect(ws *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for userID, set := range cm.connections {
		if _, ok := set[ws]; ok {
			delete(set, ws)
			_ = ws.Close()
			if len(set) == 0 {
				delete(cm.connections, userID)
			}
			log.Printf("User %s disconnected. Total connections: %d", userID, cm.CountAllConnectionsLocked())
			break
		}
	}
}

func (cm *ConnectionManager) SendToUser(userID string, message map[string]any) int {
	cm.mu.RLock()
	set, ok := cm.connections[userID]
	if !ok {
		cm.mu.RUnlock()
		return 0
	}

	// copy snapshot
	sockets := make([]*websocket.Conn, 0, len(set))
	for ws := range set {
		sockets = append(sockets, ws)
	}
	cm.mu.RUnlock()

	payload, err := json.Marshal(message)
	if err != nil {
		log.Printf("json marshal error for user %s: %v", userID, err)
		return 0
	}

	log.Printf("Send message to user %s: %s", userID, string(payload))
	var toDisconnect []*websocket.Conn
	success := 0

	for _, ws := range sockets {
		_ = ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err := ws.WriteMessage(websocket.TextMessage, payload); err != nil {
			log.Printf("Send error to %s: %v", userID, err)
			toDisconnect = append(toDisconnect, ws)
		} else {
			success++
		}
	}

	for _, ws := range toDisconnect {
		cm.Disconnect(ws)
	}
	return success
}

func (cm *ConnectionManager) Broadcast(message map[string]any) {
	log.Println("Start broadcast WebSocket")

	// snapshot all connect
	cm.mu.RLock()
	var sockets []*websocket.Conn
	for _, set := range cm.connections {
		for ws := range set {
			sockets = append(sockets, ws)
		}
	}
	cm.mu.RUnlock()

	payload, _ := json.Marshal(message)
	for _, ws := range sockets {
		_ = ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err := ws.WriteMessage(websocket.TextMessage, payload); err != nil {
			log.Printf("Broadcast error: %v", err)
			cm.Disconnect(ws)
		}
	}
}

func (cm *ConnectionManager) PushTaskToUsers(users []string, message map[string]any) {
	for _, uid := range users {
		cm.SendToUser(uid, message)
	}
}

func (cm *ConnectionManager) GetOnlineUsers() []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	users := make([]string, 0, len(cm.connections))
	for uid := range cm.connections {
		users = append(users, uid)
	}
	return users
}

func (cm *ConnectionManager) CountAllConnections() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.CountAllConnectionsLocked()
}

func (cm *ConnectionManager) CountAllConnectionsLocked() int {
	total := 0
	for _, set := range cm.connections {
		total += len(set)
	}
	return total
}

// InitWebSocketManager initializes the global WebSocket manager
func InitWebSocketManager() {
	global.WsManager = NewConnectionManager()
}

func WebSocketHandler(c *gin.Context) {

	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user_id parameter"})
		return
	}

	ws, err := global.WsManager.Connect(c.Writer, c.Request, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "websocket upgrade failed"})
		return
	}

	// Goroutine reads to detect client close and cleans up
	go func(userID string, ws *websocket.Conn) {
		defer global.WsManager.Disconnect(ws)
		for {
			_, msgData, err := ws.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("websocket error: %v", err)
				}
				break
			}

			var payload map[string]any
			if err := json.Unmarshal(msgData, &payload); err != nil {
				log.Printf("invalid message from %s: %v", userID, err)
				continue
			}

			// handle type
			switch payload["type"] {
			case "direct":
				to := payload["to"].(string)
				global.WsManager.SendToUser(to, payload)
			case "broadcast":
				global.WsManager.Broadcast(payload)
			default:
				log.Printf("unknown message type from %s: %+v", userID, payload)
			}
		}
	}(userID, ws)
}
