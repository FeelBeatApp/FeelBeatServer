package networking_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Just an example test
func TestHubBroadcastsMessages(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(2+2, 4)
}
