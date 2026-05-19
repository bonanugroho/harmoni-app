package auth

import (
	"strings"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
)

func TestCasbinPolicy_RTOfficerOwnTerritory(t *testing.T) {
	// Create enforcer with rt-01 substituted for {{territory_id}}
	e := newEnforcerForTerritory(t, "rt-01")

	// RT officer should be allowed to read/write their own territory
	tests := []struct {
		role      string
		resource  string
		action    string
		territory string
		want      bool
	}{
		{"rt_officer", "tenant", "read", "rt-01", true},
		{"rt_officer", "tenant", "write", "rt-01", true},
		{"rt_officer", "income", "read", "rt-01", true},
		{"rt_officer", "income", "write", "rt-01", true},
		{"rt_officer", "expenditure", "read", "rt-01", true},
		{"rt_officer", "expenditure", "write", "rt-01", true},
		{"rt_officer", "report", "read", "rt-01", true},
		{"rt_officer", "report", "write", "rt-01", true},
	}

	for _, tt := range tests {
		ok, err := e.Enforce(tt.role, tt.resource, tt.action, tt.territory)
		if err != nil {
			t.Errorf("Enforce(%s, %s, %s, %s) error: %v", tt.role, tt.resource, tt.action, tt.territory, err)
			continue
		}
		if ok != tt.want {
			t.Errorf("Enforce(%s, %s, %s, %s) = %v, want %v", tt.role, tt.resource, tt.action, tt.territory, ok, tt.want)
		}
	}
}

func TestCasbinPolicy_RTOfficerOtherTerritory(t *testing.T) {
	// Create enforcer with rt-01 substituted for {{territory_id}}
	e := newEnforcerForTerritory(t, "rt-01")

	// RT officer should NOT be allowed to access other territories
	ok, err := e.Enforce("rt_officer", "tenant", "read", "rt-02")
	if err != nil {
		t.Fatalf("Enforce error: %v", err)
	}
	if ok {
		t.Error("rt_officer should NOT have access to rt-02 territory")
	}
}

func TestCasbinPolicy_RWOfficerAllTerritories(t *testing.T) {
	// Create enforcer with rt-01 substituted (RW uses * so territory doesn't matter)
	e := newEnforcerForTerritory(t, "rt-01")

	// RW officer should be allowed to read/write ANY territory
	territories := []string{"rt-01", "rt-02", "rt-03", "rw-01"}
	resources := []string{"tenant", "income", "expenditure", "report"}
	actions := []string{"read", "write"}

	for _, territory := range territories {
		for _, resource := range resources {
			for _, action := range actions {
				ok, err := e.Enforce("rw_officer", resource, action, territory)
				if err != nil {
					t.Errorf("Enforce(rw_officer, %s, %s, %s) error: %v", resource, action, territory, err)
					continue
				}
				if !ok {
					t.Errorf("Enforce(rw_officer, %s, %s, %s) = %v, want true", resource, action, territory, ok)
				}
			}
		}
	}
}

func TestCasbinPolicy_ResidentReadOnly(t *testing.T) {
	// Create enforcer with rt-01 substituted for {{territory_id}}
	e := newEnforcerForTerritory(t, "rt-01")

	// Resident should be allowed to read own territory
	ok, err := e.Enforce("resident", "tenant", "read", "rt-01")
	if err != nil {
		t.Fatalf("Enforce error: %v", err)
	}
	if !ok {
		t.Error("resident should be allowed to read own territory")
	}

	// Resident should NOT be allowed to write
	ok, err = e.Enforce("resident", "tenant", "write", "rt-01")
	if err != nil {
		t.Fatalf("Enforce error: %v", err)
	}
	if ok {
		t.Error("resident should NOT be allowed to write")
	}
}

func TestCasbinPolicy_ResidentOtherTerritory(t *testing.T) {
	// Create enforcer with rt-01 substituted for {{territory_id}}
	e := newEnforcerForTerritory(t, "rt-01")

	// Resident should NOT be allowed to access other territory data
	ok, err := e.Enforce("resident", "tenant", "read", "rt-02")
	if err != nil {
		t.Fatalf("Enforce error: %v", err)
	}
	if ok {
		t.Error("resident should NOT have access to other territory")
	}
}

func TestCasbinPolicy_RoleInheritance(t *testing.T) {
	// Create enforcer with rt-01 substituted for {{territory_id}}
	e := newEnforcerForTerritory(t, "rt-01")

	// RW officer inherits RT officer permissions (should have rt_officer access too)
	ok, err := e.Enforce("rw_officer", "tenant", "read", "rt-01")
	if err != nil {
		t.Fatalf("Enforce error: %v", err)
	}
	if !ok {
		t.Error("rw_officer should inherit rt_officer permissions")
	}

	// RT officer inherits resident permissions (should have resident read access)
	ok, err = e.Enforce("rt_officer", "tenant", "read", "rt-01")
	if err != nil {
		t.Fatalf("Enforce error: %v", err)
	}
	if !ok {
		t.Error("rt_officer should inherit resident permissions")
	}
}

// newEnforcerForTerritory creates a Casbin enforcer with the given territory ID
// substituted for all {{territory_id}} placeholders in the policy CSV.
// This simulates what the middleware does at runtime.
func newEnforcerForTerritory(t *testing.T, territoryID string) *casbin.Enforcer {
	t.Helper()

	m, err := model.NewModelFromFile("../../../rbac_model.conf")
	if err != nil {
		t.Fatalf("failed to load model: %v", err)
	}

	// Build policy model string with substituted territory
	policyText := buildPolicyWithTerritory(territoryID)

	// Create a string adapter from the policy text
	adapter := &stringAdapter{policies: parsePolicyLines(policyText)}

	e, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		t.Fatalf("failed to create enforcer: %v", err)
	}

	err = e.LoadPolicy()
	if err != nil {
		t.Fatalf("failed to load policy: %v", err)
	}

	return e
}

// stringAdapter is a simple in-memory policy adapter for testing.
type stringAdapter struct {
	policies [][]string
}

func (a *stringAdapter) LoadPolicy(model model.Model) error {
	for _, rule := range a.policies {
		if len(rule) < 2 {
			continue
		}
		pType := rule[0]
		if pType == "p" && len(rule) >= 5 {
			model.AddPolicy("p", "p", rule[1:])
		} else if pType == "g" && len(rule) >= 3 {
			model.AddPolicy("g", "g", rule[1:])
		}
	}
	return nil
}

func (a *stringAdapter) SavePolicy(model model.Model) error {
	return nil
}

func (a *stringAdapter) AddPolicy(sec string, pType string, rule []string) error {
	return nil
}

func (a *stringAdapter) RemovePolicy(sec string, pType string, rule []string) error {
	return nil
}

func (a *stringAdapter) RemoveFilteredPolicy(sec string, pType string, fieldIndex int, fieldValues ...string) error {
	return nil
}

func (a *stringAdapter) RemovePolicies(sec string, pType string, rules [][]string) error {
	return nil
}

func (a *stringAdapter) AddPolicies(sec string, pType string, rules [][]string) error {
	return nil
}

func (a *stringAdapter) IsFiltered() bool {
	return false
}

func parsePolicyLines(text string) [][]string {
	var policies [][]string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, ",")
		for i, part := range parts {
			parts[i] = strings.TrimSpace(part)
		}
		policies = append(policies, parts)
	}
	return policies
}

func buildPolicyWithTerritory(territoryID string) string {
	return strings.ReplaceAll(`p, resident, tenant, read, {{territory_id}}
p, resident, income, read, {{territory_id}}
p, resident, expenditure, read, {{territory_id}}
p, resident, report, read, {{territory_id}}
p, rt_officer, tenant, read, {{territory_id}}
p, rt_officer, tenant, write, {{territory_id}}
p, rt_officer, income, read, {{territory_id}}
p, rt_officer, income, write, {{territory_id}}
p, rt_officer, expenditure, read, {{territory_id}}
p, rt_officer, expenditure, write, {{territory_id}}
p, rt_officer, report, read, {{territory_id}}
p, rt_officer, report, write, {{territory_id}}
p, rw_officer, tenant, read, *
p, rw_officer, tenant, write, *
p, rw_officer, income, read, *
p, rw_officer, income, write, *
p, rw_officer, expenditure, read, *
p, rw_officer, expenditure, write, *
p, rw_officer, report, read, *
p, rw_officer, report, write, *
g, rw_officer, rt_officer
g, rt_officer, resident`, "{{territory_id}}", territoryID)
}
