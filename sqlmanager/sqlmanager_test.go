package sqlmanager

import (
	"reflect"
	"testing"

	"github.com/blushft/redtape"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type SqlManagerSuite struct {
	suite.Suite
}

func TestSqlManagerSuite(t *testing.T) {
	suite.Run(t, new(SqlManagerSuite))
}

func (s *SqlManagerSuite) TestAPolicyOptions() {
	id := uuid.NewString()

	opts := redtape.NewPolicyOptions(
		redtape.PolicyID(id),
		redtape.PolicyName("test_policy"),
		redtape.PolicyDescription("just a test"),
		redtape.SetActions("create", "delete", "update", "read"),
		redtape.SetResources("database"),
		redtape.PolicyAllow(),
		redtape.WithCondition(redtape.ConditionOptions{
			Name: "test_cond",
			Type: "bool",
			Options: map[string]interface{}{
				"value": true,
			},
		}),
		redtape.WithRole(redtape.NewRole("allow_test")),
	)

	// get sql manager
	man, err := NewSqlManager(
		SetDialect("postgres"),
		SetConnString("host=localhost port=5432 user=admin dbname=policy password=password sslmode=disable"),
	)
	s.Require().NoError(err)

	// new policy
	policy := redtape.MustNewPolicy(redtape.SetPolicyOptions(opts))

	// create new policy
	s.Require().NoError(man.Create(policy))

	// get created policy
	getPolicy, err := man.Get(opts.ID)
	s.Require().NoError(err)

	s.Require().True(reflect.DeepEqual(policy, getPolicy))

	// update policy
	uptOpts := redtape.NewPolicyOptions(
		redtape.PolicyID(id),
		redtape.PolicyName("updated policy name"),
		redtape.PolicyDescription("updated description"),
		redtape.SetActions("update"),
		redtape.SetResources("database", "another resource"),
		redtape.PolicyDeny(),
		redtape.WithCondition(redtape.ConditionOptions{
			Name: "updated condition name",
			Type: "bool",
			Options: map[string]interface{}{
				"value": false,
			},
		}),
		redtape.WithRole(redtape.NewRole("updated role")),
	)

	updatedPolicy := redtape.MustNewPolicy(redtape.SetPolicyOptions(uptOpts))

	s.Require().NoError(man.Update(updatedPolicy))

	uptPolicy, err := man.Get(opts.ID)
	s.Require().NoError(err)

	s.Require().True(reflect.DeepEqual(updatedPolicy, uptPolicy))

	// delete policy
	delOpts := redtape.NewPolicyOptions(
		redtape.PolicyID(uuid.NewString()),
		redtape.PolicyName("test_policy"),
		redtape.PolicyDescription("just a test"),
		redtape.SetActions("create", "delete", "update", "read"),
		redtape.SetResources("database"),
		redtape.PolicyAllow(),
		redtape.WithCondition(redtape.ConditionOptions{
			Name: "test_cond",
			Type: "bool",
			Options: map[string]interface{}{
				"value": true,
			},
		}),
		redtape.WithRole(redtape.NewRole("allow_test")),
	)

	delPolicy := redtape.MustNewPolicy(redtape.SetPolicyOptions(delOpts))

	s.Require().NoError(man.Create(delPolicy))
	s.Require().NoError(man.Delete(delPolicy.ID()))

	_, err = man.Get(delPolicy.ID())
	s.Require().Error(err)

	// all policy
	policies, err := man.All(10, 0)
	s.Require().NoError(err)
	s.Require().Greater(len(policies), 0, "should only be at least one policy")
}

func (s *SqlManagerSuite) TestBPolicyOptions() {
	id := uuid.NewString()

	opts := redtape.NewPolicyOptions(
		redtape.PolicyID(id),
		redtape.PolicyName("Test Policy"),
		redtape.PolicyDescription("Test Description"),
		redtape.SetActions("Test Action"),
		redtape.SetResources("Test Resource"),
		redtape.SetScopes("Test Scope"),
		redtape.PolicyAllow(),
		redtape.WithCondition(redtape.ConditionOptions{
			Name: "test_cond",
			Type: "bool",
			Options: map[string]interface{}{
				"value": true,
			},
		}),
		redtape.WithRole(redtape.NewRole("Test Role")),
	)

	// get sql manager
	man, err := NewSqlManager(
		SetDialect("postgres"),
		SetConnString("host=localhost port=5432 user=admin dbname=policy password=password sslmode=disable"),
	)
	s.Require().NoError(err)

	// new policy
	policy := redtape.MustNewPolicy(redtape.SetPolicyOptions(opts))

	// create new policy
	s.Require().NoError(man.Create(policy))

	// get created policy
	getPolicy, err := man.Get(opts.ID)
	s.Require().NoError(err)

	s.Require().True(reflect.DeepEqual(policy, getPolicy))

	// find by request
	req := redtape.NewRequest(opts.Resources[0], opts.Actions[0], opts.Roles[0].ID, opts.Scopes[0])

	policies, err := man.FindByRequest(req)
	s.Require().NoError(err)

	s.Require().GreaterOrEqual(len(policies), 0, "should be at least one policy")

	s.Require().True(reflect.DeepEqual(policy, policies[0]))

	// find by empty request
	policies, err = man.FindByRequest(&redtape.Request{})
	s.Require().Error(err)

	s.Require().Equal(len(policies), 0, "should not have found a policy")

	// find by random request
	policies, err = man.FindByRequest(redtape.NewRequest(opts.Resources[0], "random", "random", "random"))
	s.Require().NoError(err)

	s.Require().Equal(len(policies), 0, "should not have found a policy 2")

	policies, err = man.FindByRequest(redtape.NewRequest("random", opts.Actions[0], "random", "random"))
	s.Require().NoError(err)

	s.Require().Equal(len(policies), 0, "should not have found a policy 3")

	// find by resource
	policies, err = man.FindByResource(req.Resource)
	s.Require().NoError(err)

	s.Require().True(reflect.DeepEqual(policy, policies[0]))

	// find by scope
	policies, err = man.FindByScope(req.Scope)
	s.Require().NoError(err)

	s.Require().True(reflect.DeepEqual(policy, policies[0]))
}
