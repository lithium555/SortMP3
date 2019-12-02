package postgres

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

var detailParse = regexp.MustCompile(`(?m)\((\w+(, \w+)*)\)`)

const (
	errCodeDuplicate       pq.ErrorCode = "23505"
	errCodeWrongForeignKey pq.ErrorCode = "23503"
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
			result := detailParse.FindString(err.Detail)

			field := strings.NewReplacer(
				"(", "",
				")", "",
			).Replace(result)

			return &DuplicateError{Field: field}

		// create record with wrong Foreign Key
		case errCodeWrongForeignKey:
			result := detailParse.FindString(err.Detail)

			field := strings.NewReplacer(
				"(", "",
				")", "",
			).Replace(result)

			return &DuplicateError{Field: field}

		}
	}
	return err
}
