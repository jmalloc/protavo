package protavo

import "github.com/jmalloc/protavo/src/protavo/filter"

// IsOneOf matches all of the documents with the given IDs.
func IsOneOf(ids ...string) filter.Condition {
	return &filter.IsOneOf{
		Values: filter.NewSet(ids...),
	}
}

// HasUniqueKeyIn matches all of the documents that have one of the given unique
// keys.
func HasUniqueKeyIn(keys ...string) filter.Condition {
	return &filter.HasUniqueKeyIn{
		Values: filter.NewSet(keys...),
	}
}

// HasKeys matches documents that have all of the given keys, regardless of key
// type.
//
// Note that the parameters to this condition form a logical AND. That is, a
// document is required to have ALL of the keys in order to match.
func HasKeys(keys ...string) filter.Condition {
	return &filter.HasKeys{
		Values: filter.NewSet(keys...),
	}
}

// TODO(jmalloc): implement HasKeyIn() and HasSharedKeyIn()
// TODO(jmalloc): find some way to implement a logical OR of key sets, something
// conceptually like HasKeys(set1, set2, ...)
