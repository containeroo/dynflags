package dynflags_test

import (
	"testing"

	"github.com/containeroo/dynflags"
	"github.com/stretchr/testify/assert"
)

func TestListenAddrValue(t *testing.T) {
	t.Parallel()

	t.Run("Parse valid listen address", func(t *testing.T) {
		t.Parallel()

		value := dynflags.ListenAddrValue{}
		parsed, err := value.Parse(":8080")
		assert.NoError(t, err)
		assert.NotNil(t, parsed)
		assert.Equal(t, ":8080", *(parsed.(*string)))
	})

	t.Run("Parse invalid listen address", func(t *testing.T) {
		t.Parallel()

		value := dynflags.ListenAddrValue{}
		parsed, err := value.Parse("no-port")
		assert.Error(t, err)
		assert.Nil(t, parsed)
	})

	t.Run("Set valid listen address", func(t *testing.T) {
		t.Parallel()

		bound := ":9090"
		value := dynflags.ListenAddrValue{Bound: &bound}

		newVal := ":8080"
		err := value.Set(&newVal)
		assert.NoError(t, err)
		assert.Equal(t, ":8080", *value.Bound)
	})

	t.Run("Set invalid value type", func(t *testing.T) {
		t.Parallel()

		bound := ":9090"
		value := dynflags.ListenAddrValue{Bound: &bound}

		err := value.Set(1234)
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid value type: expected string pointer for listen address")
	})
}

func TestGroupConfigListenAddr(t *testing.T) {
	t.Parallel()

	t.Run("Define listen address flag with valid default", func(t *testing.T) {
		t.Parallel()

		group := &dynflags.ConfigGroup{Flags: make(map[string]*dynflags.Flag)}
		defaultAddr := ":9090"
		group.ListenAddr("listen", defaultAddr, "Listen address")

		assert.Contains(t, group.Flags, "listen")
		assert.Equal(t, "Listen address", group.Flags["listen"].Usage)
		assert.Equal(t, defaultAddr, group.Flags["listen"].Default)
	})

	t.Run("Define listen address flag with invalid default", func(t *testing.T) {
		t.Parallel()

		group := &dynflags.ConfigGroup{Flags: make(map[string]*dynflags.Flag)}

		assert.PanicsWithValue(t,
			"listen has an invalid default listen address 'bad:address': lookup tcp/address: unknown port",
			func() {
				group.ListenAddr("listen", "bad:address", "Broken default")
			})
	})
}

func TestParsedGroupGetListenAddr(t *testing.T) {
	t.Parallel()

	t.Run("Get existing listen address flag", func(t *testing.T) {
		t.Parallel()

		group := &dynflags.ParsedGroup{
			Values: map[string]any{
				"listen": ":9090",
			},
		}
		addr, err := group.GetListenAddr("listen")
		assert.NoError(t, err)
		assert.Equal(t, ":9090", addr)
	})

	t.Run("Get non-existent listen address flag", func(t *testing.T) {
		t.Parallel()

		group := &dynflags.ParsedGroup{
			Values: map[string]any{},
		}
		addr, err := group.GetListenAddr("listen")
		assert.Error(t, err)
		assert.Equal(t, "", addr)
		assert.EqualError(t, err, "flag 'listen' not found in group ''")
	})

	t.Run("Get listen address flag with invalid type", func(t *testing.T) {
		t.Parallel()

		group := &dynflags.ParsedGroup{
			Values: map[string]any{
				"listen": 9090,
			},
		}
		addr, err := group.GetListenAddr("listen")
		assert.Error(t, err)
		assert.Equal(t, "", addr)
		assert.EqualError(t, err, "flag 'listen' is not a string listen address")
	})
}

func TestListenAddrGetBound(t *testing.T) {
	t.Run("ListenAddrValue - GetBound", func(t *testing.T) {
		bound := ":9090"
		value := dynflags.ListenAddrValue{Bound: &bound}
		assert.Equal(t, ":9090", value.GetBound())

		value = dynflags.ListenAddrValue{Bound: nil}
		assert.Nil(t, value.GetBound())
	})
}
