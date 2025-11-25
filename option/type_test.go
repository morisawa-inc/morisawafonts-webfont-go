package option

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientOptions_Merge(t *testing.T) {
	base := NewClientOptions(WithAPIToken("token-1"))
	merged := base.Merge(WithAPIToken("token-2"))

	assert.Equal(t, base.APIToken, "token-1")
	assert.Equal(t, merged.APIToken, "token-2")
}
