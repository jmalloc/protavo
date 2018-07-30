package drivertest

// 	g.Describe("Delete", func() {
// 		g.BeforeEach(func() {
// 			savedDoc := protavo.NewDocument(
// 				"doc-id",
// 				&TestContent{
// 					Data: "<content-1>",
// 				},
// 			).WithUniqueKeys("<uniq>")

// 			var err error
// 			savedDoc, err = db.Save(ctx, savedDoc)
// 			m.Expect(err).ShouldNot(m.HaveOccurred())

// 			err = db.Delete(ctx, savedDoc)
// 			m.Expect(err).ShouldNot(m.HaveOccurred())
// 		})

// 		g.It("removes the document", func() {
// 			_, ok, err := db.Load(ctx, "doc-id")
// 			m.Expect(err).ShouldNot(m.HaveOccurred())
// 			m.Expect(ok).To(m.BeFalse())
// 		})

// 		g.It("allows other documents with the same unique key to be saved", func() {
// 			doc := protavo.NewDocument(
// 				"doc-id",
// 				&TestContent{
// 					Data: "<content-1>",
// 				},
// 			).WithUniqueKeys("<uniq>")

// 			_, err := db.Save(ctx, doc)
// 			m.Expect(err).ShouldNot(m.HaveOccurred())
// 		})
// 	})
// })
// }
