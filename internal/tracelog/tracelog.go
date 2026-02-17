package tracelog

import (
	"fmt"
	"log"
	"os"
	"time"
)

func write(module string, message string, endpoint string) {
	_ = os.MkdirAll("logs", os.ModePerm)

	date := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("logs/%s_%s.txt", module, date)

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Gagal membuka file log: %v", err)
		return
	}
	defer file.Close()

	logger := log.New(file, "", 0)

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	line := fmt.Sprintf("%s --- %s --- Endpoint: %s", timestamp, message, endpoint)

	logger.Println(line)
}

// Untuk Auth
func AuthLog(message, endpoint string) {
	write("tracelog_auth", message, endpoint)
}

// Untuk User Management
func UserManagementLog(message, endpoint string) {
	write("tracelog_user_management", message, endpoint)
}

// Untuk Dokumen Management
func DokumenManagementLog(message, endpoint string) {
	write("tracelog_dokumen_management", message, endpoint)
}

// Untuk General Management
func HomeLog(message, endpoint string) {
	write("tracelog_home", message, endpoint)
}

// Untuk Topic
func TopicLog(message, endpoint string) {
	write("tracelog_topic", message, endpoint)
}

// Untuk Category
func CategoryLog(message, endpoint string) {
	write("tracelog_category", message, endpoint)
}

// Untuk General Management
func WebSocketLog(message, endpoint string) {
	write("tracelog_websocket", message, endpoint)
}
