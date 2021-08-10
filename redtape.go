package redtape

var (
	// DefaultMatcher is a simple matcher.
	DefaultMatcher = NewMatcher()
	// DefaultPolicyEffect is the policy effect to apply when no other matches can be found.
	DefaultPolicyEffect = PolicyEffectDeny
)

// MatchPolicy is a utility function that uses DefaultMatcher to evaluate whether p can be matched by val.
func MatchPolicy(p Policy, def []string, val string) (bool, error) {
	return DefaultMatcher.MatchPolicy(p, def, val)
}
