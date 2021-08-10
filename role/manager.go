package role

import (
	"fmt"
	"sort"
	"sync"
)

// RoleManager provides methods to store and retrieve role sets.
type Manager interface {
	Create(*Role) error
	Update(*Role) error
	Get(string) (*Role, error)
	GetByName(string) (*Role, error)
	Delete(string) error
	All(limit, offset int) ([]*Role, error)

	GetMatching(string) ([]*Role, error)
}

type defaultRoleManager struct {
	roles map[string]*Role
	mu    sync.RWMutex
}

func NewRoleManager() Manager {
	return &defaultRoleManager{
		roles: make(map[string]*Role),
	}
}

func (m *defaultRoleManager) Create(r *Role) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.roles[r.ID]; exists {
		return fmt.Errorf("role %s already registered", r.ID)
	}

	m.roles[r.ID] = r

	return nil
}

func (m *defaultRoleManager) Update(r *Role) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.roles[r.ID] = r

	return nil
}

func (m *defaultRoleManager) Get(id string) (*Role, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	r, ok := m.roles[id]
	if !ok {
		return nil, fmt.Errorf("role %s does not exist", id)
	}

	return r, nil
}

func (m *defaultRoleManager) GetByName(name string) (*Role, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, r := range m.roles {
		if name == r.Name {
			return r, nil
		}
	}

	return nil, fmt.Errorf("role %s does not exist", name)
}

func (m *defaultRoleManager) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.roles, id)
	return nil
}

func (m *defaultRoleManager) All(limit, offset int) ([]*Role, error) {
	m.mu.RLock()

	rkeys := make([]string, len(m.roles))
	i := 0
	for k := range m.roles {
		rkeys[i] = k
		i++
	}

	start, end := limitIndices(limit, offset, len(m.roles))
	sort.Strings(rkeys)

	roles := make([]*Role, 0, len(rkeys[start:end]))
	for _, r := range rkeys[start:end] {
		roles = append(roles, m.roles[r])
	}

	m.mu.RUnlock()

	return roles, nil
}

func (m *defaultRoleManager) GetMatching(id string) ([]*Role, error) {
	panic("not implemented")
}

func limitIndices(limit, offset, length int) (int, int) {
	if offset > length {
		return length, length
	}

	if limit+offset > length {
		return offset, length
	}

	return offset, offset + limit
}
