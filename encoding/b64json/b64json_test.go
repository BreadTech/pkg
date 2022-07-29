package b64json_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/breadtech/pkg/encoding/b64json"
)

type resp struct {
	Status, Schemes, Scope string
}

var _ = Describe("B64json", func() {
	Context("Marshal and Unmarshal", func() {
		It("should work", func() {
			testobj := &resp{
				Status:  "400",
				Schemes: "Basic",
				Scope:   "https://mail.breadtech.com/",
			}
			b, err := b64json.Marshal(testobj)
			Expect(err).ToNot(HaveOccurred())

			obj := new(resp)
			Expect(b64json.Unmarshal(b, obj)).ToNot(HaveOccurred())
			Expect(obj.Status).To(Equal(testobj.Status))
			Expect(obj.Schemes).To(Equal(testobj.Schemes))
			Expect(obj.Scope).To(Equal(testobj.Scope))
		})
	})
})
