package e

import "fmt"

func Wrap(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}

func WrapIfErr(op string, err error) error {
	if err == nil {
		return nil
	}

	return Wrap(op, err)
}
