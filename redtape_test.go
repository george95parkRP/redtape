package redtape

import (
	"context"
	"testing"

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
	sub, err := NewSubject(
		"allow_test",
		WithConditions(ConditionOptions{
			Name: "test_cond",
			Type: "bool",
			Options: map[string]interface{}{
				"value": true,
			},
		}),
	)
	s.Require().NoError(err)

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
				WithSubjects(sub),
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

	subA, err := NewSubject("subA")
	s.Require().NoError(err)

	subB, err := NewSubject("subB")
	s.Require().NoError(err)

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
		Subjects: []string{"subA"},
		Context: NewRequestContext(context.TODO(), map[string]interface{}{
			"match_me": true,
		}),
	}

	err = e.Enforce(req)
	s.Require().NoError(err, "should be allowed")

	req.Subjects = []string{"subB"}

	err = e.Enforce(req)
	s.Require().Error(err, "should be denied")
}
