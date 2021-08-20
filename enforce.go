package redtape

import (
	"errors"
)

type EnforceFunc func(*Request) error

// Enforcer interface provides methods to enforce policies against a request.
type Enforcer interface {
	Enforce(*Request) error
}

type enforcer struct {
	manager PolicyManager
	matcher Matcher
	auditor Auditor
}

// NewEnforcer returns a default Enforcer combining a PolicyManager, Matcher, and Auditor.
func NewEnforcer(manager PolicyManager, matcher Matcher, auditor Auditor) Enforcer {
	return &enforcer{
		manager: manager,
		matcher: matcher,
		auditor: auditor,
	}
}

func NewDefaultEnforcer(manager PolicyManager) Enforcer {
	return NewEnforcer(manager, DefaultMatcher, NewConsoleAuditor(AuditAll))
}

// Enforce fulfills the Enforce method of Enforcer. The default implementation matches the Request against
// the range of stored Policies and evaluating each.
// Polices are matched first by Action, then Role, Resource, Scope and finally Condition. If a match is found, the
// configured Policy Effect is applied.
// TODO: return explicit PolicyEffect and use error to indicate processing failures
func (e *enforcer) Enforce(r *Request) error {
	allow := false
	matched := []Policy{}

	e.auditReq(r)

	pols, err := e.manager.FindByRequest(r)
	if err != nil {
		return err
	}

	for _, p := range pols {
		match, err := e.evalPolicy(r, p)
		if err != nil {
			return err
		}

		if !match {
			continue
		}

		matched = append(matched, p)

		// deny overrides all
		if p.Effect() == PolicyEffectDeny {
			e.auditEffect(r, PolicyEffectDeny)
			return NewErrRequestDeniedExplicit(p)
		}

		allow = true
	}

	if !allow && DefaultPolicyEffect == PolicyEffectDeny {
		e.auditEffect(r, PolicyEffectDeny)
		return NewErrRequestDeniedImplicit(errors.New("access denied because no policy allowed access"))
	}

	e.auditEffect(r, PolicyEffectAllow)

	return nil
}

func (e *enforcer) evalPolicy(r *Request, p Policy) (bool, error) {
	// match actions
	am, err := e.matcher.MatchPolicy(p, p.Actions(), r.Action)
	if err != nil {
		return false, err
	}

	if !am {
		return false, nil
	}

	// match resources
	resm, err := e.matcher.MatchPolicy(p, p.Resources(), r.Resource)
	if err != nil {
		return false, err
	}
	if !resm {
		return false, nil
	}

	// match scopes
	scm, err := e.matcher.MatchPolicy(p, p.Scopes(), r.Scope)
	if err != nil {
		return false, err
	}
	if !scm {
		return false, nil
	}

	// find the subject
	// check conditions
	// TODO: could there be more than one subject with the same type ?
	condMeets := false
	for _, sub := range p.Subjects() {
		if sub.Type() == r.Subject.Type() {
			if sub.Conditions().Meets(r) {
				condMeets = true
				break
			}
		}
	}

	if !condMeets {
		return false, nil
	}

	return true, nil
}

func (e *enforcer) auditReq(req *Request) {
	if e.auditor != nil {
		e.auditor.LogRequest(req)
	}
}

func (e *enforcer) auditEffect(req *Request, effect PolicyEffect) {
	if e.auditor != nil {
		e.auditor.LogPolicyEffect(req, effect)
	}
}
