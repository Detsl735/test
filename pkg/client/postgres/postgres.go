package postgres

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client struct {
	DB *gorm.DB
}

func NewClient(ctx context.Context, dsn string) (*Client, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm open: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("db.DB(): %w", err)
	}

	if err := pingWithTimeout(ctx, sqlDB, 5*time.Second); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}

	return &Client{DB: db}, nil
}

func (c *Client) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func pingWithTimeout(ctx context.Context, sqlDB interface {
	PingContext(context.Context) error
}, timeout time.Duration) error {
	ctxPing, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return sqlDB.PingContext(ctxPing)
}
