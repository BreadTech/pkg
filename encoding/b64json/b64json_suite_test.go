package b64json_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestB64json(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "B64json Suite")
}
