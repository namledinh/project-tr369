package uspstore

import (
	"context"
	"time"
	apperrors "usp-management-device-api/common/app_errors"

	"gorm.io/gorm"
)

type contextKey string

const (
	defaultQueryTimeout            = 30 * time.Second
	maxQueryTimeout                = 5 * time.Minute
	txKey               contextKey = "tx_key"
)

func (s *store) GetDBName() string {
	return "usp_system_db"
}

type store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *store {
	return &store{
		db: db,
	}
}

// Helper to get DB from context if transaction exists
func (s *store) getDBFromContext(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(txKey).(*gorm.DB)
	if ok {
		return tx
	}
	return s.db
}

// BeginTx starts a transaction and stores it in the context
func (s *store) BeginTx(ctx context.Context) (context.Context, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		return ctx, apperrors.NewDBError(tx.Error, "SQL")
	}
	return context.WithValue(ctx, txKey, tx), nil
}

// CommitTx commits a transaction stored in the context
func (s *store) CommitTx(ctx context.Context) error {
	tx, ok := ctx.Value(txKey).(*gorm.DB)
	if !ok {
		return apperrors.NewInternalError(nil, "No transaction found in context")
	}

	if err := tx.Commit().Error; err != nil {
		return apperrors.NewDBError(err, "SQL")
	}
	return nil
}

// RollbackTx rolls back a transaction stored in the context
func (s *store) RollbackTx(ctx context.Context) error {
	tx, ok := ctx.Value(txKey).(*gorm.DB)
	if !ok {
		return apperrors.NewInternalError(nil, "No transaction found in context")
	}

	if err := tx.Rollback().Error; err != nil {
		return apperrors.NewDBError(err, "SQL")
	}
	return nil
}
