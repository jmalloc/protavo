package drivertest

// g.Describe(name, func() {
// 	ctx := context.Background()

// 	var (
// 		driver driver.Driver
// 		db     *protavo.DB
// 	)

// 	g.BeforeEach(func() {
// 		var err error
// 		driver, err = before()
// 		m.Expect(err).ShouldNot(m.HaveOccurred())

// 		db, err = protavo.NewDB(driver)
// 		m.Expect(err).ShouldNot(m.HaveOccurred())
// 	})

// 	g.AfterEach(func() {
// 		db.Close()

// 		if after != nil {
// 			after()
// 		}
// 	})

// 	g.Describe("Load", func() {
// 		g.When("the document exists", func() {
// 			var doc *document.Document

// 			g.BeforeEach(func() {
// 				doc = &protavo.Document{
// 					ID: "doc-id",
// 					Content: &TestContent{
// 						Data: "<content-1>",
// 					},
// 				}

// 				var err error
// 				doc, err = db.Save(ctx, doc)
// 				m.Expect(err).ShouldNot(m.HaveOccurred())
// 			})

// 			g.It("returns true", func() {
// 				_, ok, err := db.Load(ctx, "doc-id")
// 				m.Expect(err).ShouldNot(m.HaveOccurred())
// 				m.Expect(ok).To(m.BeTrue())
// 			})

// 			g.It("returns the document", func() {
// 				d, _, err := db.Load(ctx, "doc-id")
// 				m.Expect(err).ShouldNot(m.HaveOccurred())
// 				m.Expect(d.Equal(doc)).To(m.BeTrue())
// 			})
// 		})

// 		g.When("the document does not exist", func() {
// 			g.It("returns false", func() {
// 				_, ok, err := db.Load(ctx, "doc-id")
// 				m.Expect(err).ShouldNot(m.HaveOccurred())
// 				m.Expect(ok).To(m.BeFalse())
// 			})
// 		})
// 	})

// 	g.Describe("Save", func() {
// 		g.When("a new document is created", func() {
// 			var (
// 				savedDoc    *document.Document
// 				returnedDoc *document.Document
// 			)

// 			g.BeforeEach(func() {
// 				savedDoc = &protavo.Document{
// 					ID: "doc-id",
// 					Content: &TestContent{
// 						Data: "<content-1>",
// 					},
// 				}

// 				var err error
// 				returnedDoc, err = db.Save(ctx, savedDoc)
// 				m.Expect(err).ShouldNot(m.HaveOccurred())
// 			})

// 			g.It("returns a copy of the document at version 1", func() {
// 				m.Expect(returnedDoc.Revision).To(
// 					m.BeNumerically("==", 1),
// 				)
// 			})

// 			g.It("sets the created and updated timestamps", func() {
// 				m.Expect(returnedDoc.CreatedAt).To(
// 					m.BeTemporally("~", time.Now(), 50*time.Millisecond),
// 				)

// 				m.Expect(returnedDoc.UpdatedAt).To(
// 					m.BeTemporally("==", returnedDoc.CreatedAt),
// 				)
// 			})

// 			g.It("persists the document as returned", func() {
// 				d, _, err := db.Load(ctx, "doc-id")
// 				m.Expect(err).ShouldNot(m.HaveOccurred())
// 				m.Expect(d.Equal(returnedDoc)).To(m.BeTrue())
// 			})

// 			g.It("returns an optimistic lock error if the version is non-zero", func() {
// 				// first delete the document
// 				err := db.Delete(ctx, returnedDoc)
// 				m.Expect(err).ShouldNot(m.HaveOccurred())

// 				// then try to save it with the version of the document @ v1
// 				_, err = db.Save(ctx, returnedDoc)
// 				m.Expect(err).To(
// 					m.MatchError(&protavo.OptimisticLockError{
// 						DocumentID:    "doc-id",
// 						GivenVersion:  1,
// 						ActualVersion: 0,
// 						Action:        "save",
// 					}),
// 				)
// 			})
// 		})
// 	})

// 	g.When("an existing document is updated", func() {
// 		var (
// 			savedDoc    *document.Document
// 			returnedDoc *document.Document
// 		)

// 		g.BeforeEach(func() {
// 			savedDoc = &protavo.Document{
// 				ID: "doc-id",
// 				Content: &TestContent{
// 					Data: "<content-1>",
// 				},
// 			}

// 			var err error
// 			savedDoc, err = db.Save(ctx, savedDoc)
// 			m.Expect(err).ShouldNot(m.HaveOccurred())

// 			savedDoc = savedDoc.WithContent(
// 				&TestContent{
// 					Data: "<content-2>",
// 				},
// 			)

// 			returnedDoc, err = db.Save(ctx, savedDoc)
// 			m.Expect(err).ShouldNot(m.HaveOccurred())
// 		})

// 		g.It("returns a copy of the document at version 2", func() {
// 			m.Expect(returnedDoc.Version()).To(
// 				m.BeNumerically("==", 2),
// 			)
// 		})

// 		g.It("does not modify the created timestamp", func() {
// 			after, _ := returnedDoc.CreatedAt()
// 			before, _ := savedDoc.CreatedAt()

// 			m.Expect(after).To(m.Equal(before))
// 		})

// 		g.It("sets the updated timestamp", func() {
// 			t, ok := returnedDoc.UpdatedAt()
// 			m.Expect(ok).To(m.BeTrue())
// 			m.Expect(t).To(
// 				m.BeTemporally("~", time.Now(), 50*time.Millisecond),
// 			)
// 		})

// 		g.It("persists the document as returned", func() {
// 			d, _, err := db.Load(ctx, "doc-id")
// 			m.Expect(err).ShouldNot(m.HaveOccurred())
// 			m.Expect(d.Equal(returnedDoc)).To(m.BeTrue())
// 		})

// 		g.It("returns an optimistic lock error if the version is not equal", func() {
// 			_, err := db.Save(ctx, savedDoc)
// 			m.Expect(err).To(
// 				m.MatchError(&protavo.OptimisticLockError{
// 					DocumentID:    "doc-id",
// 					GivenVersion:  1,
// 					ActualVersion: 2,
// 					Action:        "save",
// 				}),
// 			)
// 		})
// 	})

// 	g.When("the document has a unique key", func() {
// 		var (
// 			savedDoc *document.Document
// 		)

// 		g.BeforeEach(func() {
// 			savedDoc = protavo.NewDocument(
// 				"doc-1",
// 				&TestContent{
// 					Data: "<content-1>",
// 				},
// 			).WithUniqueKeys("<uniq>")

// 			var err error
// 			savedDoc, err = db.Save(ctx, savedDoc)
// 			m.Expect(err).ShouldNot(m.HaveOccurred())
// 		})

// 		g.It("can be updated while retaining the unique key", func() {
// 			doc := savedDoc.WithContent(
// 				&TestContent{
// 					Data: "<content-3>",
// 				},
// 			)

// 			_, err := db.Save(ctx, doc)
// 			m.Expect(err).ShouldNot(m.HaveOccurred())
// 		})

// 		g.It("returns a duplicate key error if another document is saved with the same unique key", func() {
// 			doc := protavo.NewDocument(
// 				"doc-2",
// 				&TestContent{
// 					Data: "<content-1>",
// 				},
// 			).WithUniqueKeys("<uniq>")

// 			_, err := db.Save(ctx, doc)
// 			m.Expect(err).To(
// 				m.MatchError(&protavo.DuplicateKeyError{
// 					DocumentID:            "doc-2",
// 					ConflictingDocumentID: "doc-1",
// 					UniqueKey:             "<uniq>",
// 				}),
// 			)
// 		})
// 	})

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
