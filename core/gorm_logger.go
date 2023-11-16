package core

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"gorm.io/gorm/logger"
)

var (
	_ logger.Interface = GormLogger{}

	SlowThresholdMilliseconds int64 = 500
)

type GormLogger struct {
	log log.Logger
}

func NewGormLogger(log log.Logger) GormLogger {
	return GormLogger{log}
}

func (l GormLogger) LogMode(lvl logger.LogLevel) logger.Interface {
	return l
}

func (l GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.log.Info(fmt.Sprintf(msg, data...))
}

func (l GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.log.Warn(fmt.Sprintf(msg, data...))
}

func (l GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.log.Error(fmt.Sprintf(msg, data...))
}

func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsedMs := time.Since(begin).Milliseconds()

	// omit any values for batch inserts as they can be very long
	sql, rows := fc()
	if i := strings.Index(strings.ToLower(sql), "values"); i > 0 {
		sql = fmt.Sprintf("%sVALUES (...)", sql[:i])
	}

	if elapsedMs < SlowThresholdMilliseconds {
		l.log.Debug("database operation", "duration_ms", elapsedMs, "rows_affected", rows, "sql", sql)
	} else {
		l.log.Warn("database operation", "duration_ms", elapsedMs, "rows_affected", rows, "sql", sql)
	}
}
