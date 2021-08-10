package redtape

type Engine interface {
	Verify(*Request) error
	Grant(Policy) error
	Revoke(string) error
	List() ([]Policy, error)

	// TODO: allow adds
	//RegisterCondition(name string, ctor ConditionBuilder)
}

type engine struct {
	mgr      PolicyManager
	enforcer Enforcer
	auditor  Auditor
	cond     ConditionRegistry
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
		cond:     options.Registry,
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
