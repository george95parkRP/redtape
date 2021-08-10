package file

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/blushft/redtape/role"
)

type Options struct {
	Name string
	Path string
}

type Option func(*Options)

func NewOptions(opts ...Option) Options {
	o := Options{
		Name: "redtape",
	}

	for _, opt := range opts {
		opt(&o)
	}

	return o
}

type Manager struct {
	options Options
}

func New(opts ...Option) (*Manager, error) {
	f := &Manager{
		options: NewOptions(opts...),
	}

	if !fileExists(f.RolePath()) {
		if err := os.WriteFile(f.RolePath(), []byte("{}"), os.ModePerm); err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (f *Manager) RolePath() string {
	fn := fmt.Sprintf("%s.roles", f.options.Name)
	return filepath.Join(f.options.Path, fn)
}

func (f *Manager) loadRoles() (map[string]*role.Role, error) {
	b, err := os.ReadFile(f.RolePath())
	if err != nil {
		return nil, err
	}

	m := make(map[string]*role.Role)
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	return m, nil
}

func (f *Manager) saveRoles(roles map[string]*role.Role) error {
	b, err := json.Marshal(roles)
	if err != nil {
		return err
	}

	return os.WriteFile(f.RolePath(), b, os.ModePerm)
}

func (f *Manager) Create(role *role.Role) error {
	return f.writeRole(role, false)
}

func (f *Manager) Update(role *role.Role) error {
	return f.writeRole(role, true)
}

func (f *Manager) writeRole(role *role.Role, overwrite bool) error {
	m, err := f.loadRoles()
	if err != nil {
		return err
	}

	_, ok := m[role.ID]
	if ok && !overwrite {
		return fmt.Errorf("role %s already registered", role.ID)
	}

	m[role.ID] = role

	return f.saveRoles(m)
}

func (f *Manager) Get(id string) (*role.Role, error) {
	m, err := f.loadRoles()
	if err != nil {
		return nil, err
	}

	r, ok := m[id]
	if !ok {
		return nil, fmt.Errorf("role %s not found", id)
	}

	return r, nil
}

func (f *Manager) GetByName(name string) (*role.Role, error) {
	m, err := f.loadRoles()
	if err != nil {
		return nil, err
	}

	for _, r := range m {
		if r.Name == name {
			return r, nil
		}
	}

	return nil, fmt.Errorf("role name %s not found", name)
}

func (f *Manager) Delete(id string) error {
	m, err := f.loadRoles()
	if err != nil {
		return err
	}

	delete(m, id)

	return f.saveRoles(m)
}

func (f *Manager) All(limit, offset int) ([]*role.Role, error) {
	m, err := f.loadRoles()
	if err != nil {
		return nil, err
	}

	rkeys := make([]string, len(m))
	i := 0
	for k := range m {
		rkeys[i] = k
		i++
	}

	start, end := limitIndices(limit, offset, len(m))
	sort.Strings(rkeys)

	roles := make([]*role.Role, 0, len(rkeys[start:end]))
	for _, r := range rkeys[start:end] {
		roles = append(roles, m[r])
	}

	return roles, nil
}

func (f *Manager) GetMatching(_ string) ([]*role.Role, error) {
	panic("not implemented") // TODO: Implement
}

func fileExists(fp string) bool {
	i, err := os.Stat(fp)
	if err != nil || err == os.ErrNotExist {
		return false
	}

	if i.IsDir() {
		return false
	}

	return true

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
