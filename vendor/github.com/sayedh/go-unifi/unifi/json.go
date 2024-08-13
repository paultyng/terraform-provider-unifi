package unifi

import (
	"errors"
	"strconv"
	"strings"
)

func emptyBoolToTrue(b *bool) bool {
	if b == nil {
		return true
	}
	return *b
}

// numberOrString handles strings that can also accept JSON numbers.
// For example a field may contain a number or the string "auto".
type numberOrString string

func (e *numberOrString) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	s := string(b)
	if s == `""` {
		*e = ""
		return nil
	}
	var err error
	if strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`) {
		s, err = strconv.Unquote(s)
		if err != nil {
			return err
		}
		*e = numberOrString(s)
		return nil
	}
	*e = numberOrString(string(b))
	return nil
}

// emptyStringInt was created due to the behavior change in
// Go 1.14 with json.Number's handling of empty string.
type emptyStringInt int

func (e *emptyStringInt) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	s := string(b)
	if s == `""` {
		*e = 0
		return nil
	}
	var err error
	if strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`) {
		s, err = strconv.Unquote(s)
		if err != nil {
			return err
		}
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*e = emptyStringInt(i)
	return nil
}

func (e *emptyStringInt) MarshalJSON() ([]byte, error) {
	if e == nil || *e == 0 {
		return []byte(`""`), nil
	}

	return []byte(strconv.Itoa(int(*e))), nil
}

type booleanishString bool

func (e *booleanishString) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == `"enabled"` {
		*e = booleanishString(true)
		return nil
	} else if s == `"disabled"` {
		*e = booleanishString(false)
		return nil
	}
	return errors.New("Could not unmarshal JSON value.")
}
