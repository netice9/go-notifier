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
	It("should notify listeners of changes", func(done Done) {
		n := notifier.NewNotifier()
		c := make(chan interface{})
		n.AddListener(c)
		n.Notify("test")
		notification := <-c
		Expect(notification).To(Equal("test"))
		close(done)
	})
})
