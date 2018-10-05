package lang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAcceptLanguage(t *testing.T) {
	assert := assert.New(t)

	langQPair := ParseAcceptLanguage("en-US,en-US;q=0.2")
	expectedLangQPair := []LangQPair{
		LangQPair{Lang: "en-US", Q: 1},
		LangQPair{Lang: "en-US", Q: 0.2},
	}
	assert.Equal(expectedLangQPair, langQPair)
}
