package filter

// Optimize performs basic optimization of the given filter.
func Optimize(f *Filter) *Filter {
	o := &optimizer{}

	possible, err := f.Accept(o)
	if err != nil {
		panic(err)
	}

	// if the filter presents an impossible combination, return an empty filter
	if !possible {
		return &Filter{}
	}

	var conds []Condition

	if o.isOneOf != nil {
		conds = append(conds, o.isOneOf)
	}

	if o.hasUniqueKeyIn != nil {
		conds = append(conds, o.hasUniqueKeyIn)
	}

	if o.hasKeys != nil {
		conds = append(conds, o.hasKeys)
	}

	// TODO(jmalloc): we could scan o.hasKeys to look for any keys that are also in
	// o.hasUniqueKeyIn and remove it from the set.

	// if all of the conditions have been optimized away return a nil filter
	if len(conds) == 0 {
		return nil
	}

	return &Filter{conds}
}

type optimizer struct {
	isOneOf             *IsOneOf
	isOneOfCount        int
	hasUniqueKeyIn      *HasUniqueKeyIn
	hasUniqueKeyInCount int
	hasKeys             *HasKeys
	hasKeysCount        int
}

func (o *optimizer) IsOneOf(c *IsOneOf) (bool, error) {
	// if there are no IDs, there can be no matches
	if len(c.Values) == 0 {
		return false, nil
	}

	o.isOneOfCount++

	// if this is the first condition of this type we've seen, use it as is.
	if o.isOneOfCount == 1 {
		o.isOneOf = c
		return true, nil
	}

	// if this is the second condition we've seen, perform a copy-on-write
	if o.isOneOfCount == 2 {
		o.isOneOf = &IsOneOf{
			Values: o.isOneOf.Values.Copy(),
		}
	}

	// compute the intersection of the values from the existing condition and this
	// new one. bail early if the intersection is empty.
	o.isOneOf.Values.IntersectInPlace(c.Values)

	return len(o.isOneOf.Values) > 0, nil
}

func (o *optimizer) HasUniqueKeyIn(c *HasUniqueKeyIn) (bool, error) {
	// if there are no IDs, there can be no matches
	if len(c.Values) == 0 {
		return false, nil
	}

	o.hasUniqueKeyInCount++

	// if this is the first condition of this type we've seen, use it as is.
	if o.hasUniqueKeyInCount == 1 {
		o.hasUniqueKeyIn = c
		return true, nil
	}

	// if this is the second condition we've seen, perform a copy-on-write
	if o.hasUniqueKeyInCount == 2 {
		o.hasUniqueKeyIn = &HasUniqueKeyIn{
			Values: o.hasUniqueKeyIn.Values.Copy(),
		}
	}

	// compute the intersection of the values from the existing condition and this
	// new one. bail early if the intersection is empty.
	o.hasUniqueKeyIn.Values.IntersectInPlace(c.Values)

	return len(o.hasUniqueKeyIn.Values) > 0, nil
}

func (o *optimizer) HasKeys(c *HasKeys) (bool, error) {
	if len(c.Values) == 0 {
		return true, nil
	}

	o.hasKeysCount++

	// if this is the first condition of this type we've seen, use it as is.
	if o.hasKeysCount == 1 {
		o.hasKeys = c
		return true, nil
	}

	// if this is the second condition we've seen, perform a copy-on-write
	if o.hasKeysCount == 2 {
		o.hasKeys = &HasKeys{
			Values: o.hasKeys.Values.Copy(),
		}
	}

	// compute the union of the values from the existing condition and this one.
	o.hasKeys.Values.UnionInPlace(c.Values)

	return true, nil
}
