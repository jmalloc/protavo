package drivertest

// // describeFetch defines the standard test suite for the Fetch operation.
// func describeFetch(
// 	before func() (driver.Driver, error),
// 	after func(),
// ) {
// 	ctx := context.Background()

// 	g.Describe("fetch operation", func() {
// 		var (
// 			dr driver.Driver
// 			db *protavo.DB
// 		)

// 		g.BeforeEach(func() {
// 			var err error
// 			dr, err = before()
// 			m.Expect(err).ShouldNot(m.HaveOccurred())

// 			db, err = protavo.NewDB(dr)
// 			m.Expect(err).ShouldNot(m.HaveOccurred())
// 		})

// 		g.AfterEach(func() {
// 			_ = db.Close()

// 			if after != nil {
// 				after()
// 			}
// 		})

// 		g.When("there are no documents in the database", func() {
// 			g.When("a nil filter is used", func() {
// 				g.It("does not return an error", func() {
// 					op := protavo.FetchAll(
// 						nil, // should never be invoked
// 					)

// 					err := db.Read(ctx, op)
// 					m.Expect(err).ShouldNot(m.HaveOccurred())
// 				})
// 			})

// 			g.When("an empty filter is used", func() {
// 				g.It("does not return an error", func() {
// 					op := protavo.FetchWhere(
// 						nil, // should never be invoked
// 					)

// 					err := db.Read(ctx, op)
// 					m.Expect(err).ShouldNot(m.HaveOccurred())
// 				})
// 			})
// 		})

// 		g.When("there are documents in the database", func() {
// 			g.BeforeEach(func() {
// 			})

// 			g.When("a nil filter is used", func() {
// 				g.It("does not return an error", func() {
// 					op := protavo.FetchAll(
// 						Each: func(*document.Document) (bool, error) {
// 							return false, nil
// 						},
// 						Filter: nil, // matches everything
// 					)

// 					err := db.Read(ctx, op)
// 					m.Expect(err).ShouldNot(m.HaveOccurred())
// 				})
// 			})

// 			g.When("an empty filter is used", func() {
// 				g.It("does not return an error", func() {
// 					op := protavo.FetchWhere(
// 						func(*document.Document) (bool, error) {
// 							return false, nil
// 						},
// 					)

// 					err := db.Read(ctx, op)
// 					m.Expect(err).ShouldNot(m.HaveOccurred())
// 				})
// 			})
// 		})
// 	})
// }
