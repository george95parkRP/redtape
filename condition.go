package redtape

import (
	"fmt"
	"net"

	"github.com/mitchellh/mapstructure"
)

// ConditionBuilder is a typed function that returns a Condition.
type ConditionBuilder func() Condition

// ConditionRegistry is a map contiaining named ConditionBuilders.
type ConditionRegistry map[string]ConditionBuilder

// NewConditionRegistry returns a ConditionRegistry containing the default Conditions and accepts an array
// of map[string]ConditionBuilder to add custom conditions to the set.
func NewConditionRegistry(conds ...map[string]ConditionBuilder) ConditionRegistry {
	reg := ConditionRegistry{
		new(BoolCondition).Name(): func() Condition {
			return new(BoolCondition)
		},
		new(SubjectEqualsCondition).Name(): func() Condition {
			return new(SubjectEqualsCondition)
		},
		new(StringEqualsCondition).Name(): func() Condition {
			return new(StringEqualsCondition)
		},
		new(CIDRCondition).Name(): func() Condition {
			return new(CIDRCondition)
		},
	}

	for _, ce := range conds {
		for k, c := range ce {
			reg[k] = c
		}
	}

	return reg
}

// Condition is the interface allowing different types of conditional expressions.
type Condition interface {
	Name() string
	Meets(interface{}, *Request) bool
}

// Conditions is a map of named Conditions.
type Conditions map[string]Condition

// NewConditions accepts an array of options and an optional ConditionRegistry and returns a Conditions map.
func NewConditions(opts []ConditionOptions, reg ConditionRegistry) (Conditions, error) {
	if reg == nil {
		reg = NewConditionRegistry()
	}

	cond := make(map[string]Condition)

	for _, co := range opts {
		cf, ok := reg[co.Type]
		if !ok {
			return nil, fmt.Errorf("unknown condition type %s, is it registered?", co.Type)
		}

		nc := cf()
		if len(co.Options) > 0 {
			if err := mapstructure.Decode(co.Options, &nc); err != nil {
				return nil, err
			}
		}

		cond[co.Name] = nc
	}

	return cond, nil
}

func (c Conditions) Meets(r *Request) bool {
	meta := RequestMetadataFromContext(r.Context)
	for key, cond := range c {
		if ok := cond.Meets(meta[key], r); !ok {
			return false
		}
	}

	return true
}

// ConditionOptions contains the values used to build a Condition.
type ConditionOptions struct {
	Name    string                 `json:"name"`
	Type    string                 `json:"type"`
	Options map[string]interface{} `json:"options"`
}

// BoolCondition matches a boolean value from context to the preconfigured value.
type BoolCondition struct {
	Value bool `json:"value"`
}

// Name fulfills the Name method of Condition.
func (c *BoolCondition) Name() string {
	return "bool"
}

// Meets evaluates whether parameter val matches the Condition Value.
func (c *BoolCondition) Meets(val interface{}, _ *Request) bool {
	v, ok := val.(bool)

	return ok && v == c.Value
}

// SubjectEqualsCondition matches the Request subject against the required subject passed to the condition.
type SubjectEqualsCondition struct{}

// Name fulfills the Name method of Condition.
func (c *SubjectEqualsCondition) Name() string {
	return "subject_equals"
}

// Meets evaluates true when the subject val matches Request#Subject.
func (c *SubjectEqualsCondition) Meets(val interface{}, r *Request) bool {
	v, ok := val.([]string)
	if !ok {
		return false
	}

	for _, sub := range v {
		found := false
		for _, rs := range r.Subjects {
			if rs == sub {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

type StringEqualsCondition struct {
	Equals string `json:"equals"`
}

func (c *StringEqualsCondition) Name() string {
	return "string_equals_condition"
}

func (c *StringEqualsCondition) Meets(val interface{}, _ *Request) bool {
	s, ok := val.(string)

	return ok && s == c.Equals
}

type CIDRCondition struct {
	CIDR string `json:"cidr"`
}

func (c *CIDRCondition) Meets(value interface{}, _ *Request) bool {
	ips, ok := value.(string)
	if !ok {
		return false
	}

	_, cidrnet, err := net.ParseCIDR(c.CIDR)
	if err != nil {
		return false
	}

	ip := net.ParseIP(ips)
	if ip == nil {
		return false
	}

	return cidrnet.Contains(ip)
}

func (c *CIDRCondition) Name() string {
	return "cidr_condition"
}
