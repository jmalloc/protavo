package filter

// Visitor an interface for performing operations with filter conditions.
type Visitor interface {
	MatchDocumentID(*MatchDocumentID) error
	MatchAllKeys(*MatchAllKeys) error
	MatchUniqueKey(*MatchUniqueKey) error
	MatchNothing(*MatchNothing) error
}
