package logger

import (
	"log"
	"os"
	"sync"
)

type GlobalLogger struct {
	filename string
	*log.Logger
}

var logger *GlobalLogger
var once sync.Once

func GetInstance() *GlobalLogger {
	once.Do(func() {
		logger = createLogger("log.txt")
	})
	return logger
}

func createLogger(fileName string) *GlobalLogger {
	file, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)

	return &GlobalLogger{
		filename: fileName,
		Logger:   log.New(file, "sunset-wallpaper-changer-go ", log.Lshortfile),
	}
}
