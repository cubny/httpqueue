package config

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewWithContext(t *testing.T) {
	got, err := New(context.Background())
	require.NoError(t, err)

	assert.Equal(t, got.DB.TimerMaxTTLDays, 180)
}
