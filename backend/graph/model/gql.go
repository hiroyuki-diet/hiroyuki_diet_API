package model

import (
	"database/sql/driver"
	"fmt"
	"io"

	"github.com/google/uuid"
)

type UUID uuid.UUID

// func ParseUUIDFromString(s string) (UUID, error) {
// 	u, err := uuid.Parse(s)
// 	if err != nil {
// 		return UUID{}, err
// 	}
// 	return UUID(u), nil
// }

func (u *UUID) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("UUID must be a string")
	}
	id, err := uuid.Parse(str)
	if err != nil {
		return err
	}
	*u = UUID(id)
	return nil
}

func (u UUID) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, `"`+uuid.UUID(u).String()+`"`)
}

func (u UUID) String() string {
	return uuid.UUID(u).String()
}

func (u *UUID) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("UUID should be a string, got %T", value)
	}
	id, err := uuid.Parse(str)
	if err != nil {
		return err
	}
	*u = UUID(id)
	return nil
}

func (u UUID) Value() (driver.Value, error) {
	return uuid.UUID(u).String(), nil
}
