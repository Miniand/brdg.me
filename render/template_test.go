package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Testlexer_readUntil(t *testing.T) {
	l := &lexer{input: "The {{c 'gray'}}egg{{_c}}"}
	assert.Equal(t, "", l.readUntil("The"))
	assert.Equal(t, "The ", l.readUntil("{{"))
	assert.Equal(t, "{{", l.readN(2))
	assert.Equal(t, "c", l.readWord())
	assert.Equal(t, " 'gray'", l.readUntil("}}"))
	assert.Equal(t, "}}", l.readN(2))
}
