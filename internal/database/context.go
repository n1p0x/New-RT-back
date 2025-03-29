package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type contextKey = string

const dbKey = contextKey("db")

func WithDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, dbKey, db)
}

func FromContext(ctx context.Context, db *gorm.DB) *gorm.DB {
	if ctx == nil {
		return db
	}
	if stored, ok := ctx.Value(dbKey).(*gorm.DB); ok {
		return stored
	}
	return db
}

func RunInTx(ctx context.Context, db *gorm.DB, f func(ctx context.Context) error) error {
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start tx: %v", tx.Error)
	}

	ctx = WithDB(ctx, tx)
	if err := f(ctx); err != nil {
		if err1 := tx.Rollback().Error; err1 != nil {
			return fmt.Errorf("rollback tx: %v", err1)
		}
		return fmt.Errorf("invoke function: %v", err)
	}
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit tx: %v", err)
	}
	return nil
}
