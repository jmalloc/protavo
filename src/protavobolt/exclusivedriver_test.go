package protavobolt_test

import (
	"github.com/jmalloc/protavo/src/protavo/driver"
	"github.com/jmalloc/protavo/src/protavo/driver/drivertest"
	. "github.com/jmalloc/protavo/src/protavobolt"
)

func init() {
	drivertest.Describe(
		"ExclusiveDriver",
		func() (driver.Driver, error) {
			return OpenTemp(0600, nil)
		},
		nil,
	)
}
