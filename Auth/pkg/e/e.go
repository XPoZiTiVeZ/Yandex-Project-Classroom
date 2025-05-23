package e

import "fmt"

func Wrap(err error, msg string) error {
	return fmt.Errorf("%s: %w", msg, err)
}

func WrapIfErr(err error, msg string) error {
	if err != nil {
		return Wrap(err, msg)
	}
	return nil
}
