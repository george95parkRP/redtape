package redtape

import (
	"context"
	"testing"

	"github.com/blushft/redtape/role"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type RedtapeSuite struct {
	suite.Suite
}

func TestRedtapeSuite(t *testing.T) {
	suite.Run(t, new(RedtapeSuite))
}

func (s *RedtapeSuite) TestBPolicies() {
	table := []struct {
		opts PolicyOptions
	}{
		{
			opts: NewPolicyOptions(
				PolicyID(uuid.NewString()),
				PolicyName("test_policy"),
				PolicyDescription("just a test"),
				SetActions("create", "delete", "update", "read"),
				SetResources("database"),
				PolicyAllow(),
				WithCondition(ConditionOptions{
					Name: "test_cond",
					Type: "bool",
					Options: map[string]interface{}{
						"value": true,
					},
				}),
				WithSubject(NewSubject("allow_test")),
			),
		},
	}

	man := NewPolicyManager()

	for _, tt := range table {
		p, err := NewPolicy(SetPolicyOptions(tt.opts))
		s.Require().NoError(err)

		err = man.Create(p)
		s.Require().NoError(err)
	}
}

func (s *RedtapeSuite) TestCEnforce() {
	m := NewMatcher()
	pm := NewPolicyManager()

	allow := role.New("test.A")
	deny := role.New("test.B")

	subA := NewSubject(
		uuid.NewString(),
		SubjectRole(allow),
		SubjectEnforcer(func(r *Request) error {
			return nil
		}),
	)

	subB := NewSubject(
		uuid.NewString(),
		SubjectRole(deny),
	)

	popts := []PolicyOptions{
		{
			ID:          uuid.NewString(),
			Name:        "test_policy_allow",
			Description: "testing",
			Subjects: []Subject{
				subA,
			},
			Resources: []string{
				"test_resource",
			},
			Actions: []string{
				"test",
			},
			Scopes: []string{
				"test_scope",
			},
			Effect: "allow",
			Conditions: []ConditionOptions{
				{
					Name: "match_me",
					Type: "bool",
					Options: map[string]interface{}{
						"value": true,
					},
				},
			},
		},
		{
			ID:          uuid.NewString(),
			Name:        "test_policy",
			Description: "testing",
			Subjects: []Subject{
				subB,
			},
			Resources: []string{
				"test_resource",
			},
			Actions: []string{
				"test",
			},
			Scopes: []string{
				"test_scope",
			},
			Effect: "deny",
			Conditions: []ConditionOptions{
				{
					Name: "match_me",
					Type: "bool",
					Options: map[string]interface{}{
						"value": true,
					},
				},
			},
		},
	}

	for _, po := range popts {
		err := pm.Create(MustNewPolicy(SetPolicyOptions(po)))
		s.Require().NoError(err)
	}

	e := NewEnforcer(pm, m, nil)

	req := &Request{
		Resource: "test_resource",
		Action:   "test",
		Scope:    "test_scope",
		Subject:  subA,
		Context: NewRequestContext(context.TODO(), map[string]interface{}{
			"match_me": true,
		}),
	}

	err := e.Enforce(req)
	s.Require().NoError(err, "should be allowed")

	req.Subject = subB

	err = e.Enforce(req)
	s.Require().Error(err, "should be denied")
}
