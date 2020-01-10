package postgres

import (
	"database/sql"
	"fmt"
	"regexp"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

var (
	detailParse = regexp.MustCompile(`(?m)\((\w+(, \w+)*)\)`)

	DuplicateValueErr   = errors.New("Trying to write Duplicate Value to the DB")
	WrongForeignKeyErr  = errors.New("Trying to create record with wrong Foreign Key.")
	TableDoesntExistErr = errors.New("Undefined table. Table doesnt exist in database.")
	NotFoundErr         = errors.New("Not found.")
)

const (
	errCodeDuplicate       pq.ErrorCode = "23505"
	errCodeWrongForeignKey pq.ErrorCode = "23503"
	errCodeUndefinedTable  pq.ErrorCode = "42P01"
)

type DuplicateError struct {
	Field string
}

func (e *DuplicateError) Error() string {
	return fmt.Sprintf("duplicate value on field %q", e.Field)
}

func convertError(err error) error {
	if err == nil {
		return nil
	}

	switch err := errors.Cause(err).(type) {
	case *pq.Error:
		switch err.Code {
		// try to write duplicate value into table
		case errCodeDuplicate:
			return DuplicateValueErr

		// create record with wrong Foreign Key
		case errCodeWrongForeignKey:
			return WrongForeignKeyErr

		case errCodeUndefinedTable:
			return TableDoesntExistErr
		}
	}

	if err == sql.ErrNoRows {
		return NotFoundErr
	}
	return err
}

// TODO: delete this old functions after all will be fine with new one
//func convertError(err error) error {
//	if err == nil {
//		return nil
//	}
//
//	switch err := errors.Cause(err).(type) {
//	case *pq.Error:
//		switch err.Code {
//		// try to write duplicate value into table
//		case errCodeDuplicate:
//			result := detailParse.FindString(err.Detail)
//
//			field := strings.NewReplacer(
//				"(", "",
//				")", "",
//			).Replace(result)
//
//			return &DuplicateError{Field: field}
//
//		// create record with wrong Foreign Key
//		case errCodeWrongForeignKey:
//			result := detailParse.FindString(err.Detail)
//
//			field := strings.NewReplacer(
//				"(", "",
//				")", "",
//			).Replace(result)
//
//			return &DuplicateError{Field: field}
//		}
//	}
//	return err
//}
