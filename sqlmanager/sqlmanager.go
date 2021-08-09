package sqlmanager

import (
	"context"
	"errors"
	"fmt"

	"github.com/blushft/redtape"
	"github.com/blushft/redtape/sqlmanager/ent"
	cent "github.com/blushft/redtape/sqlmanager/ent/conditions"
	poent "github.com/blushft/redtape/sqlmanager/ent/policyoptions"
	rent "github.com/blushft/redtape/sqlmanager/ent/roles"
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
	// Let's first insert roles/conditions in order to create the edges.
	ctx := context.Background()

	roles, conditions, err := pm.createConditionsRoles(p.Roles(), p.Conditions(), ctx)
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
		AddRoles(roles...).
		AddConditions(conditions...).
		Save(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Update updates a policy given an ID.
func (pm *sqlPolicyMgr) Update(p redtape.Policy) error {
	// Delete current conditions and roles associated with this policy first.
	ctx := context.Background()

	policies, err := pm.client.PolicyOptions.Query().
		WithConditions().
		WithRoles().
		Where(poent.ID(p.ID())).
		First(ctx)
	if err != nil {
		return err
	}

	for _, c := range policies.Edges.Conditions {
		_, err := pm.client.Conditions.Delete().Where(cent.Name(c.Name)).Exec(ctx)
		if err != nil {
			return err
		}
	}

	for _, r := range policies.Edges.Roles {
		_, err := pm.client.Roles.Delete().Where(rent.Name(r.Name)).Exec(ctx)
		if err != nil {
			return err
		}
	}

	// Create the new conditions and roles to associate them back to this policy.
	roles, conditions, err := pm.createConditionsRoles(p.Roles(), p.Conditions(), ctx)
	if err != nil {
		return err
	}

	// Update policy.
	_, err = pm.client.PolicyOptions.UpdateOneID(p.ID()).
		SetName(p.Name()).
		SetDescription(p.Description()).
		SetResources(p.Resources()).
		SetActions(p.Actions()).
		SetScopes(p.Scopes()).
		SetEffect(string(p.Effect())).
		AddConditions(conditions...).
		AddRoles(roles...).
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
		WithConditions().
		WithRoles().
		First(context.Background())
	if err != nil {
		return nil, err
	}

	return entPolicyToTape(policy), nil
}

// Delete will delete a policy from the database given an ID.
func (pm *sqlPolicyMgr) Delete(id string) error {
	return pm.client.PolicyOptions.DeleteOneID(id).Exec(context.Background())
}

func (pm *sqlPolicyMgr) All(limit, offset int) ([]redtape.Policy, error) {
	policies, err := pm.client.PolicyOptions.Query().
		WithConditions().
		WithRoles().
		Limit(limit).
		Offset(offset).
		All(context.Background())
	if err != nil {
		return nil, err
	}

	result := []redtape.Policy{}
	for _, p := range policies {
		result = append(result, entPolicyToTape(p))
	}

	return result, nil
}

// FindByRequest will search the database for a policy that has the exact same data as the request.
func (pm *sqlPolicyMgr) FindByRequest(req *redtape.Request) ([]redtape.Policy, error) {
	if req.Resource == "" || req.Action == "" || req.Scope == "" || req.Role == "" {
		return nil, errors.New(fmt.Sprintf("Request had an empty field: %v", req))
	}

	policies, err := pm.client.PolicyOptions.Query().
		WithConditions().
		WithRoles().
		All(context.Background())
	if err != nil {
		return nil, err
	}

	result := []redtape.Policy{}

	// Traverse each policy.
	for _, p := range policies {
		found := false

		// Check if any of the resources are the same.
		for _, r := range p.Resources {
			if r == req.Resource {
				found = true
				break
			}
		}

		// If any time the req's data isn't found, we can just not return this policy.
		if found {
			for i, a := range p.Actions {
				if a == req.Action {
					break
				} else if i == len(p.Actions)-1 {
					found = false
				}
			}
		}

		if found {
			for i, s := range p.Scopes {
				if s == req.Scope {
					break
				} else if i == len(p.Scopes)-1 {
					found = false
				}
			}
		}

		if found {
			for i, role := range p.Edges.Roles {
				if role.ID == req.Role {
					break
				} else if i == len(p.Edges.Roles)-1 {
					found = false
				}
			}
		}

		if found {
			result = append(result, entPolicyToTape(p))
		}
	}

	return result, nil
}

// FindByRole will return a policy from the database with the same role name.
func (pm *sqlPolicyMgr) FindByRole(role string) ([]redtape.Policy, error) {
	policies, err := pm.client.PolicyOptions.Query().
		WithConditions().
		WithRoles(func(q *ent.RolesQuery) {
			q.Where(rent.Name(role))
		}).
		All(context.Background())
	if err != nil {
		return nil, err
	}

	result := []redtape.Policy{}

	for _, p := range policies {
		result = append(result, entPolicyToTape(p))
	}

	return result, nil
}

// FindByResource will return a policy from the database that has the resource in the resources field.
func (pm *sqlPolicyMgr) FindByResource(resource string) ([]redtape.Policy, error) {
	policies, err := pm.client.PolicyOptions.Query().
		WithConditions().
		WithRoles().
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
			result = append(result, entPolicyToTape(p))
		}
	}

	return result, nil
}

// FindByResource will return a policy from the database that has the scope in the scopes field.
func (pm *sqlPolicyMgr) FindByScope(scope string) ([]redtape.Policy, error) {
	policies, err := pm.client.PolicyOptions.Query().
		WithConditions().
		WithRoles().
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
			result = append(result, entPolicyToTape(p))
		}
	}

	return result, nil
}

// Creates Conditions and Roles in the database and returns their ent reference.
func (pm *sqlPolicyMgr) createConditionsRoles(roles []*redtape.Role, conditions redtape.Conditions,
	ctx context.Context) ([]*ent.Roles, []*ent.Conditions, error) {
	entRoles := []*ent.Roles{}
	entConds := []*ent.Conditions{}

	for _, role := range roles {
		r, err := pm.client.Roles.Create().
			SetID(role.ID).
			SetName(role.Name).
			SetDescription(role.Description).
			Save(ctx)
		if err != nil {
			return nil, nil, err
		}

		entRoles = append(entRoles, r)
	}

	for name, cond := range conditions {
		typ, val := getTypeAndVal(cond)
		opts := map[string]interface{}{
			"value": val,
		}

		c, err := pm.client.Conditions.Create().
			SetName(name).
			SetType(typ).
			SetOptions(opts).
			Save(ctx)
		if err != nil {
			return nil, nil, err
		}

		entConds = append(entConds, c)
	}

	return entRoles, entConds, nil
}

// Translate ent's Policy to redtape's Policy.
func entPolicyToTape(p *ent.PolicyOptions) redtape.Policy {
	po := redtape.PolicyOptions{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Resources:   p.Resources,
		Actions:     p.Actions,
		Scopes:      p.Scopes,
		Effect:      p.Effect,
	}

	rtRoles := []*redtape.Role{}
	for _, r := range p.Edges.Roles {
		rtRoles = append(rtRoles, entRoleToTape(r))
	}

	rtConds := []redtape.ConditionOptions{}
	for _, c := range p.Edges.Conditions {
		rtConds = append(rtConds, entCondToTape(c))
	}

	po.Roles = rtRoles
	po.Conditions = rtConds

	return redtape.MustNewPolicy(redtape.SetPolicyOptions(po))
}

// Translate ent's Role to redtape's Role.
func entRoleToTape(role *ent.Roles) *redtape.Role {
	return &redtape.Role{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
	}
}

// Translate ent's Conditions to redtape's Conditions.
func entCondToTape(cond *ent.Conditions) redtape.ConditionOptions {
	return redtape.ConditionOptions{
		Name:    cond.Name,
		Type:    cond.Type,
		Options: cond.Options,
	}
}

// Returns condition's type and its value.
func getTypeAndVal(val interface{}) (string, bool) {
	// TODO: RoleEqualsCondition doesn't have a value right now.
	switch v := val.(type) {
	case *redtape.BoolCondition:
		return v.Name(), v.Value
	case *redtape.RoleEqualsCondition:
		return v.Name(), false
	default:
		return "", false
	}
}
