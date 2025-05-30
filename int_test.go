package dynflags_test

import (
	"testing"

	"github.com/containeroo/dynflags"
	"github.com/stretchr/testify/assert"
)

func TestIntValue_Parse(t *testing.T) {
	t.Parallel()

	t.Run("ValidInt", func(t *testing.T) {
		t.Parallel()

		var bound int
		val := &dynflags.IntValue{Bound: &bound}
		parsed, err := val.Parse("42")
		assert.NoError(t, err)
		assert.Equal(t, 42, parsed)
	})

	t.Run("InvalidInt", func(t *testing.T) {
		t.Parallel()

		var bound int
		val := &dynflags.IntValue{Bound: &bound}
		_, err := val.Parse("invalid")
		assert.Error(t, err)
	})
}

func TestIntValue_Set(t *testing.T) {
	t.Parallel()

	t.Run("ValidInt", func(t *testing.T) {
		t.Parallel()

		var bound int
		val := &dynflags.IntValue{Bound: &bound}
		assert.NoError(t, val.Set(42))
		assert.Equal(t, 42, bound)
	})

	t.Run("InvalidType", func(t *testing.T) {
		t.Parallel()

		var bound int
		val := &dynflags.IntValue{Bound: &bound}
		assert.Error(t, val.Set("not an int"))
	})
}

func TestGroupConfig_Int(t *testing.T) {
	t.Parallel()

	t.Run("DefineAndRetrieveInt", func(t *testing.T) {
		t.Parallel()

		group := &dynflags.ConfigGroup{Flags: make(map[string]*dynflags.Flag)}
		bound := group.Int("test-int", 100, "Test integer flag")
		assert.NotNil(t, bound)

		flag, exists := group.Flags["test-int"]
		assert.True(t, exists)
		assert.NotNil(t, flag)
		assert.Equal(t, dynflags.FlagTypeInt, flag.Type)
		assert.Equal(t, 100, flag.Default)
		assert.Equal(t, "Test integer flag", flag.Usage)
	})
}

func TestParsedGroup_GetInt(t *testing.T) {
	t.Parallel()

	t.Run("ValidIntRetrieval", func(t *testing.T) {
		t.Parallel()

		parsedGroup := &dynflags.ParsedGroup{
			Values: map[string]any{
				"test-int": 42,
			},
		}
		val, err := parsedGroup.GetInt("test-int")
		assert.NoError(t, err)
		assert.Equal(t, 42, val)
	})

	t.Run("FlagNotFound", func(t *testing.T) {
		t.Parallel()

		parsedGroup := &dynflags.ParsedGroup{
			Values: make(map[string]any),
		}
		_, err := parsedGroup.GetInt("non-existent")
		assert.Error(t, err)
	})

	t.Run("InvalidType", func(t *testing.T) {
		t.Parallel()

		parsedGroup := &dynflags.ParsedGroup{
			Values: map[string]any{
				"test-int": "not an int",
			},
		}
		_, err := parsedGroup.GetInt("test-int")
		assert.Error(t, err)
	})
}

func TestIntGetBound(t *testing.T) {
	t.Run("IntValue - GetBound", func(t *testing.T) {
		var i *int
		val := 42
		i = &val

		intValue := dynflags.IntValue{Bound: i}
		assert.Equal(t, 42, intValue.GetBound())

		intValue = dynflags.IntValue{Bound: nil}
		assert.Nil(t, intValue.GetBound())
	})
}
