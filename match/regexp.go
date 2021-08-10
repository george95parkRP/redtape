package match

import (
	"regexp"
	"strings"
)

type Regexp struct {
	start rune
	stop  rune
	cache map[string]*regexp.Regexp
}

func NewRegexp() *Regexp {
	return &Regexp{
		start: '<',
		stop:  '>',
		cache: make(map[string]*regexp.Regexp),
	}
}

func (r *Regexp) Match(s, val string) (bool, error) {
	if !strings.ContainsRune(s, r.start) {
		if Wildcard(s, val) {
			return true, nil
		}

		return false, nil
	}

	var reg *regexp.Regexp
	var err error

	reg, ok := r.cache[s]
	if !ok {
		reg, err = CompileDelimitedRegex(val, r.start, r.stop)
		if err != nil {
			return false, err
		}

		r.cache[s] = reg
	}

	if reg.MatchString(val) {
		return true, nil
	}

	return false, nil
}
