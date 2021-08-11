package redtape

import (
	"github.com/blushft/redtape/match"
	"github.com/blushft/redtape/role"
)

type Subject interface {
	ID() string
	Type() string
	Enforcer() EnforceFunc
	EffectiveRoles() []*role.Role
	MatchRole(string) bool
}

type subject struct {
	Id           string                 `json:"id"`
	Name         string                 `json:"name"`
	Roles        []*role.Role           `json:"roles"`
	Meta         map[string]interface{} `json:"meta"`
	EnforcerFunc EnforceFunc            `json:"enforcer"`
}

// func (s *subject) Enforce(req *Request) error {
// 	for _, r := range s.EffectiveRoles() {
// 		if r.ID == "test" {
// 			return nil
// 		}
// 	}

// 	return NewErrRequestDeniedImplicit(errors.New("access denied because no roles match assignement"))
// }

func NewSubject(id string, opts ...SubjectOption) Subject {
	options := NewSubjectOptions(opts...)

	sub := &subject{
		Id:           id,
		Name:         options.Name,
		Roles:        options.Roles,
		EnforcerFunc: options.Enforcer,
	}

	if options.Meta == nil {
		sub.Meta = make(map[string]interface{})
	}

	return sub
}

func (s *subject) ID() string {
	return s.Id
}

func (s *subject) Type() string {
	return ""
}

func (s *subject) Enforcer() EnforceFunc {
	return s.EnforcerFunc
}

func (s *subject) EffectiveRoles() []*role.Role {
	var er []*role.Role
	for _, r := range s.Roles {
		er = append(er, r.EffectiveRoles()...)
	}

	return er
}

func (s *subject) MatchRole(roleID string) bool {
	for _, r := range s.EffectiveRoles() {
		if match.Wildcard(roleID, r.ID) {
			return true
		}
	}

	return false
}

type SubjectOptions struct {
	Enforcer EnforceFunc
	Name     string
	Meta     map[string]interface{}
	Roles    []*role.Role
}

type SubjectOption func(*SubjectOptions)

func NewSubjectOptions(opts ...SubjectOption) SubjectOptions {
	options := SubjectOptions{}

	for _, o := range opts {
		o(&options)
	}

	return options
}

func SubjectEnforcer(e EnforceFunc) SubjectOption {
	return func(o *SubjectOptions) {
		o.Enforcer = e
	}
}

func SubjectName(n string) SubjectOption {
	return func(o *SubjectOptions) {
		o.Name = n
	}
}

func SubjectRole(role ...*role.Role) SubjectOption {
	return func(o *SubjectOptions) {
		o.Roles = append(o.Roles, role...)
	}
}

func SubjectMeta(meta ...map[string]interface{}) SubjectOption {
	return func(o *SubjectOptions) {
		if o.Meta == nil {
			o.Meta = make(map[string]interface{})
		}

		for _, md := range meta {
			for k, v := range md {
				o.Meta[k] = v
			}
		}
	}
}
