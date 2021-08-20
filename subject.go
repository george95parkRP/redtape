package redtape

type Subject interface {
	Type() string
	// Metadata() map[string]interface{}
	// Enforcer() EnforceFunc
	Conditions() Conditions
	// EffectiveRoles() []*role.Role
	// MatchRole(string) bool
}

type subject struct {
	typ string
	// roles      []*role.Role
	// meta       map[string]interface{}
	// enforcer   EnforceFunc
	conditions Conditions
}

// func (s *subject) Enforce(req *Request) error {
// 	for _, r := range s.EffectiveRoles() {
// 		if r.ID == "test" {
// 			return nil
// 		}
// 	}

// 	return NewErrRequestDeniedImplicit(errors.New("access denied because no roles match assignement"))
// }

func NewSubject(typ string, opts ...SubjectOption) (Subject, error) {
	options := NewSubjectOptions(opts...)

	sub := &subject{
		typ: typ,
		// roles:    options.Roles,
		// meta:     options.Meta,
		// enforcer: options.Enforcer,
	}

	conds, err := NewConditions(options.Conditions, nil)
	if err != nil {
		return nil, err
	}

	sub.conditions = conds

	return sub, nil
}

func (s *subject) Type() string {
	return s.typ
}

// func (s *subject) Metadata() map[string]interface{} {
// 	return s.meta
// }

// func (s *subject) Enforcer() EnforceFunc {
// 	return s.enforcer
// }

func (s *subject) Conditions() Conditions {
	return s.conditions
}

// func (s *subject) EffectiveRoles() []*role.Role {
// 	var er []*role.Role
// 	for _, r := range s.roles {
// 		er = append(er, r.EffectiveRoles()...)
// 	}

// 	return er
// }

// func (s *subject) MatchRole(roleID string) bool {
// 	for _, r := range s.EffectiveRoles() {
// 		if match.Wildcard(roleID, r.ID) {
// 			return true
// 		}
// 	}

// 	return false
// }

type SubjectOptions struct {
	// Enforcer   EnforceFunc
	// Meta       map[string]interface{}
	// Roles      []*role.Role
	Conditions []ConditionOptions
}

type SubjectOption func(*SubjectOptions)

func NewSubjectOptions(opts ...SubjectOption) SubjectOptions {
	options := SubjectOptions{}

	for _, o := range opts {
		o(&options)
	}

	// if options.Meta == nil {
	// 	options.Meta = make(map[string]interface{})
	// }

	return options
}

// func SubjectEnforcer(e EnforceFunc) SubjectOption {
// 	return func(o *SubjectOptions) {
// 		o.Enforcer = e
// 	}
// }

// func SubjectMeta(meta ...map[string]interface{}) SubjectOption {
// 	return func(o *SubjectOptions) {
// 		if o.Meta == nil {
// 			o.Meta = make(map[string]interface{})
// 		}

// 		for _, md := range meta {
// 			for k, v := range md {
// 				o.Meta[k] = v
// 			}
// 		}
// 	}
// }

// func SubjectRole(role ...*role.Role) SubjectOption {
// 	return func(o *SubjectOptions) {
// 		o.Roles = append(o.Roles, role...)
// 	}
// }

func WithConditions(co ...ConditionOptions) SubjectOption {
	return func(o *SubjectOptions) {
		o.Conditions = append(o.Conditions, co...)
	}
}
