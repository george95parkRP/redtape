package sql

import (
	"context"
	"errors"
	"fmt"

	"github.com/blushft/redtape"
	"github.com/blushft/redtape/manager/sql/ent"
	poent "github.com/blushft/redtape/manager/sql/ent/policyoptions"
	"github.com/blushft/redtape/match"
	_ "github.com/lib/pq"
)

type sqlPolicyMgr struct {
	client *ent.Client
}

// NewSqlManager returns an implementation of the PolicyManager interface
// with an ent client to make calls to the database.
func NewSqlManager(opts ...SqlManagerOption) (redtape.PolicyManager, error) {
	options := NewSqlManagerOptions(opts...)

	pm := &sqlPolicyMgr{}

	c, err := ent.Open(options.Dialect, options.ConnString)
	if err != nil {
		return nil, err
	}

	if err := c.Schema.Create(context.Background()); err != nil {
		return nil, err
	}

	pm.client = c

	return pm, nil
}

// Create creates a policy in the database.
func (pm *sqlPolicyMgr) Create(p redtape.Policy) error {
	// Let's first insert subjects/conditions in order to create the edges.
	ctx := context.Background()

	subs, err := pm.createSubjects(p.Subjects(), ctx)
	if err != nil {
		return err
	}

	// We are ready to create our policy option along with its edges.
	_, err = pm.client.PolicyOptions.Create().
		SetID(p.ID()).
		SetName(p.Name()).
		SetDescription(p.Description()).
		SetResources(p.Resources()).
		SetActions(p.Actions()).
		SetScopes(p.Scopes()).
		SetEffect(string(p.Effect())).
		AddSubjects(subs...).
		Save(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Update updates a policy given an ID.
func (pm *sqlPolicyMgr) Update(p redtape.Policy) error {
	_, err := pm.client.PolicyOptions.UpdateOneID(p.ID()).
		SetName(p.Name()).
		SetDescription(p.Description()).
		SetResources(p.Resources()).
		SetActions(p.Actions()).
		SetScopes(p.Scopes()).
		SetEffect(string(p.Effect())).
		Save(context.Background())
	if err != nil {
		return err
	}

	return nil
}

// Get gets a policy from the database given an ID.
func (pm *sqlPolicyMgr) Get(id string) (redtape.Policy, error) {
	policy, err := pm.client.PolicyOptions.Query().
		Where(poent.ID(id)).
		WithSubjects(func(q *ent.SubjectsQuery) {
			q.WithConditions()
		}).
		First(context.Background())
	if err != nil {
		return nil, err
	}

	rp, err := entPolicyToTape(policy)
	if err != nil {
		return nil, err
	}

	return rp, nil
}

// Delete will delete a policy from the database given an ID.
func (pm *sqlPolicyMgr) Delete(id string) error {
	return pm.client.PolicyOptions.DeleteOneID(id).Exec(context.Background())
}

func (pm *sqlPolicyMgr) All() ([]redtape.Policy, error) {
	policies, err := pm.client.PolicyOptions.Query().
		WithSubjects(func(q *ent.SubjectsQuery) {
			q.WithConditions()
		}).
		All(context.Background())
	if err != nil {
		return nil, err
	}

	result := []redtape.Policy{}
	for _, p := range policies {
		rp, err := entPolicyToTape(p)
		if err != nil {
			return nil, err
		}
		result = append(result, rp)
	}

	return result, nil
}

// FindByRequest will search the database for a policy that has the exact same data as the request.
func (pm *sqlPolicyMgr) FindByRequest(req *redtape.Request) ([]redtape.Policy, error) {
	if req.Resource == "" || req.Action == "" || req.Scope == "" || len(req.Subjects) < 1 {
		return nil, errors.New(fmt.Sprintf("Request had an empty field: %v", req))
	}

	policies, err := pm.client.PolicyOptions.Query().
		WithSubjects(func(q *ent.SubjectsQuery) {
			q.WithConditions()
		}).
		All(context.Background())
	if err != nil {
		return nil, err
	}

	result := []redtape.Policy{}
	for _, p := range policies {
		if found := findByResource(req.Resource, p.Resources); !found {
			continue
		}

		if found := findByScope(req.Scope, p.Scopes); !found {
			continue
		}

		if found := findByAction(req.Action, p.Actions); !found {
			continue
		}

		rp, err := entPolicyToTape(p)
		if err != nil {
			return nil, err
		}
		result = append(result, rp)
	}

	return result, nil
}

func findByAction(act string, actions []string) bool {
	for _, a := range actions {
		if match.Wildcard(a, act) {
			return true
		}
	}

	return false
}
func findByResource(res string, resources []string) bool {
	for _, r := range resources {
		if match.Wildcard(r, res) {
			return true
		}
	}

	return false
}
func findByScope(scope string, scopes []string) bool {
	for _, s := range scopes {
		if match.Wildcard(s, scope) {
			return true
		}
	}

	return false
}

// FindByRole will return a policy from the database with the same role name.
func (pm *sqlPolicyMgr) FindByRole(role string) ([]redtape.Policy, error) {
	return nil, nil
}

// FindByResource will return a policy from the database that has the resource in the resources field.
func (pm *sqlPolicyMgr) FindByResource(resource string) ([]redtape.Policy, error) {
	policies, err := pm.client.PolicyOptions.Query().
		WithSubjects(func(q *ent.SubjectsQuery) {
			q.WithConditions()
		}).
		All(context.Background())
	if err != nil {
		return nil, err
	}

	result := []redtape.Policy{}

	for _, p := range policies {
		found := false
		for _, r := range p.Resources {
			if r == resource {
				found = true
				break
			}
		}

		if found {
			rp, err := entPolicyToTape(p)
			if err != nil {
				return nil, err
			}
			result = append(result, rp)
		}
	}

	return result, nil
}

// FindByResource will return a policy from the database that has the scope in the scopes field.
func (pm *sqlPolicyMgr) FindByScope(scope string) ([]redtape.Policy, error) {
	policies, err := pm.client.PolicyOptions.Query().
		WithSubjects(func(q *ent.SubjectsQuery) {
			q.WithConditions()
		}).
		All(context.Background())
	if err != nil {
		return nil, err
	}

	result := []redtape.Policy{}

	for _, p := range policies {
		found := false
		for _, s := range p.Scopes {
			if s == scope {
				found = true
				break
			}
		}

		if found {
			rp, err := entPolicyToTape(p)
			if err != nil {
				return nil, err
			}
			result = append(result, rp)
		}
	}

	return result, nil
}

// Creates Conditions and Roles in the database and returns their ent reference.
func (pm *sqlPolicyMgr) createSubjects(subjects []redtape.Subject, ctx context.Context) ([]*ent.Subjects, error) {
	entSubs := []*ent.Subjects{}

	for _, sub := range subjects {
		entConds := []*ent.Conditions{}

		for _, cond := range sub.Conditions() {
			typ, val := getTypeAndVal(cond)
			opts := map[string]interface{}{
				"value": val,
			}

			c, err := pm.client.Conditions.Create().
				SetName(cond.Name()).
				SetType(typ).
				SetOptions(opts).
				Save(ctx)
			if err != nil {
				return nil, err
			}

			entConds = append(entConds, c)
		}

		s, err := pm.client.Subjects.Create().
			SetType(sub.Type()).
			AddConditions(entConds...).
			Save(ctx)
		if err != nil {
			return nil, err
		}

		entSubs = append(entSubs, s)
	}

	return entSubs, nil
}

// Translate ent's Policy to redtape's Policy.
func entPolicyToTape(p *ent.PolicyOptions) (redtape.Policy, error) {
	po := redtape.PolicyOptions{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Resources:   p.Resources,
		Actions:     p.Actions,
		Scopes:      p.Scopes,
		Effect:      p.Effect,
	}

	rtSubjects := []redtape.Subject{}
	for _, s := range p.Edges.Subjects {
		rs, err := entSubjectToTape(s)
		if err != nil {
			return nil, err
		}
		rtSubjects = append(rtSubjects, rs)
	}

	po.Subjects = rtSubjects

	return redtape.MustNewPolicy(redtape.SetPolicyOptions(po)), nil
}

// Translate ent's Role to redtape's Role.
func entSubjectToTape(subject *ent.Subjects) (redtape.Subject, error) {
	s, err := redtape.NewSubject(subject.Type, redtape.WithConditions(entCondsToTape(subject.Edges.Conditions)...))
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Translate ent's Conditions to redtape's Conditions.
func entCondsToTape(conds []*ent.Conditions) []redtape.ConditionOptions {
	res := []redtape.ConditionOptions{}

	for _, c := range conds {
		res = append(res, redtape.ConditionOptions{
			Name:    c.Name,
			Type:    c.Type,
			Options: c.Options,
		})
	}

	return res
}

// Returns condition's type and its value.
func getTypeAndVal(val interface{}) (string, bool) {
	// TODO: RoleEqualsCondition doesn't have a value right now.
	switch v := val.(type) {
	case *redtape.BoolCondition:
		return v.Name(), v.Value
	default:
		return "", false
	}
}
