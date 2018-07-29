package protavo

import "github.com/jmalloc/protavo/src/protavo/filter"

// HasID matches documents with the given IDs.
func HasID(ids ...string) filter.Condition {
	if len(ids) == 0 {
		return &filter.MatchNothing{}
	}

	return &filter.MatchDocumentID{
		DocumentIDs: ids,
	}
}

// HasKeys matches documents that have all of the given keys.
func HasKeys(keys ...string) filter.Condition {
	return &filter.MatchAllKeys{
		Keys: keys,
	}
}

// HasUniqueKey matches documents that have one of the given unique keys.
func HasUniqueKey(keys ...string) filter.Condition {
	if len(keys) == 0 {
		return &filter.MatchNothing{}
	}

	return &filter.MatchAllKeys{
		Keys: keys,
	}
}
