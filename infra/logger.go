package infra

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log/slog"
	"time"
)

func ResolverLogger(context context.Context) *slog.Logger {
	val := context.Value("logger")

	if val == nil {
		return slog.Default()
	}

	log, ok := val.(*slog.Logger)
	if !ok {
		return slog.Default()
	}

	return log
}

type SlogGorm struct {
	Logger                    *slog.Logger
	LogLevel                  logger.LogLevel
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
}

func (s *SlogGorm) LogMode(level logger.LogLevel) logger.Interface {
	return &SlogGorm{
		Logger:        s.Logger,
		SlowThreshold: s.SlowThreshold,
		LogLevel:      level,
	}
}

func (s *SlogGorm) Info(ctx context.Context, query string, args ...interface{}) {
	if s.LogLevel < logger.Info {
		return
	}

	s.Logger.InfoContext(ctx, query, args)
}

func (s *SlogGorm) Warn(ctx context.Context, query string, args ...interface{}) {
	if s.LogLevel < logger.Warn {
		return
	}

	s.Logger.WarnContext(ctx, query, args)
}

func (s *SlogGorm) Error(ctx context.Context, query string, args ...interface{}) {
	if s.LogLevel < logger.Error {
		return
	}

	s.Logger.ErrorContext(ctx, query, args)
}

func (s *SlogGorm) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if s.LogLevel <= 0 {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		if s.LogLevel >= logger.Error && (!s.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)) {
			s.Logger.ErrorContext(ctx, "error during query executing",
				slog.Duration("elapsed", elapsed),
				slog.String("query", sql),
				slog.Int64("rows_affected", rows),
				slog.Any("err", err),
			)
		}
		return
	}

	if s.SlowThreshold > 0 && elapsed > s.SlowThreshold && s.LogLevel >= logger.Warn {
		s.Logger.WarnContext(ctx, "slow query detected",
			slog.Duration("elapsed", elapsed),
			slog.Int64("rows", rows),
			slog.String("sql", sql),
		)
		return
	}

	if s.LogLevel >= logger.Info {
		s.Logger.WarnContext(ctx, "query executed",
			slog.Duration("elapsed", elapsed),
			slog.Int64("rows", rows),
			slog.String("sql", sql),
		)
	}
}
