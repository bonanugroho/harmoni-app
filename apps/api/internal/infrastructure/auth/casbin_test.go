package auth

import (
	"sync"
	"testing"
)

func TestCasbinEnforcer_InitAndEnforce(t *testing.T) {
	// Reset singleton for test
	enforcerOnce = sync.Once{}
	enforcerInstance = nil
	enforcerErr = nil

	ce, err := InitEnforcer("../../../rbac_model.conf", "../../../policy.csv")
	if err != nil {
		t.Fatalf("InitEnforcer failed: %v", err)
	}
	if ce == nil {
		t.Fatal("expected non-nil CasbinEnforcer")
	}
}

func TestCasbinEnforcer_Enforce(t *testing.T) {
	// Reset singleton for test
	enforcerOnce = sync.Once{}
	enforcerInstance = nil
	enforcerErr = nil

	_, err := InitEnforcer("../../../rbac_model.conf", "../../../policy.csv")
	if err != nil {
		t.Fatalf("InitEnforcer failed: %v", err)
	}

	ce := GetEnforcer()
	if ce == nil {
		t.Fatal("expected non-nil enforcer after InitEnforcer")
	}

	// Test RT officer accessing own territory with literal placeholder
	// The placeholder is stored as-is in the file adapter
	// This verifies the enforcer loads correctly
	ok, err := ce.Enforce("rt_officer", "tenant", "read", "{{territory_id}}")
	if err != nil {
		t.Fatalf("Enforce error: %v", err)
	}
	if !ok {
		t.Error("rt_officer should match {{territory_id}} placeholder when domain is also {{territory_id}}")
	}
}

func TestCasbinEnforcer_EnforceWithTerritory(t *testing.T) {
	// Reset singleton for test
	enforcerOnce = sync.Once{}
	enforcerInstance = nil
	enforcerErr = nil

	// For testing, we need a fresh enforcer with substituted territory
	// Use the test helper from casbin_policy_test.go approach
	e := newEnforcerForTerritoryTest(t, "rt-01")

	// RT officer accessing own territory
	ok, err := e.Enforce("rt_officer", "tenant", "read", "rt-01")
	if err != nil {
		t.Fatalf("Enforce error: %v", err)
	}
	if !ok {
		t.Error("rt_officer should be allowed to read tenant in rt-01")
	}

	// RT officer accessing other territory
	ok, err = e.Enforce("rt_officer", "tenant", "read", "rt-02")
	if err != nil {
		t.Fatalf("Enforce error: %v", err)
	}
	if ok {
		t.Error("rt_officer should NOT be allowed to read tenant in rt-02")
	}

	// RW officer accessing any territory (using * domain)
	ok, err = e.Enforce("rw_officer", "tenant", "read", "*")
	if err != nil {
		t.Fatalf("Enforce error: %v", err)
	}
	if !ok {
		t.Error("rw_officer should be allowed to read tenant with * domain")
	}
}

func TestCasbinEnforcer_AddRemovePolicy(t *testing.T) {
	// Reset singleton for test
	enforcerOnce = sync.Once{}
	enforcerInstance = nil
	enforcerErr = nil

	ce, err := InitEnforcer("../../../rbac_model.conf", "../../../policy.csv")
	if err != nil {
		t.Fatalf("InitEnforcer failed: %v", err)
	}

	// Add a new policy
	err = ce.AddPolicy("custom_role", "tenant", "read", "rt-03")
	if err != nil {
		t.Fatalf("AddPolicy failed: %v", err)
	}

	// Verify the policy works
	ok, err := ce.Enforce("custom_role", "tenant", "read", "rt-03")
	if err != nil {
		t.Fatalf("Enforce error: %v", err)
	}
	if !ok {
		t.Error("custom_role should be allowed to read tenant in rt-03")
	}

	// Remove the policy
	err = ce.RemovePolicy("custom_role", "tenant", "read", "rt-03")
	if err != nil {
		t.Fatalf("RemovePolicy failed: %v", err)
	}

	// Verify the policy is gone
	ok, err = ce.Enforce("custom_role", "tenant", "read", "rt-03")
	if err != nil {
		t.Fatalf("Enforce error: %v", err)
	}
	if ok {
		t.Error("custom_role should NOT be allowed after policy removal")
	}
}

func TestCasbinEnforcer_AddRemoveRoleLink(t *testing.T) {
	// Reset singleton for test
	enforcerOnce = sync.Once{}
	enforcerInstance = nil
	enforcerErr = nil

	ce, err := InitEnforcer("../../../rbac_model.conf", "../../../policy.csv")
	if err != nil {
		t.Fatalf("InitEnforcer failed: %v", err)
	}

	// Add a custom role link
	err = ce.AddRoleLink("super_admin", "rw_officer")
	if err != nil {
		t.Fatalf("AddRoleLink failed: %v", err)
	}

	// Verify super_admin inherits rw_officer permissions (wildcard access)
	ok, err := ce.Enforce("super_admin", "tenant", "read", "*")
	if err != nil {
		t.Fatalf("Enforce error: %v", err)
	}
	if !ok {
		t.Error("super_admin should inherit rw_officer permissions")
	}

	// Remove the role link
	err = ce.RemoveRoleLink("super_admin", "rw_officer")
	if err != nil {
		t.Fatalf("RemoveRoleLink failed: %v", err)
	}

	// Verify the link is gone
	ok, err = ce.Enforce("super_admin", "tenant", "read", "*")
	if err != nil {
		t.Fatalf("Enforce error: %v", err)
	}
	if ok {
		t.Error("super_admin should NOT have rw_officer permissions after link removal")
	}
}

func TestCasbinEnforcer_GetEnforcer_NotInitialized(t *testing.T) {
	// Reset singleton
	enforcerOnce = sync.Once{}
	enforcerInstance = nil
	enforcerErr = nil

	ce := GetEnforcer()
	if ce != nil {
		t.Error("expected nil enforcer before initialization")
	}
}

func TestCasbinEnforcer_EnforceWithTerritory_RWOfficer(t *testing.T) {
	// Reset singleton for test
	enforcerOnce = sync.Once{}
	enforcerInstance = nil
	enforcerErr = nil

	ce, err := InitEnforcer("../../../rbac_model.conf", "../../../policy.csv")
	if err != nil {
		t.Fatalf("InitEnforcer failed: %v", err)
	}

	// RW officer with * domain should have access
	ok, err := ce.EnforceWithTerritory("rw_officer", "tenant", "read", "rt-01", true)
	if err != nil {
		t.Fatalf("EnforceWithTerritory error: %v", err)
	}
	if !ok {
		t.Error("rw_officer should be allowed to read tenant (uses * domain)")
	}

	// RW officer writing should also succeed
	ok, err = ce.EnforceWithTerritory("rw_officer", "income", "write", "rt-02", true)
	if err != nil {
		t.Fatalf("EnforceWithTerritory error: %v", err)
	}
	if !ok {
		t.Error("rw_officer should be allowed to write income (uses * domain)")
	}
}

func TestCasbinEnforcer_ResetEnforcerForTest(t *testing.T) {
	// Initialize first
	enforcerOnce = sync.Once{}
	enforcerInstance = nil
	enforcerErr = nil

	_, err := InitEnforcer("../../../rbac_model.conf", "../../../policy.csv")
	if err != nil {
		t.Fatalf("InitEnforcer failed: %v", err)
	}

	// Verify it's initialized
	if enforcerInstance == nil {
		t.Fatal("enforcer should be initialized")
	}

	// Reset
	ResetEnforcerForTest()

	// Verify it's reset
	if enforcerInstance != nil {
		t.Error("enforcer should be nil after reset")
	}
	if enforcerErr != nil {
		t.Error("enforcerErr should be nil after reset")
	}
}

func TestCasbinEnforcer_Singleton(t *testing.T) {
	// Reset singleton
	enforcerOnce = sync.Once{}
	enforcerInstance = nil
	enforcerErr = nil

	ce1, err := InitEnforcer("../../../rbac_model.conf", "../../../policy.csv")
	if err != nil {
		t.Fatalf("InitEnforcer failed: %v", err)
	}

	ce2 := GetEnforcer()
	if ce2 == nil {
		t.Fatal("expected non-nil enforcer after initialization")
	}

	// Both should wrap the same underlying enforcer
	if ce1.enforcer != ce2.enforcer {
		t.Error("expected same enforcer instance (singleton)")
	}
}

// newEnforcerForTerritoryTest creates an enforcer with territory substitution for testing.
// This is a duplicate of newEnforcerForTerritory from casbin_policy_test.go to avoid
// cross-file test dependencies, but uses the CasbinEnforcer wrapper.
func newEnforcerForTerritoryTest(t *testing.T, territoryID string) *CasbinEnforcer {
	t.Helper()

	// Use the test helper from casbin_policy_test.go
	e := newEnforcerForTerritory(t, territoryID)
	return &CasbinEnforcer{enforcer: e}
}
