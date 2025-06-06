package dynflags_test

import (
	"testing"

	"github.com/containeroo/dynflags"
	"github.com/stretchr/testify/assert"
)

func TestStringValue(t *testing.T) {
	t.Parallel()

	t.Run("Parse valid string", func(t *testing.T) {
		t.Parallel()

		stringValue := dynflags.StringValue{}
		parsed, err := stringValue.Parse("example")
		assert.NoError(t, err)
		assert.Equal(t, "example", parsed)
	})

	t.Run("Set valid string", func(t *testing.T) {
		t.Parallel()

		bound := "initial"
		stringValue := dynflags.StringValue{Bound: &bound}

		err := stringValue.Set("updated")
		assert.NoError(t, err)
		assert.Equal(t, "updated", bound)
	})

	t.Run("Set invalid type", func(t *testing.T) {
		t.Parallel()

		bound := "initial"
		stringValue := dynflags.StringValue{Bound: &bound}

		err := stringValue.Set(123) // Invalid type
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid value type: expected string")
	})
}

func TestGroupConfigString(t *testing.T) {
	t.Parallel()

	t.Run("Define string flag", func(t *testing.T) {
		t.Parallel()

		group := &dynflags.ConfigGroup{Flags: make(map[string]*dynflags.Flag)}
		defaultValue := "default"
		group.String("stringFlag", defaultValue, "A string flag")

		assert.Contains(t, group.Flags, "stringFlag")
		assert.Equal(t, "A string flag", group.Flags["stringFlag"].Usage)
		assert.Equal(t, defaultValue, group.Flags["stringFlag"].Default)
	})
}

func TestParsedGroupGetString(t *testing.T) {
	t.Parallel()

	t.Run("Get existing string flag", func(t *testing.T) {
		t.Parallel()

		group := &dynflags.ParsedGroup{
			Values: map[string]any{
				"stringFlag": "value",
			},
		}
		str, err := group.GetString("stringFlag")
		assert.NoError(t, err)
		assert.Equal(t, "value", str)
	})

	t.Run("Get non-existent string flag", func(t *testing.T) {
		t.Parallel()

		group := &dynflags.ParsedGroup{
			Values: map[string]any{},
		}
		str, err := group.GetString("stringFlag")
		assert.Error(t, err)
		assert.Equal(t, "", str)
		assert.EqualError(t, err, "flag 'stringFlag' not found in group ''")
	})

	t.Run("Get string flag with invalid type", func(t *testing.T) {
		t.Parallel()

		group := &dynflags.ParsedGroup{
			Values: map[string]any{
				"stringFlag": 123, // Invalid type
			},
		}
		str, err := group.GetString("stringFlag")
		assert.Error(t, err)
		assert.Equal(t, "", str)
		assert.EqualError(t, err, "flag 'stringFlag' is not a string")
	})
}

func TestStringGetBound(t *testing.T) {
	t.Run("StringValue - GetBound", func(t *testing.T) {
		var str *string
		value := "test"
		str = &value

		stringValue := dynflags.StringValue{Bound: str}
		assert.Equal(t, "test", stringValue.GetBound())

		stringValue = dynflags.StringValue{Bound: nil}
		assert.Nil(t, stringValue.GetBound())
	})
}
