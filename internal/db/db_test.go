package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpen(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		db, err := Open(`wtf`)
		assert.Error(t, err)
		assert.Nil(t, db)
		t.Run("auto migrate/clean", func(t *testing.T) {
			assert.Error(t, AutoMigrate(db))
			assert.Error(t, AutoClean(db))
		})
	})
	t.Run("valid", func(t *testing.T) {
		db, err := Open(`host=/var/run/postgresql dbname=testing sslmode=disable`)
		assert.NoError(t, err)
		defer func() { assert.NoError(t, db.Close()) }()
		t.Run("auto migrate/clean", func(t *testing.T) {
			assert.NoError(t, AutoMigrate(db))
			assert.NoError(t, AutoClean(db))
		})
	})
}
