package file

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/blushft/redtape"
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

func New(opts ...Option) *Manager {
	return &Manager{
		options: NewOptions(opts...),
	}
}

func (f *Manager) policyPath() string {
	fn := fmt.Sprintf("%s.policy", f.options.Name)
	return filepath.Join(f.options.Path, fn)
}

func (f *Manager) loadPolicies() (map[string]redtape.Policy, error) {
	b, err := os.ReadFile(f.policyPath())
	if err != nil {
		return nil, err
	}

	m := make(map[string]redtape.Policy)
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	return m, nil
}

func (f *Manager) savePolicies(m map[string]redtape.Policy) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return os.WriteFile(f.policyPath(), b, os.ModePerm)
}

func (f *Manager) Create(p redtape.Policy) error {
	panic("not implemented") // TODO: Implement
}

func (f *Manager) Update(_ redtape.Policy) error {
	panic("not implemented") // TODO: Implement
}

func (f *Manager) Get(_ string) (redtape.Policy, error) {
	panic("not implemented") // TODO: Implement
}

func (f *Manager) Delete(_ string) error {
	panic("not implemented") // TODO: Implement
}

func (f *Manager) All(limit int, offset int) ([]redtape.Policy, error) {
	panic("not implemented") // TODO: Implement
}

func (f *Manager) FindByRequest(_ *redtape.Request) ([]redtape.Policy, error) {
	panic("not implemented") // TODO: Implement
}

func (f *Manager) FindByRole(_ string) ([]redtape.Policy, error) {
	panic("not implemented") // TODO: Implement
}

func (f *Manager) FindByResource(_ string) ([]redtape.Policy, error) {
	panic("not implemented") // TODO: Implement
}

func (f *Manager) FindByScope(_ string) ([]redtape.Policy, error) {
	panic("not implemented") // TODO: Implement
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
