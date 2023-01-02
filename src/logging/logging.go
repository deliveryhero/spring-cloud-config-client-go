package logging

import "context"

type Logger interface {
	NewLoggerWithLevel(string, string) Logger
	ErrorContext(ctx context.Context, msg string, keysAndValues ...interface{})
}
