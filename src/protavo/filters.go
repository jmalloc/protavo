package protavo

import "github.com/jmalloc/protavo/src/protavo/filter"

// IsOneOf matches documents that have one of the given IDs.
func IsOneOf(ids ...string) filter.Condition {
	return &filter.IsOneOf{
		Values: filter.NewSet(ids...),
	}
}

// HasUniqueKeyIn matches documents that have one of the given unique keys.
func HasUniqueKeyIn(keys ...string) filter.Condition {
	return &filter.HasUniqueKeyIn{
		Values: filter.NewSet(keys...),
	}
}

// HasKeys matches documents that have all of the given keys.
func HasKeys(keys ...string) filter.Condition {
	return &filter.HasKeys{
		Values: filter.NewSet(keys...),
	}
}
