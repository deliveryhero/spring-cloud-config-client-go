package logging

import "context"

type Logger interface {
	ErrorContext(ctx context.Context, msg string, keysAndValues ...interface{})
}
