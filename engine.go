package redtape

import (
	"errors"
	"fmt"
)

type Engine interface {
	Verify(*Request) error
	Grant(Policy) error
	Revoke(string) error
	List() ([]Policy, error)

	RegisterCondition(string, ConditionBuilder) error
	ConditionRegistry() ConditionRegistry
}

type engine struct {
	mgr      PolicyManager
	enforcer Enforcer
	auditor  Auditor
	conds    ConditionRegistry
}

func NewEngine(mgr PolicyManager, opts ...EngineOption) Engine {
	return newEngine(mgr, opts...)
}

func newEngine(mgr PolicyManager, opts ...EngineOption) *engine {
	options := NewEngineOptions(opts...)

	return &engine{
		mgr:      mgr,
		enforcer: NewEnforcer(mgr, options.Matcher, options.Auditor),
		auditor:  options.Auditor,
		conds:    options.Registry,
	}
}

func (e *engine) Verify(req *Request) error {
	return e.enforcer.Enforce(req)
}

func (e *engine) Grant(p Policy) error {
	return e.mgr.Create(p)
}

func (e *engine) Revoke(p string) error {
	return e.mgr.Delete(p)
}

func (e *engine) List() ([]Policy, error) {
	return e.mgr.All()
}

func (e *engine) RegisterCondition(name string, ctor ConditionBuilder) error {
	if val, ok := e.conds[name]; ok {
		return errors.New(fmt.Sprintf("Condition already in registry with key %s and value %v", name, val))
	}

	e.conds[name] = ctor

	return nil
}

func (e *engine) ConditionRegistry() ConditionRegistry {
	return e.conds
}

type EngineOptions struct {
	Auditor  Auditor
	Matcher  Matcher
	Registry ConditionRegistry
}

type EngineOption func(*EngineOptions)

func NewEngineOptions(opts ...EngineOption) EngineOptions {
	options := EngineOptions{
		Auditor: NewConsoleAuditor(AuditAll),
		Matcher: NewRegexMatcher(),
	}

	for _, o := range opts {
		o(&options)
	}

	if options.Registry == nil {
		options.Registry = NewConditionRegistry()
	}

	return options
}

func EngineAuditor(a Auditor) EngineOption {
	return func(o *EngineOptions) {
		o.Auditor = a
	}
}

func EngineConditions(c ConditionRegistry) EngineOption {
	return func(o *EngineOptions) {
		o.Registry = c
	}
}
