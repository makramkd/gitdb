package gitdb_test

import (
	"testing"

	"github.com/makramkd/gitdb"
	"github.com/stretchr/testify/assert"
)

func TestGetInstance(t *testing.T) {
	i := gitdb.GetInstance()
	i2 := gitdb.GetInstance()

	assert.Equal(t, i, i2)
}
