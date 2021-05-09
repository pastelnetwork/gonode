package storage_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/storage/definition"
)

var _ = Describe("Memory storage test", func() {
	var memStore definition.KeyValue

	When("Create memory storage", func() {
		It("memory instance should not nil", func() {
			memStore = storage.NewKeyValue()
			Expect(memStore).ShouldNot(BeNil())
		})
	})

	When("Set key value to memory storage", func() {
		It("Should not throw error", func() {
			err := memStore.Set("test", []byte("value"))
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	When("Get key value from memory storage", func() {
		var data []byte
		var err error

		It("should not throw error", func() {
			data, err = memStore.Get("test")
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("data string shuold value", func() {
			Expect(string(data)).Should(Equal("value"))
		})
	})

	When("Delete key value from storage", func() {

		It("should not throw error", func() {
			err := memStore.Delete("test")
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("test key should not exist", func() {
			_, err := memStore.Get("test")
			Expect(err).Should(HaveOccurred())
		})
	})

})
