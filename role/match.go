package role

import (
	"github.com/blushft/redtape/match"
)

func Match(r *Role, val string) (bool, error) {
	m := &simpleMatcher{}
	return m.Match(r, val)
}

type Matcher interface {
	Match(r *Role, val string) (bool, error)
}

type simpleMatcher struct{}

// NewMatcher returns the default Matcher implementation.
func NewMatcher() Matcher {
	return &simpleMatcher{}
}

// Match evaluates true when the provided val wildcard matches at least one role in Role#EffectiveRoles.
func (m *simpleMatcher) Match(r *Role, val string) (bool, error) {
	er := r.EffectiveRoles()

	for _, rr := range er {
		if match.Wildcard(val, rr.ID) {
			return true, nil
		}
	}

	return false, nil
}

type regexMatcher struct {
	r *match.Regexp
}

// NewRegexMatcher returns a Matcher using delimited regex for matching.
func NewRegexMatcher() Matcher {
	return &regexMatcher{
		r: match.NewRegexp(),
	}
}

// MatchRole evaluates true when the provided val regex matches at least one role in Role#EffectiveRoles.
func (m *regexMatcher) Match(r *Role, val string) (bool, error) {
	ef := r.EffectiveRoles()

	def := make([]string, 0, len(ef))
	for _, rr := range ef {
		def = append(def, rr.ID)
	}

	return m.match(def, val)
}

func (m *regexMatcher) match(def []string, val string) (bool, error) {
	for _, h := range def {
		ok, err := m.r.Match(h, val)
		if err != nil {
			return false, err
		}

		if ok {
			return true, nil
		}
	}

	return false, nil
}
