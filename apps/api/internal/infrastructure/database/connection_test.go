package database

import (
	"context"
	"testing"
)

func TestNewConnection_InvalidURL(t *testing.T) {
	ctx := context.Background()
	_, err := NewConnection(ctx, "invalid://not-a-valid-url")
	if err == nil {
		t.Error("expected error for invalid database URL, got nil")
	}
}

func TestNewConnection_EmptyURL(t *testing.T) {
	ctx := context.Background()
	_, err := NewConnection(ctx, "")
	if err == nil {
		t.Error("expected error for empty database URL, got nil")
	}
}

func TestDB_Close_NilPool(t *testing.T) {
	db := &DB{Pool: nil}
	// Should not panic
	db.Close()
}

func TestHealthCheck_NilPool(t *testing.T) {
	db := &DB{Pool: nil}
	ctx := context.Background()
	err := db.HealthCheck(ctx)
	if err == nil {
		t.Error("expected error for nil pool health check, got nil")
	}
}
