package drivertest

import (
	"github.com/jmalloc/protavo/src/protavo"
	g "github.com/onsi/ginkgo"
)

// Describe defines the standard test suite for an implementation of
// protavo.Driver.
func Describe(
	name string,
	before func() (*protavo.DB, error),
	after func(),
) {
	describe(name, before, after)

	var db *protavo.DB

	describe(
		name+" (sub-namespace)",
		func() (*protavo.DB, error) {
			var err error
			db, err = before()
			if err != nil {
				return nil, err
			}

			return db.Namespace("sub-namespace"), nil
		},
		func() {
			if db != nil {
				db.Close()
			}

			if after != nil {
				after()
			}
		},
	)
}

func describe(
	name string,
	before func() (*protavo.DB, error),
	after func(),
) {
	g.Describe(name, func() {
		describeFetchAll(before, after)
		describeFetchWhere(before, after)
		describeSave(before, after)
		describeForceSave(before, after)
		describeDelete(before, after)
		describeDeleteWhere(before, after)
		describeDeleteNamespace(before, after)

		describeFilters(before, after)
	})
}
