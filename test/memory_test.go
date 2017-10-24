package test

import (
	"testing"

	_ "github.com/go-kivik/memorydb"
)

func init() {
	RegisterMemoryDBSuite()
}

func TestMemory(t *testing.T) {
	MemoryTest(t)
}
