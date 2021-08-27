package redtape

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type EngineSuite struct {
	suite.Suite
}

func TestEngineSuite(t *testing.T) {
	suite.Run(t, new(EngineSuite))
}

func (s *EngineSuite) TestAEngine() {
	engine := NewEngine(NewPolicyManager())

	sub, err := NewSubject(
		"rpx.example.api.hello",
		WithConditions([]ConditionOptions{
			{
				Name: "service",
				Type: "string_equals_condition",
				Options: map[string]interface{}{
					"equals": "rpx.example.rpc.hello",
				},
			},
			{
				Name: "match_me",
				Type: "bool",
				Options: map[string]interface{}{
					"value": true,
				},
			},
		}...),
	)
	s.Require().NoError(err)

	po := NewPolicyOptions(
		PolicyID(uuid.NewString()),
		PolicyName("test_policy"),
		PolicyDescription("just a test"),
		PolicyAllow(),
		SetResources("test_resource"),
		SetActions("create", "delete"),
		SetScopes("test_scope"),
		WithSubjects(sub),
	)

	err = engine.Grant(MustNewPolicy(SetPolicyOptions(po)))
	s.Require().NoError(err)

	req := &Request{
		Resource: "test_resource",
		Action:   "create",
		Scope:    "test_scope",
		Subjects: []string{"rpx.example.api.hello"},
		Context: NewRequestContext(context.TODO(), map[string]interface{}{
			"service":  "rpx.example.rpc.hello",
			"match_me": true,
		}),
	}

	err = engine.Verify(req)
	s.Require().NoError(err, "should be allowed")

	fakeReq := &Request{
		Resource: "test_resource",
		Action:   "create",
		Scope:    "test_scope",
		Subjects: []string{"fake_test_subject"},
		Context: NewRequestContext(context.TODO(), map[string]interface{}{
			"match_me": true,
		}),
	}

	err = engine.Verify(fakeReq)
	s.Require().Error(err, "should not be allowed")
}
