package filter

// Set is a set of strings.
type Set map[string]struct{}

// NewSet returns a string set from a slice of its members.
func NewSet(members ...string) Set {
	set := make(Set, len(members))

	for _, m := range members {
		set[m] = struct{}{}
	}

	return set
}

// Copy returns a copy of s.
func (s Set) Copy() Set {
	set := make(map[string]struct{}, len(s))

	for m := range s {
		set[m] = struct{}{}
	}

	return set
}

// IntersectInPlace updates s to the intersection of x and itself.
func (s Set) IntersectInPlace(x Set) {
	for m := range s {
		if _, ok := x[m]; !ok {
			delete(s, m)
		}
	}
}

// UnionInPlace updates s to the union of x and itself.
func (s Set) UnionInPlace(x Set) {
	for v := range x {
		s[v] = struct{}{}
	}
}
