package commander

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenize(t *testing.T) {
	tokens := NewCommand("say    <input>   ").Tokenize()
	assert.Equal(t, len(tokens), 2)
	assert.Equal(t, tokens[0].Word, "say")
	assert.False(t, tokens[0].IsParameter)
	assert.Equal(t, tokens[1].Word, "input")
	assert.True(t, tokens[1].IsParameter)

	tokens = NewCommand("search <pattern>").Tokenize()
	assert.Equal(t, len(tokens), 2)
	assert.Equal(t, tokens[0].Word, "search")
	assert.False(t, tokens[0].IsParameter)
	assert.Equal(t, tokens[1].Word, "pattern")
	assert.True(t, tokens[1].IsParameter)

	tokens = NewCommand("<a> <123> <a123> <a-123> <a.123> b> <c e").Tokenize()
	assert.Equal(t, len(tokens), 8)
	assert.Equal(t, tokens[0].Word, "a")
	assert.True(t, tokens[0].IsParameter)
	assert.Equal(t, tokens[1].Word, "123")
	assert.True(t, tokens[1].IsParameter)
	assert.Equal(t, tokens[2].Word, "a123")
	assert.True(t, tokens[2].IsParameter)
	assert.Equal(t, tokens[3].Word, "a-123")
	assert.True(t, tokens[3].IsParameter)
	assert.Equal(t, tokens[4].Word, "a.123")
	assert.True(t, tokens[4].IsParameter)
	assert.Equal(t, tokens[5].Word, "b>")
	assert.False(t, tokens[5].IsParameter)
	assert.Equal(t, tokens[6].Word, "<c")
	assert.False(t, tokens[6].IsParameter)
	assert.Equal(t, tokens[7].Word, "e")
	assert.False(t, tokens[7].IsParameter)

	tokens = NewCommand("\\ ( ) { } [ ] ? . + | ^ $").Tokenize()
	assert.Equal(t, len(tokens), 13)
	assert.Equal(t, tokens[0].Word, "\\")
	assert.False(t, tokens[0].IsParameter)
	assert.Equal(t, tokens[1].Word, "(")
	assert.False(t, tokens[1].IsParameter)
	assert.Equal(t, tokens[2].Word, ")")
	assert.False(t, tokens[2].IsParameter)
	assert.Equal(t, tokens[3].Word, "{")
	assert.False(t, tokens[3].IsParameter)
	assert.Equal(t, tokens[4].Word, "}")
	assert.False(t, tokens[4].IsParameter)
	assert.Equal(t, tokens[5].Word, "[")
	assert.False(t, tokens[5].IsParameter)
	assert.Equal(t, tokens[6].Word, "]")
	assert.False(t, tokens[6].IsParameter)
	assert.Equal(t, tokens[7].Word, "?")
	assert.False(t, tokens[7].IsParameter)
	assert.Equal(t, tokens[8].Word, ".")
	assert.False(t, tokens[8].IsParameter)
	assert.Equal(t, tokens[9].Word, "+")
	assert.False(t, tokens[9].IsParameter)
	assert.Equal(t, tokens[10].Word, "|")
	assert.False(t, tokens[10].IsParameter)
	assert.Equal(t, tokens[11].Word, "^")
	assert.False(t, tokens[11].IsParameter)
	assert.Equal(t, tokens[12].Word, "$")
	assert.False(t, tokens[12].IsParameter)
}

func TestMatch(t *testing.T) {
	properties, isMatch := NewCommand("").Match("ping")
	assert.False(t, isMatch)
	assert.Nil(t, properties)

	properties, isMatch = NewCommand("").Match("")
	assert.False(t, isMatch)
	assert.Nil(t, properties)

	properties, isMatch = NewCommand("ping").Match("ping")
	assert.True(t, isMatch)
	assert.NotNil(t, properties)

	properties, isMatch = NewCommand("ping").Match("pong")
	assert.False(t, isMatch)
	assert.Nil(t, properties)

	properties, isMatch = NewCommand("abc").Match(".abc.")
	assert.False(t, isMatch)
	assert.Nil(t, properties)

	properties, isMatch = NewCommand("help").Match("helpful")
	assert.False(t, isMatch)
	assert.Nil(t, properties)

	properties, isMatch = NewCommand("search all").Match("search all")
	assert.True(t, isMatch)
	assert.NotNil(t, properties)

	properties, isMatch = NewCommand("search all").Match("search     all")
	assert.True(t, isMatch)
	assert.NotNil(t, properties)

	properties, isMatch = NewCommand("search all").Match("search for all")
	assert.False(t, isMatch)
	assert.Nil(t, properties)

	properties, isMatch = NewCommand("search all").Match("searchall")
	assert.False(t, isMatch)
	assert.Nil(t, properties)

	properties, isMatch = NewCommand("help").Match("Could you help me?")
	assert.True(t, isMatch)
	assert.NotNil(t, properties)

	properties, isMatch = NewCommand("help me").Match("Could you help me?")
	assert.False(t, isMatch)
	assert.Nil(t, properties)

	properties, isMatch = NewCommand("help me").Match("please help me")
	assert.True(t, isMatch)
	assert.NotNil(t, properties)

	properties, isMatch = NewCommand("echo <word>").Match("echo")
	assert.True(t, isMatch)
	assert.NotNil(t, properties)

	properties, isMatch = NewCommand("echo <word>").Match("echo hey")
	assert.True(t, isMatch)
	assert.Equal(t, properties.StringParam("word", ""), "hey")

	properties, isMatch = NewCommand("search <pattern>").Match("search *")
	assert.True(t, isMatch)
	assert.Equal(t, properties.StringParam("pattern", ""), "*")

	properties, isMatch = NewCommand("repeat <word> <number>").Match("repeat hey 5")
	assert.True(t, isMatch)
	assert.Equal(t, properties.StringParam("word", ""), "hey")
	assert.Equal(t, properties.IntegerParam("number", 0), 5)

	properties, isMatch = NewCommand("repeat <word> <number>").Match("repeat hey")
	assert.True(t, isMatch)
	assert.Equal(t, properties.StringParam("word", ""), "hey")
	assert.Equal(t, properties.IntegerParam("number", 0), 0)

	properties, isMatch = NewCommand("calculate <number1> plus <number2>").Match("calculate 10 plus 5")
	assert.True(t, isMatch)
	assert.Equal(t, properties.IntegerParam("number1", 0), 10)
	assert.Equal(t, properties.IntegerParam("number2", 0), 5)

	properties, isMatch = NewCommand("<number1> + <number2>").Match("10 + 5")
	assert.True(t, isMatch)
	assert.Equal(t, properties.IntegerParam("number1", 0), 10)
	assert.Equal(t, properties.IntegerParam("number2", 0), 5)

	properties, isMatch = NewCommand("<number1> + <number2>").Match("+")
	assert.True(t, isMatch)
	assert.Equal(t, properties.IntegerParam("number1", 0), 0)
	assert.Equal(t, properties.IntegerParam("number2", 0), 0)

	properties, isMatch = NewCommand("\\ ( ) { } [ ] ? . + | ^ $").Match("\\ ( ) { } [ ] ? . + | ^ $")
	assert.True(t, isMatch)
	assert.NotNil(t, properties)
}
