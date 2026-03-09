package tables

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/finance-dashboard/backend/internal/pkg/models"
	"github.com/finance-dashboard/backend/internal/pkg/postgres"
)

const paymentsTableName = "public.payment"

type Payments interface {
	GetByUserID(ctx context.Context, userID string) ([]*models.Payment, error)
	Create(ctx context.Context, payment models.Payment) (*models.Payment, error)
	DeleteByID(ctx context.Context, id string) error
}

type paymentsTable struct {
	psql  *postgres.Service
	table string
}

func NewPayments(psql *postgres.Service) Payments {
	return &paymentsTable{psql: psql, table: paymentsTableName}
}

func (s *paymentsTable) Create(ctx context.Context, payment models.Payment) (*models.Payment, error) {
	statement, args, err := s.psql.Builder.
		Insert(s.table).
		Columns("user_id", "name", "amount", "due_day", "category", "color").
		Values(
			payment.UserID,
			payment.Name,
			payment.Amount,
			payment.DueDay,
			payment.Category,
			payment.Color).
		Suffix(returningAllColumns(payment)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("CreatePayment.s.psql.Builder err: %w", err)
	}

	res := &models.Payment{}
	err = s.psql.DB.QueryRowContext(ctx, statement, args...).Scan(
		&res.ID,
		&res.UserID,
		&res.Name,
		&res.Amount,
		&res.DueDay,
		&res.Category,
		&res.Color,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("CreatePayment.s.psql.DB.QueryRowContext err: %w", err)
	}

	return res, nil
}

func (s *paymentsTable) GetByUserID(ctx context.Context, userID string) ([]*models.Payment, error) {
	statement, args, err := s.psql.Builder.
		Select(columns(new(models.Payment))...).
		From(s.table).
		Where("user_id = ?", userID).
		Where("deleted_at is null").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetByUserID.s.psql.Builder err: %w", err)
	}

	rows, err := s.psql.DB.QueryContext(ctx, statement, args...)
	if err != nil {
		return nil, fmt.Errorf("GetByUserID.s.psql.DB.QueryContext err: %w", err)
	}

	res, err := scanPaymentRows(rows)
	if err != nil {
		return nil, fmt.Errorf("GetByUserID.scanRows err: %w", err)
	}

	return res, nil
}

func (s *paymentsTable) DeleteByID(ctx context.Context, id string) error {
	statement, args, err := s.psql.Builder.
		Update(s.table).
		Where("id = ?", id).
		Set("deleted_at", time.Now()).
		ToSql()
	if err != nil {
		return fmt.Errorf("DeleteByID.s.psql.Builder err: %w", err)
	}

	_, err = s.psql.DB.QueryContext(ctx, statement, args...)
	if err != nil {
		return fmt.Errorf("CreatePayment.s.psql.DB.QueryRowContext err: %w", err)
	}

	return nil
}

func scanPaymentRows(rows *sql.Rows) (res []*models.Payment, err error) {
	res = make([]*models.Payment, 0)
	for rows.Next() {
		var current models.Payment
		err = rows.Scan(
			&current.ID,
			&current.UserID,
			&current.Name,
			&current.Amount,
			&current.DueDay,
			&current.Category,
			&current.Color,
			&current.CreatedAt,
			&current.UpdatedAt,
			&current.DeletedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("GetByUserID.rows.Scan err: %w", err)
		}
		res = append(res, &current)
	}

	return res, nil
}
