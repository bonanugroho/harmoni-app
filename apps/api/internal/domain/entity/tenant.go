package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

// Tenant represents a tenant unit in the Harmoni system.
type Tenant struct {
	ID          string          `json:"id"`
	Block       string          `json:"block"`
	UnitNumber  string          `json:"unit_number"`
	Occupancy   string          `json:"occupancy"`
	MonthlyFee  decimal.Decimal `json:"monthly_fee"`
	TerritoryID string          `json:"territory_id"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// Sanitize returns a copy of the tenant with all fields.
// Use this when returning tenant data to clients.
func (t *Tenant) Sanitize() *Tenant {
	return &Tenant{
		ID:          t.ID,
		Block:       t.Block,
		UnitNumber:  t.UnitNumber,
		Occupancy:   t.Occupancy,
		MonthlyFee:  t.MonthlyFee,
		TerritoryID: t.TerritoryID,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// MandatoryFee represents a mandatory fee assigned to a tenant unit.
type MandatoryFee struct {
	ID            string          `json:"id"`
	TenantID      string          `json:"tenant_id"`
	Amount        decimal.Decimal `json:"amount"`
	Description   string          `json:"description"`
	EffectiveDate time.Time       `json:"effective_date"`
	PaidAt        *time.Time      `json:"paid_at"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// Sanitize returns a copy of the mandatory fee with all fields.
// Use this when returning fee data to clients.
func (f *MandatoryFee) Sanitize() *MandatoryFee {
	return &MandatoryFee{
		ID:            f.ID,
		TenantID:      f.TenantID,
		Amount:        f.Amount,
		Description:   f.Description,
		EffectiveDate: f.EffectiveDate,
		PaidAt:        f.PaidAt,
		CreatedAt:     f.CreatedAt,
		UpdatedAt:     f.UpdatedAt,
	}
}

// VoluntaryFee represents a voluntary contribution from a tenant unit.
type VoluntaryFee struct {
	ID            string          `json:"id"`
	TenantID      string          `json:"tenant_id"`
	Amount        decimal.Decimal `json:"amount"`
	Description   string          `json:"description"`
	EffectiveDate time.Time       `json:"effective_date"`
	PaidAt        *time.Time      `json:"paid_at"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// Sanitize returns a copy of the voluntary fee with all fields.
// Use this when returning fee data to clients.
func (f *VoluntaryFee) Sanitize() *VoluntaryFee {
	return &VoluntaryFee{
		ID:            f.ID,
		TenantID:      f.TenantID,
		Amount:        f.Amount,
		Description:   f.Description,
		EffectiveDate: f.EffectiveDate,
		PaidAt:        f.PaidAt,
		CreatedAt:     f.CreatedAt,
		UpdatedAt:     f.UpdatedAt,
	}
}
