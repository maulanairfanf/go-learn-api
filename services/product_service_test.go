package services

import (
	"errors"
	"myapi/models"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateCategoryIDs_DuplicateIDs(t *testing.T) {
	err := validateCategoryIDs([]uint{1, 1, 2, 2}, []models.Category{
		{ID: 1},
		{ID: 2},
	})

	var validationErr *CategoryValidationError
	require.ErrorAs(t, err, &validationErr)
	assert.Equal(t, "categories", validationErr.Field)
	assert.Equal(t, "DUPLICATE_CATEGORY_IDS", validationErr.Code)
	assert.Equal(t, []uint{1, 2}, validationErr.DuplicateIDs)
	assert.Empty(t, validationErr.InvalidIDs)
	assert.Equal(t, "Categories must not contain duplicate IDs", validationErr.Error())
}

func TestMapProductConflict(t *testing.T) {
	tests := []struct {
		name       string
		constraint string
		wantField  string
		wantError  string
	}{
		{
			name:       "slug",
			constraint: "uni_products_slug",
			wantField:  "slug",
			wantError:  "Product slug already exists",
		},
		{
			name:       "sku",
			constraint: "uni_products_sku",
			wantField:  "sku",
			wantError:  "Product SKU already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mapProductConflict(&pgconn.PgError{
				Code:           "23505",
				ConstraintName: tt.constraint,
			})

			var conflictErr *ProductConflictError
			require.ErrorAs(t, err, &conflictErr)
			assert.Equal(t, tt.wantField, conflictErr.Field)
			assert.Equal(t, tt.wantError, conflictErr.Error())
			assert.ErrorIs(t, err, ErrConflict)
		})
	}
}

func TestMapProductConflict_UnknownError(t *testing.T) {
	err := errors.New("database unavailable")
	assert.ErrorIs(t, mapProductConflict(err), err)
}

func TestValidateCategoryIDs_AllIDsMissing(t *testing.T) {
	err := validateCategoryIDs([]uint{9, 10}, []models.Category{{ID: 1}})

	var validationErr *CategoryValidationError
	require.ErrorAs(t, err, &validationErr)
	assert.Equal(t, "INVALID_CATEGORY_IDS", validationErr.Code)
	assert.Equal(t, []uint{9, 10}, validationErr.InvalidIDs)
	assert.Equal(t, "One or more category IDs do not exist", validationErr.Error())
}

func TestValidateCategoryIDs_PartiallyMissingIDs(t *testing.T) {
	err := validateCategoryIDs([]uint{1, 9, 2, 10}, []models.Category{
		{ID: 1},
		{ID: 2},
	})

	var validationErr *CategoryValidationError
	require.ErrorAs(t, err, &validationErr)
	assert.Equal(t, []uint{9, 10}, validationErr.InvalidIDs)
}
