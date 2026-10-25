package services

import "errors"

var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrInvalidStatus     = errors.New("invalid status")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrConflict          = errors.New("conflict")
)

type CategoryValidationError struct {
	Field        string `json:"field"`
	Code         string `json:"code"`
	InvalidIDs   []uint `json:"invalid_ids,omitempty"`
	DuplicateIDs []uint `json:"duplicate_ids,omitempty"`
}

func (e *CategoryValidationError) Error() string {
	if e.Code == "DUPLICATE_CATEGORY_IDS" {
		return "Categories must not contain duplicate IDs"
	}
	return "One or more category IDs do not exist"
}

func (e *CategoryValidationError) Unwrap() error {
	return ErrInvalidInput
}

var _ error = (*CategoryValidationError)(nil)

type ConflictError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ProductConflictError struct {
	Field   string
	Message string
}

func (e *ProductConflictError) Error() string {
	return e.Message
}

func (e *ProductConflictError) Unwrap() error {
	return ErrConflict
}
