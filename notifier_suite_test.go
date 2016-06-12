package notifier_test

import (
	"github.com/netice9/notifier-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGoNotifier(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Notifier Suite")
}

var _ = Describe("Notifier", func() {

	var n *notifier.Notifier

	BeforeEach(func() {
		n = notifier.NewNotifier()
	})

	It("should notify listeners of changes", func(done Done) {
		c := make(chan interface{})
		n.AddListener(c)
		n.Notify("test")
		notification := <-c
		Expect(notification).To(Equal("test"))
		close(done)
	})

	Context("when listener is added", func() {

		var l1 chan interface{}

		BeforeEach(func() {
			l1 = make(chan interface{})
			n.AddListener(l1)
		})

		Describe("RemoveListener()", func() {
			It("should remove existing listener", func() {
				n.RemoveListener(l1)
				Expect(n.NumberOfListeners()).To(Equal(0))
			})
		})
	})
})
