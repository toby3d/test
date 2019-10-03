package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	/* TODO(toby3d): Почему-то невалидный адрес не вызывает никаких ошибок (упреждающее прослушивание?)
	t.Run("invalid", func(t *testing.T) {
		c, err := NewClient("wtf")
		assert.Error(t, err)
		t.Run("close", func(t *testing.T) {
			assert.Error(t, c.Close())
		})
	})
	*/
	t.Run("valid", func(t *testing.T) {
		c, err := NewClient(":2368")
		assert.NoError(t, err)
		assert.NotNil(t, c)
		t.Run("close", func(t *testing.T) {
			assert.NoError(t, c.Close())
		})
	})
}
