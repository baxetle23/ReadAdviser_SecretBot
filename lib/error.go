package lib

import "fmt"

func WrapOnError(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}

func WrapIfError(msg string, err error) error {
	if err == nil {
		return nil
	}
	return WrapOnError(msg, err)
}
