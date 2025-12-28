package tables

import (
	"context"
	"fmt"

	"github.com/finance-dashboard/backend/internal/pkg/models"
	"github.com/finance-dashboard/backend/internal/pkg/postgres"
)

const usersTableName = "users"

type Payments interface {
	GetByUserID(ctx context.Context, userID string) ([]*models.Payment, error)
}

type service struct {
	psql  *postgres.Service
	table string
}

func NewPayments(psql *postgres.Service) Payments {
	return &service{psql: psql, table: usersTableName}
}

func (s *service) GetByUserID(ctx context.Context, userID string) ([]*models.Payment, error) {
	statement, args, err := s.psql.Builder.
		Select(columns(new(models.Payment))...).
		From(s.table).
		Where("user_id = ?", userID).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetByUserID.s.psql.Builder err: %w", err)
	}

	rows, err := s.psql.DB.QueryContext(ctx, statement, args)
	if err != nil {
		return nil, fmt.Errorf("GetByUserID.s.psql.DB.QueryContext err: %w", err)
	}

	res, err := scanRows[models.Payment](rows)
	if err != nil {
		return nil, fmt.Errorf("GetByUserID.scanRows err: %w", err)
	}

	return res, nil
}
