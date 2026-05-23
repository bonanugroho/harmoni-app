package repository

import (
	"context"

	"harmoni-api/internal/domain/entity"

	"github.com/jackc/pgx/v5"
)

// FeeRepository defines the interface for fee data access.
// Mandatory and voluntary fees are stored in separate tables with separate methods.
type FeeRepository interface {
	// --- Mandatory Fees ---

	// CreateMandatory inserts a new mandatory fee record.
	CreateMandatory(ctx context.Context, fee *entity.MandatoryFee) (*entity.MandatoryFee, error)

	// CreateMandatoryTx creates a mandatory fee within an existing transaction.
	// Used for atomic tenant+fee creation where the service layer manages the transaction.
	CreateMandatoryTx(ctx context.Context, tx pgx.Tx, fee *entity.MandatoryFee) (*entity.MandatoryFee, error)

	// ListMandatoryByTenant returns all mandatory fees for a given tenant.
	ListMandatoryByTenant(ctx context.Context, tenantID string) ([]*entity.MandatoryFee, error)

	// UpdateMandatory updates an existing mandatory fee record.
	// Returns sql.ErrNoRows when no fee is found.
	UpdateMandatory(ctx context.Context, fee *entity.MandatoryFee) (*entity.MandatoryFee, error)

	// DeleteMandatory removes a mandatory fee record by ID.
	// Returns sql.ErrNoRows when no fee is found.
	DeleteMandatory(ctx context.Context, id string) error

	// --- Voluntary Fees ---

	// CreateVoluntary inserts a new voluntary contribution record.
	CreateVoluntary(ctx context.Context, fee *entity.VoluntaryFee) (*entity.VoluntaryFee, error)

	// ListVoluntaryByTenant returns all voluntary contributions for a given tenant.
	ListVoluntaryByTenant(ctx context.Context, tenantID string) ([]*entity.VoluntaryFee, error)

	// UpdateVoluntary updates an existing voluntary contribution record.
	// Returns sql.ErrNoRows when no fee is found.
	UpdateVoluntary(ctx context.Context, fee *entity.VoluntaryFee) (*entity.VoluntaryFee, error)

	// DeleteVoluntary removes a voluntary contribution record by ID.
	// Returns sql.ErrNoRows when no fee is found.
	DeleteVoluntary(ctx context.Context, id string) error
}
