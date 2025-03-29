package database

import (
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrNotFound     = errors.New("record not found")
	ErrKeyConflict  = errors.New("key conflict")
	ErrFKeyConflict = errors.New("fkey conflict")
)

func IsRecordNotFoundErr(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrNotFound)
}

func IsKeyConflictErr(err error) bool {
	if errors.Is(err, ErrKeyConflict) {
		return true
	}

	translatedErr := postgres.Dialector{}.Translate(err)
	return errors.Is(translatedErr, gorm.ErrDuplicatedKey)
}

func IsFKeyConflictError(err error) bool {
	if errors.Is(err, ErrFKeyConflict) {
		return true
	}

	translatedErr := postgres.Dialector{}.Translate(err)
	return errors.Is(translatedErr, gorm.ErrForeignKeyViolated)
}
