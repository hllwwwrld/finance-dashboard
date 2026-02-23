package tables

import (
	"context"
	"fmt"

	"github.com/finance-dashboard/backend/internal/pkg/models"
	"github.com/finance-dashboard/backend/internal/pkg/postgres"
)

const usersTableName = "public.user"

type Users interface {
	Create(ctx context.Context, user models.User) (*models.User, error)
	GetByLogin(ctx context.Context, login string) (*models.User, error)
	UpdateMonthlyIncome(ctx context.Context, login string, monthlyIncome int) (*models.User, error)
}

type usersTable struct {
	psql  *postgres.Service
	table string
}

func NewUsers(psql *postgres.Service) Users {
	return &usersTable{psql: psql, table: usersTableName}
}

func (s *usersTable) Create(ctx context.Context, user models.User) (*models.User, error) {
	statement, args, err := s.psql.Builder.
		Insert(s.table).
		Columns("login", "password").
		Values(user.Login, user.Password).
		Suffix(returningAllColumns(user)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetByUserID.s.psql.Builder err: %w", err)
	}

	row := s.psql.DB.QueryRowContext(ctx, statement, args...)
	if row.Err() != nil {
		return nil, fmt.Errorf("Create.User.s.psql.DB.QueryRowContext err: %v", row.Err())
	}

	res := &models.User{}
	err = row.Scan(
		&res.ID,
		&res.Login,
		&res.Password,
		&res.MonthlyIncome,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("Create.s.psql.DB.QueryRowContext err: %w", err)
	}

	return res, nil
}

func (s *usersTable) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	statement, args, err := s.psql.Builder.
		Select(allColumnsString(new(models.User))).
		From(s.table).
		Where("login = ?", login).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetByLogin.s.psql.Builder err: %w", err)
	}

	res := &models.User{}
	err = s.psql.DB.QueryRowContext(ctx, statement, args...).Scan(
		&res.ID,
		&res.Login,
		&res.Password,
		&res.MonthlyIncome,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("GetByLogin.s.psql.DB.QueryRowContext err: %w", err)
	}

	return res, nil
}

func (s *usersTable) UpdateMonthlyIncome(ctx context.Context, login string, monthlyIncome int) (*models.User, error) {
	statement, args, err := s.psql.Builder.
		Update(s.table).
		Set("monthly_income", monthlyIncome).
		Where("login = ?", login).
		Suffix(returningAllColumns(new(models.User))).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetByLogin.s.psql.Builder err: %w", err)
	}

	res := &models.User{}
	err = s.psql.DB.QueryRowContext(ctx, statement, args...).Scan(
		&res.ID,
		&res.Login,
		&res.Password,
		&res.MonthlyIncome,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("GetByLogin.s.psql.DB.QueryRowContext err: %w", err)
	}

	return res, nil
}
