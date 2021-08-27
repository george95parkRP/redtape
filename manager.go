package redtape

import (
	"fmt"
	"sync"
	"time"

	"github.com/blushft/redtape/match"
)

// PolicyManager contains methods to allow query, update, and removal of policies.
type PolicyManager interface {
	Create(Policy) error
	Update(Policy) error
	Get(string) (Policy, error)
	Delete(string) error
	All() ([]Policy, error)

	FindByRequest(*Request) ([]Policy, error)
	FindByRole(string) ([]Policy, error)
	FindByResource(string) ([]Policy, error)
	FindByScope(string) ([]Policy, error)
}

type defaultPolicyManager struct {
	policies map[string]Policy
	mu       sync.RWMutex
}

// NewPolicyManager returns a default memory backed policy manager.
func NewPolicyManager() PolicyManager {
	return newPolicyManager()
}

func newPolicyManager() *defaultPolicyManager {
	return &defaultPolicyManager{
		policies: make(map[string]Policy),
	}
}

// Create adds a policy to the manager.
func (m *defaultPolicyManager) Create(p Policy) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.policies[p.ID()]; exists {
		return fmt.Errorf("policy %s already registered", p.ID())
	}

	m.policies[p.ID()] = p

	return nil
}

// Update replaces a named policy with the provided policy.
func (m *defaultPolicyManager) Update(p Policy) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.policies[p.ID()] = p

	return nil
}

// Get retrieves a policy by id or error if one does not exist.
func (m *defaultPolicyManager) Get(id string) (Policy, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, ok := m.policies[id]
	if !ok {
		return nil, fmt.Errorf("policy %s does not exist", id)
	}

	return p, nil
}

// Delete removes a policy by id.
func (m *defaultPolicyManager) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.policies, id)
	return nil
}

// All returns a slice containing all policies.
// TODO: refactor the extra findall func
func (m *defaultPolicyManager) All() ([]Policy, error) {
	return m.findAll()
}

func (m *defaultPolicyManager) findAll() ([]Policy, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ps := make([]Policy, 0, len(m.policies))
	for _, p := range m.policies {
		ps = append(ps, p)
	}

	return ps, nil
}

// FindByRequest returns all policies matching a Request.
func (m *defaultPolicyManager) FindByRequest(r *Request) ([]Policy, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ps := []Policy{}
	for _, p := range m.policies {
		if found := findByResource(r.Resource, p.Resources()); !found {
			continue
		}

		if found := findByScope(r.Scope, p.Scopes()); !found {
			continue
		}

		if found := findByAction(r.Action, p.Actions()); !found {
			continue
		}

		ps = append(ps, p)
	}

	return ps, nil
}

func findByAction(act string, actions []string) bool {
	for _, a := range actions {
		if match.Wildcard(a, act) {
			return true
		}
	}

	return false
}

// FindByRole returns all policies matching a Role.
func (m *defaultPolicyManager) FindByRole(_ string) ([]Policy, error) {
	return m.findAll()
}

// FindByResource returns all policies matching a Resource.
func (m *defaultPolicyManager) FindByResource(res string) ([]Policy, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ps := []Policy{}
	for _, p := range m.policies {
		if found := findByResource(res, p.Resources()); found {
			ps = append(ps, p)
			break
		}
	}

	return ps, nil
}

func findByResource(res string, resources []string) bool {
	for _, r := range resources {
		if match.Wildcard(r, res) {
			return true
		}
	}

	return false
}

// FindByResource returns all policies matching a Resource.
func (m *defaultPolicyManager) FindByScope(scope string) ([]Policy, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ps := []Policy{}
	for _, p := range m.policies {
		if found := findByScope(scope, p.Scopes()); found {
			ps = append(ps, p)
			break
		}
	}

	return ps, nil
}

func findByScope(scope string, scopes []string) bool {
	for _, s := range scopes {
		if match.Wildcard(s, scope) {
			return true
		}
	}

	return false
}

type policyCache struct {
	mgr   PolicyManager
	cache *defaultPolicyManager
	ttl   time.Time
	exp   time.Duration
}

func NewPolicyCache(mgr PolicyManager, exp time.Duration) *policyCache {
	return &policyCache{
		mgr:   mgr,
		cache: newPolicyManager(),
		exp:   exp,
		ttl:   time.Now().Add(exp),
	}
}

func (c *policyCache) resetCache() {
	c.ttl = time.Now().Add(c.exp)
	c.cache = newPolicyManager()
}

func (c *policyCache) checkExp() bool {
	if time.Now().After(c.ttl) {
		c.resetCache()
		return true
	}

	return false
}

func (c *policyCache) Create(p Policy) error {
	if err := c.mgr.Create(p); err != nil {
		return err
	}

	return c.cache.Create(p)
}

func (c *policyCache) Update(p Policy) error {
	if err := c.mgr.Update(p); err != nil {
		return err
	}

	return c.cache.Update(p)
}

func (c *policyCache) Get(id string) (Policy, error) {
	if !c.checkExp() {
		p, err := c.cache.Get(id)
		if err == nil {
			return p, nil
		}
	}

	p, err := c.mgr.Get(id)
	if err != nil {
		return nil, err
	}

	if err := c.cache.Create(p); err != nil {
		return nil, err
	}

	return p, nil
}
