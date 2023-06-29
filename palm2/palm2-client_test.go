package palm2_test

import (
	"testing"

	"github.com/raghvenders/generative-ai-go/palm2"
	"github.com/stretchr/testify/assert"
)

func TestAPI(t *testing.T) {
	t.Parallel()

	llm, err := palm2.NewPalm()
	assert.ErrorIs(t, err, palm2.ErrMissingPalmAPIKey)
	assert.Nil(t, llm)

}
