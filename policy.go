package redtape

import (
	"context"
	"encoding/json"

	"github.com/fatih/structs"
)

// PolicyEffect type is returned by Enforcer to describe the outcome of a policy evaluation.
type PolicyEffect string

const (
	// PolicyEffectAllow indicates explicit permission of the request.
	PolicyEffectAllow PolicyEffect = "allow"
	// PolicyEffectDeny indicates explicit denial of the request.
	PolicyEffectDeny PolicyEffect = "deny"
)

// NewPolicyEffect returns a PolicyEffect for a given string.
func NewPolicyEffect(s string) PolicyEffect {
	switch s {
	case "allow":
		return PolicyEffectAllow
	case "deny":
		return PolicyEffectDeny
	default:
		return PolicyEffectDeny
	}
}

// Policy provides methods to return data about a configured policy.
type Policy interface {
	ID() string
	Name() string
	Description() string
	Subjects() []Subject
	Resources() []string
	Actions() []string
	Scopes() []string
	Effect() PolicyEffect
	Context() context.Context
}

type policy struct {
	id        string
	name      string
	desc      string
	subjects  []Subject
	resources []string
	actions   []string
	scopes    []string
	effect    PolicyEffect
	ctx       context.Context
}

// NewPolicy returns a default policy implementation from a set of provided options.
func NewPolicy(opts ...PolicyOption) (Policy, error) {
	o := NewPolicyOptions(opts...)

	p := &policy{
		id:        o.ID,
		name:      o.Name,
		desc:      o.Description,
		subjects:  o.Subjects,
		resources: o.Resources,
		actions:   o.Actions,
		scopes:    o.Scopes,
		effect:    NewPolicyEffect(o.Effect),
		ctx:       o.Context,
	}

	return p, nil
}

// MustNewPolicy returns a default policy implementation or panics on error.
func MustNewPolicy(opts ...PolicyOption) Policy {
	p, err := NewPolicy(opts...)

	if err != nil {
		panic("failed to create new policy: " + err.Error())
	}

	return p
}

// MarshalJSON returns a JSON byte slice representation of the default policy implementation.
func (p *policy) MarshalJSON() ([]byte, error) {
	opts := PolicyOptions{
		ID:          p.id,
		Name:        p.name,
		Description: p.desc,
		Subjects:    p.subjects,
		Resources:   p.resources,
		Actions:     p.actions,
		Scopes:      p.scopes,
		Effect:      string(p.effect),
	}

	structs.DefaultTagName = "json"

	return json.Marshal(opts)
}

// ID returns the policy ID.
func (p *policy) ID() string {
	return p.id
}

// Name returns the policy Name.
func (p *policy) Name() string {
	return p.name
}

// Description returns the policy Description.
func (p *policy) Description() string {
	return p.desc
}

// Roles returns the roles the policy applies to.
func (p *policy) Subjects() []Subject {
	return p.subjects
}

// Resources returns the resources the policy applies to.
func (p *policy) Resources() []string {
	return p.resources
}

// Actions returns the actions the policy applies to.
func (p *policy) Actions() []string {
	return p.actions
}

// Scopes returns the scopes the policy applies to.
func (p *policy) Scopes() []string {
	return p.scopes
}

func (p *policy) Context() context.Context {
	return p.ctx
}

// Effect returns the configured PolicyEffect.
func (p *policy) Effect() PolicyEffect {
	return p.effect
}

// PolicyOptions struct allows different Policy implementations to be configured with marshalable data.
type PolicyOptions struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Subjects    []Subject       `json:"roles"`
	Resources   []string        `json:"resources"`
	Actions     []string        `json:"actions"`
	Scopes      []string        `json:"scopes"`
	Effect      string          `json:"effect"`
	Context     context.Context `json:"-"`
}

// PolicyOption is a typed function allowing updates to PolicyOptions through functional options.
type PolicyOption func(*PolicyOptions)

// NewPolicyOptions returns PolicyOptions configured with the provided functional options.
func NewPolicyOptions(opts ...PolicyOption) PolicyOptions {
	options := PolicyOptions{}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// SetPolicyOptions is a PolicyOption setting all PolicyOptions to the provided values.
func SetPolicyOptions(opts PolicyOptions) PolicyOption {
	return func(o *PolicyOptions) {
		*o = opts
	}
}

// PolicyID sets the policy ID Option.
func PolicyID(id string) PolicyOption {
	return func(o *PolicyOptions) {
		o.ID = id
	}
}

// PolicyName sets the policy Name Option.
func PolicyName(n string) PolicyOption {
	return func(o *PolicyOptions) {
		o.Name = n
	}
}

// PolicyDescription sets the policy description Option.
func PolicyDescription(d string) PolicyOption {
	return func(o *PolicyOptions) {
		o.Description = d
	}
}

func SetPolicyEffect(s string) PolicyOption {
	return func(o *PolicyOptions) {
		o.Effect = s
	}
}

// PolicyDeny sets the PolicyEffect to deny.
func PolicyDeny() PolicyOption {
	return func(o *PolicyOptions) {
		o.Effect = "deny"
	}
}

// PolicyAllow sets the PolicyEffect to allow.
func PolicyAllow() PolicyOption {
	return func(o *PolicyOptions) {
		o.Effect = "allow"
	}
}

// SetResources replaces the option Resources with the provided values.
func SetResources(s ...string) PolicyOption {
	return func(o *PolicyOptions) {
		o.Resources = s
	}
}

// SetActions replaces the option Actions with the provided values.
func SetActions(s ...string) PolicyOption {
	return func(o *PolicyOptions) {
		o.Actions = s
	}
}

func SetScopes(s ...string) PolicyOption {
	return func(o *PolicyOptions) {
		o.Scopes = s
	}
}

// SetContext sets the Context option.
func SetContext(ctx context.Context) PolicyOption {
	return func(o *PolicyOptions) {
		o.Context = ctx
	}
}

// WithSubject adds a Subject to the Subjects option.
func WithSubjects(s ...Subject) PolicyOption {
	return func(o *PolicyOptions) {
		o.Subjects = s
	}
}
