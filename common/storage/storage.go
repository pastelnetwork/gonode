package storage

import (
	"github.com/pastelnetwork/gonode/common/storage/definition"
	"github.com/pastelnetwork/gonode/common/storage/memory"
)

//NewKeyValue return new instance key value storage
func NewKeyValue() definition.KeyValue {
	return memory.NewMemory()
}
