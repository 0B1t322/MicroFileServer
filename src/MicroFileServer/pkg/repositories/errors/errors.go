package errors

import "errors"

var (
	ErrDocumentExist = errors.New("Document already exists")
	ErrDocumentNotFound = errors.New("Document not found")
	ErrNotValidID		= errors.New("ID is not valid")
)