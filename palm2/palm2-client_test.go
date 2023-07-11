package palm2_test

import (
	"testing"

	"github.com/raghvenders/generative-ai-go/palm2"
	"github.com/stretchr/testify/assert"
)

// TO DO On Tests
func TestAPIMissingKey(t *testing.T) {
	t.Parallel()

	llm, err := palm2.NewPalm()
	assert.Nil(t, llm)

	for _, err = range err.(interface{ Unwrap() []error }).Unwrap() {
		assert.ErrorIs(t, err, palm2.ErrMissingPalmAPIKey)
	}

}
