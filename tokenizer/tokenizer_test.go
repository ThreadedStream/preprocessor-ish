package tokenizer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsKeyword(t *testing.T) {
	var (
		keyword1     = "break"
		keyword2     = "case"
		keyword3     = "for"
		notAKeyword1 = "garcon"
	)

	assert.True(t, isKeyword(keyword1))
	assert.True(t, isKeyword(keyword2))
	assert.True(t, isKeyword(keyword3))
	assert.False(t, isKeyword(notAKeyword1))
}
