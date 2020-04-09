package test

import (
	"testing"

	_ "github.com/go-kivik/memorydb/v3"
)

func init() {
	RegisterMemoryDBSuite()
}

func TestMemory(t *testing.T) {
	MemoryTest(t)
}
