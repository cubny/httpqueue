package internal

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWithContext(t *testing.T) {
	app, err := Init(context.Background())
	require.NoError(t, err)
	assert.Equal(t, context.Background(), app.ctx)
}

func TestApp_StopOnError(t *testing.T) {
	app := &App{err: fmt.Errorf("bOOm")}
	testFn := func(fnToTest func() *App) func(t *testing.T) {
		return func(t *testing.T) {
			returned := fnToTest()
			assert.Equal(t, app, returned)
		}
	}

	t.Run("initConsumer", testFn(app.initConsumer))
	t.Run("initRelay", testFn(app.initRelay))
	t.Run("initConfig", testFn(app.initConfig))
	t.Run("initAPIServer", testFn(app.initAPIServer))
	t.Run("initPromHandler", testFn(app.initPromHandler))
	t.Run("initService", testFn(app.initService))
	t.Run("initRepo", testFn(app.initRepo))
}
