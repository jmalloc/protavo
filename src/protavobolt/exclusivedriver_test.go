package protavobolt_test

import (
	"github.com/jmalloc/protavo/src/protavo"
	"github.com/jmalloc/protavo/src/protavo/driver/drivertest"
	. "github.com/jmalloc/protavo/src/protavobolt"
)

func init() {
	drivertest.Describe(
		"protavobolt.ExclusiveDriver",
		func() (*protavo.DB, error) {
			return OpenTemp(0600, nil)
		},
		nil,
	)
}
