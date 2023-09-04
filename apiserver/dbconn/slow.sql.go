package dbconn

import (
	"context"
	"fmt"
)

type dbLogger struct {
	name      string
	alarmCode string
	queueName string
}

func (l dbLogger) Name() string {
	return l.name
}
func (l dbLogger) Log(ctx context.Context, elapsed int64, sql string, args ...interface{}) {
	fmt.Println(elapsed, sql, args)
}
