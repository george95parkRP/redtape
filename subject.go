package redtape

import (
	"errors"

	"github.com/blushft/redtape/role"
)

type Subject interface {
	Type() string
	Enforcer() EnforceFunc
}

type subject struct {
	ID    string                 `json:"id"`
	Name  string                 `json:"name"`
	Roles []*role.Role           `json:"roles"`
	Meta  map[string]interface{} `json:"meta"`
}

func (s *subject) Enforce(req *Request) error {
	for _, r := range s.EffectiveRoles() {
		if r.ID == "test" {
			return nil
		}
	}

	return NewErrRequestDeniedImplicit(errors.New("access denied because no roles match assignement"))
}

func NewSubject(id string, opts ...SubjectOption) Subject {
	sub := &subject{
		ID:   id,
		Meta: make(map[string]interface{}),
	}

	for _, opt := range opts {
		opt(sub)
	}

	return sub
}

func (s *subject) EffectiveRoles() []*role.Role {
	var er []*role.Role
	for _, r := range s.Roles {
		er = append(er, r.EffectiveRoles()...)
	}

	return er
}

func (s *subject) String() string {
	return s.ID
}

type SubjectOption func(*Subject)

func SubjectName(n string) SubjectOption {
	return func(s *Subject) {
		s.Name = n
	}
}

func SubjectRole(role ...*role.Role) SubjectOption {
	return func(s *Subject) {
		s.Roles = append(s.Roles, role...)
	}
}

func SubjectMeta(meta ...map[string]interface{}) SubjectOption {
	return func(s *Subject) {
		if s.Meta == nil {
			s.Meta = make(map[string]interface{})
		}

		for _, md := range meta {
			for k, v := range md {
				s.Meta[k] = v
			}
		}
	}
}
