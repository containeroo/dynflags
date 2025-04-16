package dynflags_test

import (
	"testing"

	"github.com/containeroo/dynflags"
	"github.com/stretchr/testify/assert"
)

func TestListenAddrSlicesValue(t *testing.T) {
	t.Parallel()

	t.Run("Parse valid listen address", func(t *testing.T) {
		t.Parallel()

		value := dynflags.ListenAddrSlicesValue{Bound: &[]string{}}
		parsed, err := value.Parse(":8080")
		assert.NoError(t, err)
		assert.Equal(t, ":8080", parsed)
	})

	t.Run("Parse invalid listen address", func(t *testing.T) {
		t.Parallel()

		value := dynflags.ListenAddrSlicesValue{Bound: &[]string{}}
		parsed, err := value.Parse("bad-address")
		assert.Error(t, err)
		assert.Nil(t, parsed)
	})

	t.Run("Set valid listen address", func(t *testing.T) {
		t.Parallel()

		bound := []string{":9090"}
		value := dynflags.ListenAddrSlicesValue{Bound: &bound}

		err := value.Set(":8080")
		assert.NoError(t, err)
		assert.Equal(t, []string{":9090", ":8080"}, bound)
	})

	t.Run("Set invalid type", func(t *testing.T) {
		t.Parallel()

		bound := []string{}
		value := dynflags.ListenAddrSlicesValue{Bound: &bound}

		err := value.Set(12345)
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid value type: expected string listen address")
	})
}

func TestGroupConfigListenAddrSlices(t *testing.T) {
	t.Parallel()

	t.Run("Define listen address slices flag", func(t *testing.T) {
		t.Parallel()

		group := &dynflags.ConfigGroup{Flags: make(map[string]*dynflags.Flag)}
		defaultValue := []string{":9090", "127.0.0.1:9091"}
		group.ListenAddrSlices("listenSlice", defaultValue, "Multiple listen addresses")

		assert.Contains(t, group.Flags, "listenSlice")
		assert.Equal(t, "Multiple listen addresses", group.Flags["listenSlice"].Usage)
		assert.Equal(t, ":9090,127.0.0.1:9091", group.Flags["listenSlice"].Default)
	})

	t.Run("Define listen address slices with invalid default", func(t *testing.T) {
		t.Parallel()

		group := &dynflags.ConfigGroup{Flags: make(map[string]*dynflags.Flag)}

		assert.PanicsWithValue(t,
			"listenSlice has an invalid default listen address 'bad-address': address bad-address: missing port in address",
			func() {
				group.ListenAddrSlices("listenSlice", []string{":8080", "bad-address"}, "Invalid default")
			})
	})
}

func TestGetListenAddrSlices(t *testing.T) {
	t.Parallel()

	t.Run("Retrieve []string listen address slice", func(t *testing.T) {
		t.Parallel()

		parsedGroup := &dynflags.ParsedGroup{
			Name:   "testGroup",
			Values: map[string]any{"flag1": []string{":8080", "127.0.0.1:9090"}},
		}

		result, err := parsedGroup.GetListenAddrSlices("flag1")
		assert.NoError(t, err)
		assert.Equal(t, []string{":8080", "127.0.0.1:9090"}, result)
	})

	t.Run("Retrieve single string as []string", func(t *testing.T) {
		t.Parallel()

		parsedGroup := &dynflags.ParsedGroup{
			Name:   "testGroup",
			Values: map[string]any{"flag1": ":8080"},
		}

		result, err := parsedGroup.GetListenAddrSlices("flag1")
		assert.NoError(t, err)
		assert.Equal(t, []string{":8080"}, result)
	})

	t.Run("Flag not found", func(t *testing.T) {
		t.Parallel()

		parsedGroup := &dynflags.ParsedGroup{
			Name:   "testGroup",
			Values: map[string]any{},
		}

		result, err := parsedGroup.GetListenAddrSlices("missingFlag")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "flag 'missingFlag' not found in group 'testGroup'")
	})

	t.Run("Flag value is invalid type", func(t *testing.T) {
		t.Parallel()

		parsedGroup := &dynflags.ParsedGroup{
			Name:   "testGroup",
			Values: map[string]any{"flag1": 123},
		}

		result, err := parsedGroup.GetListenAddrSlices("flag1")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "flag 'flag1' is not a []string listen address slice")
	})
}

func TestListenAddrSlicesGetBound(t *testing.T) {
	t.Run("ListenAddrSlicesValue - GetBound", func(t *testing.T) {
		val := []string{":8080", "127.0.0.1:9090"}
		bound := &val

		value := dynflags.ListenAddrSlicesValue{Bound: bound}
		assert.Equal(t, val, value.GetBound())

		value = dynflags.ListenAddrSlicesValue{Bound: nil}
		assert.Nil(t, value.GetBound())
	})
}
