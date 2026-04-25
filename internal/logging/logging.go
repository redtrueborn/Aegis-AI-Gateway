package logging


import (
	"log/slog"
)
func InitLogger() *slog.Logger {
   return slog.Default()	
}
