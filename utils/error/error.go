package error

import (
	"errors"
	"tugas-akhir/constant/messages"
)

var (
	// Password
	ErrFailedHashingPassword = errors.New(messages.FAILED_HASHING_PASSWORD )
	ErrPasswordMismatch      = errors.New(messages.PASSWORD_MISMATCH)

	// // External Service
	// ErrExternalServiceError = errors.New(message.EXTERNAL_SERVICE_ERROR)

	// // Forbidden
	// ErrForbiddenResource = errors.New(message.FORBIDDEN_RESOURCE)

	// // Token
	// ErrFailedGenerateToken = errors.New(message.FAILED_GENERATE_TOKEN)

	// // DuplicateKey
	// ErrDuplicateKey = errors.New(message.DUPLICATE_KEY)

	// pages
	ErrPageNotFound = errors.New(messages.PAGE_NOT_FOUND)

	// ErrNotFound = errors.New(message.NOT_FOUND)
)
