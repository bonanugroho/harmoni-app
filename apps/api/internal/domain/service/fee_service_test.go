package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestFeeService_CreateVoluntary_ValidDates(t *testing.T) {
	svc, _, _ := newTestTenantService(t)
	ctx := context.Background()

	futureDate := time.Now().Add(48 * time.Hour).Truncate(24 * time.Hour)
	paidAt := futureDate.Add(24 * time.Hour)

	// First create a tenant
	req := validCreateRequest()
	claims := validRTClaims()

	created, err := svc.Create(ctx, req, claims)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	volReq := &CreateVoluntaryFeeRequest{
		Amount:        decimal.NewFromInt(10000),
		Description:   "Donation",
		EffectiveDate: futureDate,
		PaidAt:        &paidAt,
	}

	fee, err := svc.CreateVoluntaryFee(ctx, created.ID, volReq, claims)
	if err != nil {
		t.Fatalf("CreateVoluntaryFee() error = %v", err)
	}

	if fee.Amount.Cmp(decimal.NewFromInt(10000)) != 0 {
		t.Errorf("Amount = %v, want 10000", fee.Amount)
	}
	if fee.Description != "Donation" {
		t.Errorf("Description = %q, want %q", fee.Description, "Donation")
	}
	if fee.ID == "" {
		t.Error("CreateVoluntaryFee() should return fee with ID")
	}
}

func TestFeeService_CreateVoluntary_PastEffectiveDate(t *testing.T) {
	svc, _, _ := newTestTenantService(t)
	ctx := context.Background()

	pastDate := time.Now().Add(-48 * time.Hour).Truncate(24 * time.Hour)

	// First create a tenant
	req := validCreateRequest()
	claims := validRTClaims()

	created, err := svc.Create(ctx, req, claims)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	volReq := &CreateVoluntaryFeeRequest{
		Amount:        decimal.NewFromInt(10000),
		Description:   "Past Donation",
		EffectiveDate: pastDate,
	}

	_, err = svc.CreateVoluntaryFee(ctx, created.ID, volReq, claims)
	if !errors.Is(err, ErrInvalidEffectiveDate) {
		t.Errorf("CreateVoluntaryFee() error = %v, want ErrInvalidEffectiveDate", err)
	}
}

func TestFeeService_CreateVoluntary_PaidAtBeforeEffective(t *testing.T) {
	svc, _, _ := newTestTenantService(t)
	ctx := context.Background()

	futureDate := time.Now().Add(48 * time.Hour).Truncate(24 * time.Hour)
	paidAt := futureDate.Add(-24 * time.Hour) // paid_at before effective_date

	// First create a tenant
	req := validCreateRequest()
	claims := validRTClaims()

	created, err := svc.Create(ctx, req, claims)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	volReq := &CreateVoluntaryFeeRequest{
		Amount:        decimal.NewFromInt(10000),
		Description:   "Bad Timing Donation",
		EffectiveDate: futureDate,
		PaidAt:        &paidAt,
	}

	_, err = svc.CreateVoluntaryFee(ctx, created.ID, volReq, claims)
	if !errors.Is(err, ErrInvalidPaidAt) {
		t.Errorf("CreateVoluntaryFee() error = %v, want ErrInvalidPaidAt", err)
	}
}
