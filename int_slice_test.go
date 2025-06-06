package dynflags_test

import (
	"testing"

	"github.com/containeroo/dynflags"
	"github.com/stretchr/testify/assert"
)

func TestIntSlicesValue(t *testing.T) {
	t.Parallel()

	t.Run("Parse valid int slice value", func(t *testing.T) {
		t.Parallel()

		intSlicesValue := dynflags.IntSlicesValue{Bound: &[]int{}}
		parsed, err := intSlicesValue.Parse("123")
		assert.NoError(t, err)
		assert.Equal(t, 123, parsed)
	})

	t.Run("Parse invalid int slice value", func(t *testing.T) {
		t.Parallel()

		intSlicesValue := dynflags.IntSlicesValue{Bound: &[]int{}}
		parsed, err := intSlicesValue.Parse("invalid")
		assert.Error(t, err)
		assert.Nil(t, parsed)
	})

	t.Run("Set valid int slice value", func(t *testing.T) {
		t.Parallel()

		bound := []int{1}
		intSlicesValue := dynflags.IntSlicesValue{Bound: &bound}

		err := intSlicesValue.Set(2)
		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2}, bound)
	})

	t.Run("Set invalid type", func(t *testing.T) {
		t.Parallel()

		bound := []int{1}
		intSlicesValue := dynflags.IntSlicesValue{Bound: &bound}

		err := intSlicesValue.Set("invalid") // Invalid type
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid value type: expected int")
	})
}

func TestGroupConfigIntSlices(t *testing.T) {
	t.Parallel()

	t.Run("Define int slices flag", func(t *testing.T) {
		t.Parallel()

		group := &dynflags.ConfigGroup{Flags: make(map[string]*dynflags.Flag)}
		defaultValue := []int{1, 2}
		group.IntSlices("intSliceFlag", defaultValue, "An int slices flag")

		assert.Contains(t, group.Flags, "intSliceFlag")
		assert.Equal(t, "An int slices flag", group.Flags["intSliceFlag"].Usage)
		assert.Equal(t, "1,2", group.Flags["intSliceFlag"].Default)
	})
}

func TestGetIntSlices(t *testing.T) {
	t.Parallel()

	t.Run("Retrieve []int value", func(t *testing.T) {
		t.Parallel()

		parsedGroup := &dynflags.ParsedGroup{
			Name:   "testGroup",
			Values: map[string]any{"flag1": []int{1, 2, 3}},
		}

		result, err := parsedGroup.GetIntSlices("flag1")
		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("Retrieve single int value as []int", func(t *testing.T) {
		t.Parallel()

		parsedGroup := &dynflags.ParsedGroup{
			Name:   "testGroup",
			Values: map[string]any{"flag1": 42},
		}

		result, err := parsedGroup.GetIntSlices("flag1")
		assert.NoError(t, err)
		assert.Equal(t, []int{42}, result)
	})

	t.Run("Flag not found", func(t *testing.T) {
		t.Parallel()

		parsedGroup := &dynflags.ParsedGroup{
			Name:   "testGroup",
			Values: map[string]any{},
		}

		result, err := parsedGroup.GetIntSlices("nonExistentFlag")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "flag 'nonExistentFlag' not found in group 'testGroup'")
	})

	t.Run("Flag value is invalid type", func(t *testing.T) {
		t.Parallel()

		parsedGroup := &dynflags.ParsedGroup{
			Name:   "testGroup",
			Values: map[string]any{"flag1": "invalid"},
		}

		result, err := parsedGroup.GetIntSlices("flag1")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "flag 'flag1' is not a []int")
	})
}

func TestIntSlicesGetBound(t *testing.T) {
	t.Run("IntSlicesValue - GetBound", func(t *testing.T) {
		var slices *[]int
		val := []int{1, 2, 3}
		slices = &val

		intSlicesValue := dynflags.IntSlicesValue{Bound: slices}
		assert.Equal(t, val, intSlicesValue.GetBound())

		intSlicesValue = dynflags.IntSlicesValue{Bound: nil}
		assert.Nil(t, intSlicesValue.GetBound())
	})
}
