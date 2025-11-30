package Tracelog

import (
	"fmt"
	"log"
	"os"
	"time"
)

var logFile *os.File

func init() {
	// Buat folder log kalau belum ada
	_ = os.MkdirAll("logs", os.ModePerm)

	var err error
	logFile, err = os.OpenFile("logs/tracelog_auth.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Gagal membuka file tracelog_auth.log: %v", err)
	}

	log.SetOutput(logFile)
}

// writeLog hanya memiliki 2 parameter:
// message → pesan log
// endpoint → nama endpoint / handler
func WriteLog(message, endpoint string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	logLine := fmt.Sprintf("%s --- %s --- Endpoint : %s", timestamp, message, endpoint)

	log.Println(logLine)
}
