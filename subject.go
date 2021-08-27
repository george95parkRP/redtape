package redtape

import "github.com/blushft/redtape/match"

type Subject interface {
	Type() string
	Conditions() Conditions
	MatchSubject(string) bool
}

type subject struct {
	typ        string
	conditions Conditions
}

func NewSubject(typ string, opts ...SubjectOption) (Subject, error) {
	options := NewSubjectOptions(opts...)

	sub := &subject{
		typ: typ,
	}

	conds, err := NewConditions(options.Conditions, options.ConditionRegistry)
	if err != nil {
		return nil, err
	}

	sub.conditions = conds

	return sub, nil
}

func (s *subject) Type() string {
	return s.typ
}

func (s *subject) Conditions() Conditions {
	return s.conditions
}

func (s *subject) MatchSubject(val string) bool {
	if match.Wildcard(s.typ, val) {
		return true
	}

	return false
}

type SubjectOptions struct {
	Conditions        []ConditionOptions
	ConditionRegistry ConditionRegistry
}

type SubjectOption func(*SubjectOptions)

func NewSubjectOptions(opts ...SubjectOption) SubjectOptions {
	options := SubjectOptions{}

	for _, o := range opts {
		o(&options)
	}

	return options
}

func WithConditions(co ...ConditionOptions) SubjectOption {
	return func(o *SubjectOptions) {
		o.Conditions = co
	}
}

func WithConditionRegistry(cr ConditionRegistry) SubjectOption {
	return func(o *SubjectOptions) {
		o.ConditionRegistry = cr
	}
}
