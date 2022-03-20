package model

import "errors"

var (
	// ErrCoinsNotFound error
	ErrCoinsNotFound = errors.New("coins does not exist")
	// ErrInvalidCoinsPair error
	ErrInvalidCoinsPair = errors.New("invalid coins pair")
	// ErrInvalidSize error
	ErrInvalidSize = errors.New("size of data is greater than max size")
)
