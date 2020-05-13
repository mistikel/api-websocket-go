package utils

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
)

// InfoContext prints info message to logs
func InfoContext(ctx context.Context, format string, v ...interface{}) {
	message := fmt.Sprintf("\033[32mINFO :\033[0m ReqID "+extractReqID(ctx)+" - "+format, v...)
	log.Print(message)
}

// ErrContext prints error message to logs
func ErrContext(ctx context.Context, format string, v ...interface{}) {
	message := []interface{}{fmt.Sprintf("\033[31mERROR:\033[0m ReqID "+extractReqID(ctx)+" - "+format, v...)}
	message = append(message, "\n", string(debug.Stack()))
	log.Print(message...)

}

func WarnContext(ctx context.Context, format string, v ...interface{}) {
	message := fmt.Sprintf("\033[33mWARN :\033[0m ReqID "+extractReqID(ctx)+" - "+format, v...)
	log.Print(message)
}

func extractReqID(ctx context.Context) string {
	requestIDKey := "x-request-id"
	str, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return ""
	}
	return str
}
