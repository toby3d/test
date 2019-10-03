package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/toby3d/test/internal/handler"
	"gitlab.com/toby3d/test/internal/store"
)

func TestNewServer(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		srv, err := NewServer("wtf", nil)
		assert.Error(t, err)
		t.Run("start/stop", func(t *testing.T) {
			assert.Error(t, srv.Start())
			assert.Error(t, srv.Stop())
		})
	})
	t.Run("valid", func(t *testing.T) {
		srv, err := NewServer(":2368", handler.NewHandler(
			store.NewInMemoryCartStore(), store.NewInMemoryProductStore(),
		))
		assert.NoError(t, err)
		assert.NotNil(t, srv)
		t.Run("start/stop", func(t *testing.T) {
			go func() { assert.NoError(t, srv.Start()) }()
			time.Sleep(100 * time.Millisecond)
			assert.NoError(t, srv.Stop())
		})
	})
}
