package redtape

import (
	"github.com/blushft/redtape/match"
)

// Matcher provides methods to facilitate matching policies to different request elements.
type Matcher interface {
	MatchPolicy(p Policy, def []string, val string) (bool, error)
}

type simpleMatcher struct{}

// NewMatcher returns the default Matcher implementation.
func NewMatcher() Matcher {
	return &simpleMatcher{}
}

// MatchPolicy evaluates true when the provided val wildcard matches at least one element in def.
// If def is nil, a match is assumed against any value.
func (m *simpleMatcher) MatchPolicy(p Policy, def []string, val string) (bool, error) {
	if def == nil {
		return true, nil
	}

	for _, h := range def {
		if match.Wildcard(h, val) {
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

// MatchPolicy evaluates true when the provided val regex matches at least one element in def.
func (m *regexMatcher) MatchPolicy(p Policy, def []string, val string) (bool, error) {
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
