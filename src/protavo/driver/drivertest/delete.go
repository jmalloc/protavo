package drivertest

import (
	"context"

	"github.com/jmalloc/protavo/src/protavo"
	"github.com/jmalloc/protavo/src/protavo/document"
	g "github.com/onsi/ginkgo"
	m "github.com/onsi/gomega"
)

// describeDelete defines the standard test suite for the protavo.Delete()
// operation.
func describeDelete(
	before func() (*protavo.DB, error),
	after func(),
) {
	ctx := context.Background()
	var doc1, doc2 *document.Document

	g.Describe("Delete", func() {
		var db *protavo.DB

		g.BeforeEach(func() {
			var err error
			db, err = before()
			m.Expect(err).ShouldNot(m.HaveOccurred())

			doc1 = &document.Document{
				ID: "doc-1",
				Keys: document.Keys{
					"<unique-key>": document.UniqueKey,
				},
				Content: document.StringContent("<content-1>"),
			}

			doc2 = &document.Document{
				ID: "doc-2",
				Keys: document.Keys{
					"<unique-key>": document.UniqueKey,
				},
				Content: document.StringContent("<content-2>"),
			}
		})

		g.AfterEach(func() {
			_ = db.Close()

			if after != nil {
				after()
			}
		})

		g.When("the document exists", func() {
			g.BeforeEach(func() {
				err := db.Save(ctx, doc1)
				m.Expect(err).ShouldNot(m.HaveOccurred())
			})

			g.It("deletes the document", func() {
				op := protavo.Delete(doc1)

				err := db.Write(ctx, op)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				_, ok, err := db.Load(ctx, "doc-1")
				m.Expect(err).ShouldNot(m.HaveOccurred())
				m.Expect(ok).To(m.BeFalse())
			})

			g.It("returns an error if the provided revision is not correct", func() {
				doc1.Revision = 123
				op := protavo.Delete(doc1)

				err := db.Write(ctx, op)
				m.Expect(err).To(m.MatchError(
					&protavo.OptimisticLockError{
						DocumentID: "doc-1",
						GivenRev:   123,
						ActualRev:  1,
						Operation:  "delete",
					},
				))
			})
		})

		g.When("the document does not exist", func() {
			g.It("does not return an error", func() {
				op := protavo.Delete(doc1)

				err := db.Write(ctx, op)
				m.Expect(err).ShouldNot(m.HaveOccurred())
			})

			g.It("returns an error if the provided revision is not correct", func() {
				doc1.Revision = 123
				op := protavo.Delete(doc1)

				err := db.Write(ctx, op)
				m.Expect(err).To(m.MatchError(
					&protavo.OptimisticLockError{
						DocumentID: "doc-1",
						GivenRev:   123,
						ActualRev:  0,
						Operation:  "delete",
					},
				))
			})
		})

		g.It("allows other documents to use the delete document's unique key", func() {
			op := protavo.Delete(doc1)

			err := db.Write(ctx, op)
			m.Expect(err).ShouldNot(m.HaveOccurred())

			err = db.Save(ctx, doc2)
			m.Expect(err).ShouldNot(m.HaveOccurred())
		})
	})
}
