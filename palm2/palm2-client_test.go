package palm2_test

import (
	"testing"

	"github.com/raghvenders/generative-ai-go/palm2"
	"github.com/stretchr/testify/assert"
)

// TO DO On Tests
func TestAPI(t *testing.T) {
	t.Parallel()

	llm, err := palm2.NewPalm()
	assert.Nil(t, llm)

	_ = err

}
