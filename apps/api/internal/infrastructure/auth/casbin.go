package auth

import (
	"fmt"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

var (
	enforcerInstance *casbin.Enforcer
	enforcerOnce     sync.Once
	enforcerErr      error
)

// CasbinEnforcer wraps the Casbin enforcer with territory-aware enforcement.
type CasbinEnforcer struct {
	enforcer *casbin.Enforcer
}

// InitEnforcer initializes the Casbin enforcer (singleton pattern).
// It loads the RBAC model from rbac_model.conf and policies from policy.csv.
// This should be called once during application startup.
func InitEnforcer(modelPath, policyPath string) (*CasbinEnforcer, error) {
	enforcerOnce.Do(func() {
		m, err := model.NewModelFromFile(modelPath)
		if err != nil {
			enforcerErr = fmt.Errorf("failed to load Casbin model: %w", err)
			return
		}

		adapter := fileadapter.NewAdapter(policyPath)
		e, err := casbin.NewEnforcer(m, adapter)
		if err != nil {
			enforcerErr = fmt.Errorf("failed to create Casbin enforcer: %w", err)
			return
		}

		err = e.LoadPolicy()
		if err != nil {
			enforcerErr = fmt.Errorf("failed to load Casbin policy: %w", err)
			return
		}

		enforcerInstance = e
	})

	if enforcerErr != nil {
		return nil, enforcerErr
	}

	return &CasbinEnforcer{enforcer: enforcerInstance}, nil
}

// GetEnforcer returns the initialized CasbinEnforcer instance.
// Returns nil if InitEnforcer has not been called.
func GetEnforcer() *CasbinEnforcer {
	if enforcerInstance == nil {
		return nil
	}
	return &CasbinEnforcer{enforcer: enforcerInstance}
}

// Enforce checks if the given role can perform the action on the resource
// within the specified territory domain.
func (ce *CasbinEnforcer) Enforce(role, resource, action, territory string) (bool, error) {
	return ce.enforcer.Enforce(role, resource, action, territory)
}

// EnforceWithTerritory substitutes {{territory_id}} in the policy with the
// user's territory ID before checking permissions. For RW officers, uses "*"
// to allow access to all territories.
func (ce *CasbinEnforcer) EnforceWithTerritory(role, resource, action, territory string, isRWOfficer bool) (bool, error) {
	domain := territory
	if isRWOfficer {
		domain = "*"
	}
	return ce.enforcer.Enforce(role, resource, action, domain)
}

// AddPolicy adds a new policy rule at runtime.
func (ce *CasbinEnforcer) AddPolicy(role, resource, action, territory string) error {
	_, err := ce.enforcer.AddPolicy(role, resource, action, territory)
	if err != nil {
		return fmt.Errorf("failed to add policy: %w", err)
	}
	return nil
}

// RemovePolicy removes a policy rule at runtime.
func (ce *CasbinEnforcer) RemovePolicy(role, resource, action, territory string) error {
	_, err := ce.enforcer.RemovePolicy(role, resource, action, territory)
	if err != nil {
		return fmt.Errorf("failed to remove policy: %w", err)
	}
	return nil
}

// AddRoleLink adds a role inheritance link (e.g., "rw_officer" inherits "rt_officer").
func (ce *CasbinEnforcer) AddRoleLink(role, inheritedRole string) error {
	_, err := ce.enforcer.AddGroupingPolicy(role, inheritedRole)
	if err != nil {
		return fmt.Errorf("failed to add role link: %w", err)
	}
	return nil
}

// RemoveRoleLink removes a role inheritance link.
func (ce *CasbinEnforcer) RemoveRoleLink(role, inheritedRole string) error {
	_, err := ce.enforcer.RemoveGroupingPolicy(role, inheritedRole)
	if err != nil {
		return fmt.Errorf("failed to remove role link: %w", err)
	}
	return nil
}
