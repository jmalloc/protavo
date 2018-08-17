package drivertest

import (
	"context"

	"github.com/jmalloc/protavo/src/protavo"
	"github.com/jmalloc/protavo/src/protavo/document"
	g "github.com/onsi/ginkgo"
	m "github.com/onsi/gomega"
)

// describeDeleteNamespace defines the standard test suite for the
// protavo.DeleteNamespace() operation.
func describeDeleteNamespace(
	before func() (*protavo.DB, error),
	after func(),
) {
	ctx := context.Background()
	var doc1, doc2 *document.Document

	g.Describe("DeleteNamespace", func() {
		var db *protavo.DB

		g.BeforeEach(func() {
			var err error
			db, err = before()
			m.Expect(err).ShouldNot(m.HaveOccurred())

			doc1 = &document.Document{
				ID:      "doc-1",
				Content: document.StringContent("content-1"),
			}

			doc2 = &document.Document{
				ID:      "doc-2",
				Content: document.StringContent("content-2"),
			}

			err = db.Save(ctx, doc1, doc2)
			m.Expect(err).ShouldNot(m.HaveOccurred())
		})

		g.AfterEach(func() {
			_ = db.Close()

			if after != nil {
				after()
			}
		})

		g.It("deletes all documents", func() {
			op := protavo.DeleteNamespace()

			err := db.Write(ctx, op)
			m.Expect(err).ShouldNot(m.HaveOccurred())

			docs, err := db.LoadAll(ctx)
			m.Expect(err).ShouldNot(m.HaveOccurred())
			m.Expect(docs).To(m.BeEmpty())
		})

		g.It("can be combined with save operations", func() {
			err := db.Write(
				ctx,
				protavo.DeleteNamespace(),
				protavo.Save(&document.Document{
					ID:      "doc-3",
					Content: document.StringContent("content-3"),
				}),
			)
			m.Expect(err).ShouldNot(m.HaveOccurred())

			docs, err := db.LoadAll(ctx)
			m.Expect(err).ShouldNot(m.HaveOccurred())
			m.Expect(docs).To(m.HaveLen(1))
			m.Expect(docs[0].ID).To(m.Equal("doc-3"))
		})
	})
}
