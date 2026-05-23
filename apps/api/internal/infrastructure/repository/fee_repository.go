package repository

import (
	"context"
	"database/sql"
	"fmt"

	"harmoni-api/internal/domain/entity"
	"harmoni-api/internal/domain/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Compile-time interface assertion.
var _ repository.FeeRepository = (*PostgresFeeRepository)(nil)

// PostgresFeeRepository implements repository.FeeRepository using PostgreSQL via pgx.
// Separate SQL methods are used for mandatory_fees and voluntary_fees tables.
type PostgresFeeRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresFeeRepository creates a new fee repository backed by PostgreSQL.
func NewPostgresFeeRepository(pool *pgxpool.Pool) *PostgresFeeRepository {
	return &PostgresFeeRepository{pool: pool}
}

// ──────────────────────────────────────────────────────────────────
// Mandatory Fee Methods
// ──────────────────────────────────────────────────────────────────

// CreateMandatory inserts a new mandatory fee record.
func (r *PostgresFeeRepository) CreateMandatory(ctx context.Context, fee *entity.MandatoryFee) (*entity.MandatoryFee, error) {
	query := `
		INSERT INTO mandatory_fees (tenant_id, amount, description, effective_date, paid_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		fee.TenantID,
		fee.Amount,
		fee.Description,
		fee.EffectiveDate,
		fee.PaidAt,
	).Scan(&fee.ID, &fee.CreatedAt, &fee.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create mandatory fee: %w", err)
	}

	return fee, nil
}

// CreateMandatoryTx creates a mandatory fee within an existing transaction.
func (r *PostgresFeeRepository) CreateMandatoryTx(ctx context.Context, tx pgx.Tx, fee *entity.MandatoryFee) (*entity.MandatoryFee, error) {
	query := `
		INSERT INTO mandatory_fees (tenant_id, amount, description, effective_date, paid_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := tx.QueryRow(ctx, query,
		fee.TenantID,
		fee.Amount,
		fee.Description,
		fee.EffectiveDate,
		fee.PaidAt,
	).Scan(&fee.ID, &fee.CreatedAt, &fee.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create mandatory fee in transaction: %w", err)
	}

	return fee, nil
}

// ListMandatoryByTenant returns all mandatory fees for a given tenant.
func (r *PostgresFeeRepository) ListMandatoryByTenant(ctx context.Context, tenantID string) ([]*entity.MandatoryFee, error) {
	query := `
		SELECT id, tenant_id, amount, description, effective_date, paid_at, created_at, updated_at
		FROM mandatory_fees
		WHERE tenant_id = $1
		ORDER BY effective_date DESC
	`

	rows, err := r.pool.Query(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list mandatory fees by tenant: %w", err)
	}
	defer rows.Close()

	var fees []*entity.MandatoryFee
	for rows.Next() {
		fee, err := scanMandatoryFee(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan mandatory fee row: %w", err)
		}
		fees = append(fees, fee)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating mandatory fee rows: %w", err)
	}

	return fees, nil
}

// UpdateMandatory updates an existing mandatory fee record.
// Returns sql.ErrNoRows when no fee is found.
func (r *PostgresFeeRepository) UpdateMandatory(ctx context.Context, fee *entity.MandatoryFee) (*entity.MandatoryFee, error) {
	query := `
		UPDATE mandatory_fees
		SET amount = $2, description = $3, effective_date = $4, paid_at = $5, updated_at = NOW()
		WHERE id = $1
		RETURNING id, tenant_id, amount, description, effective_date, paid_at, created_at, updated_at
	`

	updated, err := scanMandatoryFee(r.pool.QueryRow(ctx, query,
		fee.ID,
		fee.Amount,
		fee.Description,
		fee.EffectiveDate,
		fee.PaidAt,
	))
	if err != nil {
		return nil, fmt.Errorf("failed to update mandatory fee: %w", err)
	}

	return updated, nil
}

// DeleteMandatory removes a mandatory fee record by ID.
// Returns sql.ErrNoRows when no fee is found.
func (r *PostgresFeeRepository) DeleteMandatory(ctx context.Context, id string) error {
	query := `DELETE FROM mandatory_fees WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete mandatory fee: %w", err)
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// ──────────────────────────────────────────────────────────────────
// Voluntary Fee Methods
// ──────────────────────────────────────────────────────────────────

// CreateVoluntary inserts a new voluntary contribution record.
func (r *PostgresFeeRepository) CreateVoluntary(ctx context.Context, fee *entity.VoluntaryFee) (*entity.VoluntaryFee, error) {
	query := `
		INSERT INTO voluntary_fees (tenant_id, amount, description, effective_date, paid_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		fee.TenantID,
		fee.Amount,
		fee.Description,
		fee.EffectiveDate,
		fee.PaidAt,
	).Scan(&fee.ID, &fee.CreatedAt, &fee.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create voluntary fee: %w", err)
	}

	return fee, nil
}

// ListVoluntaryByTenant returns all voluntary contributions for a given tenant.
func (r *PostgresFeeRepository) ListVoluntaryByTenant(ctx context.Context, tenantID string) ([]*entity.VoluntaryFee, error) {
	query := `
		SELECT id, tenant_id, amount, description, effective_date, paid_at, created_at, updated_at
		FROM voluntary_fees
		WHERE tenant_id = $1
		ORDER BY effective_date DESC
	`

	rows, err := r.pool.Query(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list voluntary fees by tenant: %w", err)
	}
	defer rows.Close()

	var fees []*entity.VoluntaryFee
	for rows.Next() {
		fee, err := scanVoluntaryFee(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan voluntary fee row: %w", err)
		}
		fees = append(fees, fee)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating voluntary fee rows: %w", err)
	}

	return fees, nil
}

// UpdateVoluntary updates an existing voluntary contribution record.
// Returns sql.ErrNoRows when no fee is found.
func (r *PostgresFeeRepository) UpdateVoluntary(ctx context.Context, fee *entity.VoluntaryFee) (*entity.VoluntaryFee, error) {
	query := `
		UPDATE voluntary_fees
		SET amount = $2, description = $3, effective_date = $4, paid_at = $5, updated_at = NOW()
		WHERE id = $1
		RETURNING id, tenant_id, amount, description, effective_date, paid_at, created_at, updated_at
	`

	updated, err := scanVoluntaryFee(r.pool.QueryRow(ctx, query,
		fee.ID,
		fee.Amount,
		fee.Description,
		fee.EffectiveDate,
		fee.PaidAt,
	))
	if err != nil {
		return nil, fmt.Errorf("failed to update voluntary fee: %w", err)
	}

	return updated, nil
}

// DeleteVoluntary removes a voluntary contribution record by ID.
// Returns sql.ErrNoRows when no fee is found.
func (r *PostgresFeeRepository) DeleteVoluntary(ctx context.Context, id string) error {
	query := `DELETE FROM voluntary_fees WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete voluntary fee: %w", err)
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// ──────────────────────────────────────────────────────────────────
// Scan Helpers
// ──────────────────────────────────────────────────────────────────

// scanMandatoryFee scans a mandatory fee row (single or multi-row result).
func scanMandatoryFee(row interface {
	Scan(dest ...interface{}) error
}) (*entity.MandatoryFee, error) {
	fee := &entity.MandatoryFee{}
	var amount pgtype.Numeric

	err := row.Scan(
		&fee.ID,
		&fee.TenantID,
		&amount,
		&fee.Description,
		&fee.EffectiveDate,
		&fee.PaidAt,
		&fee.CreatedAt,
		&fee.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	dec, err := numericToDecimal(amount)
	if err != nil {
		return nil, fmt.Errorf("failed to decode mandatory fee amount: %w", err)
	}
	fee.Amount = dec

	return fee, nil
}

// scanVoluntaryFee scans a voluntary fee row (single or multi-row result).
func scanVoluntaryFee(row interface {
	Scan(dest ...interface{}) error
}) (*entity.VoluntaryFee, error) {
	fee := &entity.VoluntaryFee{}
	var amount pgtype.Numeric

	err := row.Scan(
		&fee.ID,
		&fee.TenantID,
		&amount,
		&fee.Description,
		&fee.EffectiveDate,
		&fee.PaidAt,
		&fee.CreatedAt,
		&fee.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	dec, err := numericToDecimal(amount)
	if err != nil {
		return nil, fmt.Errorf("failed to decode voluntary fee amount: %w", err)
	}
	fee.Amount = dec

	return fee, nil
}
