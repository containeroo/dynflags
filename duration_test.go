package dynflags_test

import (
	"testing"
	"time"

	"github.com/containeroo/dynflags"

	"github.com/stretchr/testify/assert"
)

func TestDurationValue_Parse(t *testing.T) {
	t.Parallel()

	t.Run("ValidDuration", func(t *testing.T) {
		t.Parallel()

		d := &dynflags.DurationValue{}
		value, err := d.Parse("2h")
		assert.NoError(t, err)
		assert.Equal(t, 2*time.Hour, value)
	})

	t.Run("InvalidDuration", func(t *testing.T) {
		t.Parallel()

		d := &dynflags.DurationValue{}
		_, err := d.Parse("invalid")
		assert.Error(t, err)
	})
}

func TestDurationValue_Set(t *testing.T) {
	t.Parallel()

	t.Run("SetValidDuration", func(t *testing.T) {
		t.Parallel()

		var bound time.Duration
		d := &dynflags.DurationValue{Bound: &bound}
		err := d.Set(1 * time.Minute)
		assert.NoError(t, err)
		assert.Equal(t, 1*time.Minute, bound)
	})

	t.Run("SetInvalidType", func(t *testing.T) {
		t.Parallel()

		var bound time.Duration
		d := &dynflags.DurationValue{Bound: &bound}
		err := d.Set("not a duration")
		assert.Error(t, err)
		assert.Equal(t, time.Duration(0), bound)
	})
}

func TestGroupConfig_Duration(t *testing.T) {
	t.Parallel()

	t.Run("DurationDefault", func(t *testing.T) {
		t.Parallel()

		group := &dynflags.ConfigGroup{Flags: make(map[string]*dynflags.Flag)}
		defaultValue := 5 * time.Second
		bound := group.Duration("timeout", defaultValue, "Timeout duration")
		assert.Equal(t, defaultValue, bound.Default)
		assert.Contains(t, group.Flags, "timeout")
		assert.Equal(t, defaultValue, group.Flags["timeout"].Default)
		assert.Equal(t, dynflags.FlagTypeDuration, group.Flags["timeout"].Type)
	})
}

func TestParsedGroup_GetDuration(t *testing.T) {
	t.Parallel()

	t.Run("GetValidDuration", func(t *testing.T) {
		t.Parallel()

		parsed := &dynflags.ParsedGroup{
			Name:   "test",
			Values: map[string]any{"timeout": 30 * time.Second},
		}
		dur, err := parsed.GetDuration("timeout")
		assert.NoError(t, err)
		assert.Equal(t, 30*time.Second, dur)
	})

	t.Run("GetDurationNotFound", func(t *testing.T) {
		t.Parallel()

		parsed := &dynflags.ParsedGroup{
			Name:   "test",
			Values: map[string]any{},
		}
		_, err := parsed.GetDuration("missing")
		assert.Error(t, err)
	})

	t.Run("GetDurationWrongType", func(t *testing.T) {
		t.Parallel()

		parsed := &dynflags.ParsedGroup{
			Name:   "test",
			Values: map[string]any{"timeout": "not a duration"},
		}
		_, err := parsed.GetDuration("timeout")
		assert.Error(t, err)
	})
}

func TestDurationGetBound(t *testing.T) {
	t.Run("DurationValue - GetBound", func(t *testing.T) {
		var d *time.Duration
		val := 2 * time.Second
		d = &val

		durationValue := dynflags.DurationValue{Bound: d}
		assert.Equal(t, val, durationValue.GetBound())

		durationValue = dynflags.DurationValue{Bound: nil}
		assert.Nil(t, durationValue.GetBound())
	})
}
