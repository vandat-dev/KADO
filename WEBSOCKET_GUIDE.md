# WebSocket Guide - KADO

## 🚀 WebSocket Endpoint

**URL:** `/ws`

**Parameters:**
- `user_id` (required): ID của user kết nối

**Example:** 
```
ws://localhost:8386/ws?user_id=user123
```

## 🏗️ Kiến trúc WebSocket Manager

### WebSocket Manager Interface
```go
type WebSocketManager interface {
    Broadcast(message map[string]any)
    SendToUser(userID string, message map[string]any) int
    Connect(w http.ResponseWriter, r *http.Request, userID string) (*websocket.Conn, error)
    Disconnect(ws *websocket.Conn)
    GetOnlineUsers() []string
    CountAllConnections() int
    PushTaskToUsers(users []string, message map[string]any)
}
```

### Khởi tạo
WebSocket Manager được khởi tạo trong `internal/initialize/run.go`:
```go
func Run() {
    LoadConfig()
    InitLogger()
    Mysql()
    Redis()
    InitWebSocketManager() // ← Khởi tạo WebSocket Manager
    
    r := InitRouter()
    r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    err := r.Run(":8386")
    handleErr(err)
}
```

WebSocket Manager chỉ lo:
- ✅ **Connect/Disconnect** users
- ✅ **Quản lý connections** (mapping user_id → websocket connections)
- ✅ **Send/Broadcast** messages
- ✅ **Real-time notifications** cho các events trong hệ thống
- ❌ **KHÔNG xử lý logic tin nhắn** - để các service khác handle

## 📝 Cách sử dụng

### 1. Kết nối WebSocket
```javascript
const ws = new WebSocket('ws://localhost:8386/ws?user_id=user123');
```

### 2. Frontend gửi tin nhắn
```javascript
// Gửi tin nhắn trực tiếp cho user khác
ws.send(JSON.stringify({
    type: 'direct',
    to: 'user456',
    message: 'Hello user456!',
    timestamp: new Date().toISOString()
}));

// Broadcast tin nhắn cho tất cả users
ws.send(JSON.stringify({
    type: 'broadcast',
    message: 'Hello everyone!',
    timestamp: new Date().toISOString()
}));
```

### 3. Nhận tin nhắn từ server
```javascript
ws.onmessage = function(event) {
    const data = JSON.parse(event.data);
    
    // Xử lý các loại message khác nhau
    switch(data.type) {
        case 'new_product':
            console.log('New product created:', data.product_name);
            // Update UI để hiển thị sản phẩm mới
            break;
        case 'notification':
            console.log('Notification:', data.message);
            // Hiển thị thông báo
            break;
        default:
            console.log('Received:', data);
    }
};
```

## 🧪 Test WebSocket

1. **Khởi động server Go:**
   ```bash
   go run cmd/server/main.go
   ```

2. **Test WebSocket connection:**
   - Mở file `websocket-test.html` trong browser
   - Nhập User ID và click "Connect"

3. **Test Product Broadcast:**
   - Kết nối WebSocket với user_id
   - Tạo product mới qua API: `POST /v1/products`
   - Kiểm tra xem có nhận được broadcast message không

4. **Test với nhiều users:** 
   - Mở nhiều tab browser với User ID khác nhau
   - Tạo product và xem tất cả users có nhận được broadcast không

## 🎯 Real-time Features

### Product Creation Broadcast
Khi tạo product mới, hệ thống sẽ tự động broadcast tin nhắn:
```json
{
    "type": "new_product",
    "message": "New product: Product Name",
    "product_id": 123,
    "product_name": "Product Name",
    "time": 1703123456
}
```

## 🔧 Sử dụng trong Go Services

WebSocket Manager được truy cập thông qua `global.WsManager`:

### Gửi tin nhắn cho một user cụ thể:
```go
import (
    "base_go_be/global"
    "time"
)

// Trong service của bạn
func (s *ChatService) SendNotification(userID string, notification Notification) {
    if global.WsManager != nil {
        message := map[string]any{
            "type": "notification",
            "data": notification,
            "timestamp": time.Now().Unix(),
        }
        
        // Gửi qua websocket
        sent := global.WsManager.SendToUser(userID, message)
        if sent == 0 {
            log.Printf("User %s not online, save to database", userID)
            // Lưu vào database để gửi sau
        }
    }
}
```

### Broadcast tin nhắn cho tất cả users:
```go
func (s *SystemService) BroadcastMaintenance(message string) {
    if global.WsManager != nil {
        broadcastMsg := map[string]any{
            "type": "system_announcement",
            "message": message,
            "timestamp": time.Now().Unix(),
        }
        global.WsManager.Broadcast(broadcastMsg)
    }
}
```

### Broadcast khi tạo Product (đã implement):
```go
func (ps *ProductService) CreateProduct(name, description string, userID uint) (uint, error) {
    // ... tạo product ...
    
    // Broadcast new product
    if global.WsManager != nil {
        global.WsManager.Broadcast(map[string]any{
            "type":         "new_product",
            "message":      "New product: " + name,
            "product_id":   createdProduct.ID,
            "product_name": name,
            "time":         time.Now().Unix(),
        })
    }
    
    return createdProduct.ID, nil
}
```

### Lấy thông tin users online:
```go
func (s *UserService) GetOnlineStatus() OnlineStatus {
    if global.WsManager == nil {
        return OnlineStatus{Users: []string{}, Total: 0}
    }
    
    onlineUsers := global.WsManager.GetOnlineUsers()
    totalConnections := global.WsManager.CountAllConnections()
    
    return OnlineStatus{
        Users: onlineUsers,
        Total: totalConnections,
    }
}
```

## 📋 WebSocket Manager Features

- ✅ **Multi-user connections** (một user có thể có nhiều connections)
- ✅ **Auto disconnect cleanup** khi client ngắt kết nối  
- ✅ **Send to specific user** - gửi tin nhắn riêng cho user cụ thể
- ✅ **Broadcast to all users** - gửi tin nhắn cho tất cả users online
- ✅ **Online users tracking** - theo dõi users đang online
- ✅ **Connection management** - quản lý mapping user_id → websocket connections
- ❌ **Message processing** - KHÔNG xử lý logic tin nhắn (để service khác làm)

## 🏗️ Kiến trúc hệ thống

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Frontend      │    │  global.WsManager│    │   Services      │
│   (React/Vue)   │◄──►│   Interface      │◄──►│ (Product/User)  │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌──────────────────┐
                       │ ConnectionManager│
                       │ Implementation   │
                       └──────────────────┘
```

**Flow khởi tạo:**
1. `Run()` → `InitWebSocketManager()`
2. `global.WsManager = NewConnectionManager()`
3. Services sử dụng `global.WsManager.Broadcast()`

**Responsibilities:**
- **global.WsManager:** Interface thống nhất cho WebSocket operations
- **ConnectionManager:** Implementation cụ thể cho connection management
- **Services:** Business logic + real-time notifications
- **Frontend:** UI và real-time updates

## 🚨 Lưu ý quan trọng

⚠️ **Luôn kiểm tra nil:** `if global.WsManager != nil` trước khi sử dụng
✅ **Thread-safe:** ConnectionManager đã handle concurrent access
🔄 **Auto-cleanup:** Connections được tự động dọn dẹp khi disconnect
📡 **Real-time ready:** Hệ thống đã sẵn sàng cho real-time features
