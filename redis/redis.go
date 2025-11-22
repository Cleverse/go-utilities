package redis

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Cleverse/go-utilities/logger"
	"github.com/Cleverse/go-utilities/logger/slogx"
	"github.com/cockroachdb/errors"
	redis "github.com/redis/go-redis/v9"
)

// Config defines redis config
type Config struct {
	Address    string `env:"ADDRESS" mapstructure:"address" json:"address"`
	Password   string `env:"PASSWORD" mapstructure:"password" json:"password"`
	Database   int    `env:"DATABASE" mapstructure:"database" json:"database"`        // optional
	AuthString string `env:"AUTHSTRING" mapstructure:"auth_string" json:"authString"` // optional
	PoolSize   int    `env:"POOL_SIZE" mapstructure:"pool_size" json:"poolSize"`      // optional
}

// RedisLogger implements github.com/redis/go-redis/v9/internal.Logging interface
type RedisLogger struct {
	logger *slog.Logger
}

func (l *RedisLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	l.logger.WarnContext(ctx, fmt.Sprintf(format, v...))
}

func SetLogger(l *slog.Logger) {
	redis.SetLogger(&RedisLogger{
		logger: l,
	})
}

func init() {
	SetLogger(logger.With(slogx.String("module", "redis")))
}

// New creates new redis instance
func New(ctx context.Context, config Config) (*redis.Client, error) {
	if config.Address == "" {
		config.Address = "localhost:6379"
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:        config.Address,
		Password:    config.Password,
		DB:          config.Database,
		PoolSize:    config.PoolSize,
		ReadTimeout: -1,
	})
	if config.AuthString != "" {
		if _, err := redisClient.Do(ctx, "AUTH", config.AuthString).Result(); err != nil {
			redisClient.Close()
			return nil, errors.WithStack(err)
		}
	}
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to redis")
	}
	return redisClient, nil
}
