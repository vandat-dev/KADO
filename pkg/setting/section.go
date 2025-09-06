package setting

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Config struct {
	Server   ServerSetting   `map_structure:"server"`
	Mysql    MySQLSetting    `map_structure:"mysql"`
	Postgres PostgresSetting `map_structure:"postgres"`
	Redis    RedisSetting    `map_structure:"redis"`
	Logger   LoggerSetting   `map_structure:"logger"`
}

type ServerSetting struct {
	Port int    `map_structure:"port"`
	Mode string `map_structure:"mode"`
}

type MySQLSetting struct {
	Host            string `map_structure:"host"`
	Port            int    `map_structure:"port"`
	UserName        string `map_structure:"username"`
	Password        string `map_structure:"password"`
	DBName          string `map_structure:"dbname"`
	MaxIdleConn     int    `map_structure:"maxIdleConn"`
	MaxOpenConn     int    `map_structure:"maxOpenConn"`
	ConnMaxLifeTime int    `map_structure:"connMaxLifeTime"`
}

type PostgresSetting struct {
	Host            string `map_structure:"host"`
	Port            int    `map_structure:"port"`
	UserName        string `map_structure:"username"`
	Password        string `map_structure:"password"`
	DBName          string `map_structure:"dbname"`
	SSLMode         string `map_structure:"sslmode"`
	MaxIdleConn     int    `map_structure:"maxIdleConn"`
	MaxOpenConn     int    `map_structure:"maxOpenConn"`
	ConnMaxLifeTime int    `map_structure:"connMaxLifeTime"`
}

type LoggerSetting struct {
	LogLevel    string `map_structure:"log_level"`
	FileLogName string `map_structure:"file_log_name"`
	MaxBackups  int    `map_structure:"max_backups"`
	MaxSize     int    `map_structure:"max_size"`
	MaxAge      int    `map_structure:"max_age"`
	Compress    bool   `map_structure:"compress"`
}

type RedisSetting struct {
	Host     string `map_structure:"host"`
	Port     int    `map_structure:"port"`
	Password string `map_structure:"password"`
	Database int    `map_structure:"database"`
}

type WebSocketManager interface {
	Broadcast(message map[string]any)
	SendToUser(userID string, message map[string]any) int
	Connect(w http.ResponseWriter, r *http.Request, userID string) (*websocket.Conn, error)
	Disconnect(ws *websocket.Conn)
	GetOnlineUsers() []string
	CountAllConnections() int
	PushTaskToUsers(users []string, message map[string]any)
}
