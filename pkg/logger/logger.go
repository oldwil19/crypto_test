package logger

import (
	"log"
	"os"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
)

// InitLogger inicializa los loggers para información y errores
func InitLogger() {
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Info registra mensajes de información
func Info(v ...interface{}) {
	infoLogger.Println(v...)
}

// Error registra mensajes de error
func Error(v ...interface{}) {
	errorLogger.Println(v...)
}

// Warn registra mensajes de advertencia.
func Warn(v ...interface{}) {
	log.Println("[WARN]:", v)
}
