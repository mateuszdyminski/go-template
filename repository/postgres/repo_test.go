package postgres

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_OK_ShouldReturnTrue(t *testing.T) {
	// given
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := &pgRepository{db: db}

	// when
	ok, err := repo.OK(context.Background())

	// then
	assert.NoError(t, err)
	assert.Truef(t, ok, "db should be up and running")
}

func Test_OK_ShouldReturnFalseAndError(t *testing.T) {
	// given
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db.Close()

	repo := &pgRepository{db: db}

	// when
	ok, err := repo.OK(context.Background())

	// then
	assert.Error(t, err, "DB should be closed and return error")
	assert.Falsef(t, ok, "status should return false")
}
