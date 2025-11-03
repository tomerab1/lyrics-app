package repositories

import "errors"

var (
	ErrOidConvFailed = errors.New("failed to convert inserted ID to ObjectID")
	ErrInsertFailed  = errors.New("failed to insert")
	ErrUpdateFailed  = errors.New("failed to update")
	ErrDeleteFailed  = errors.New("failed to delete")
	ErrFindOneFailed = errors.New("failed to find")
	ErrFindAllFailed = errors.New("failed to find all")
)
