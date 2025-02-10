package utils

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

func DirExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist: %w", dir, err)
	}

	return nil
}

func DirIsWritable(dir string) error {
	if err := unix.Access(dir, unix.W_OK); err != nil {
		return fmt.Errorf("%s is not writable: %w", dir, err)
	}

	return nil
}

func CreateFileIfNotExists(file string) error {
	_, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("cannot create or append to file %s: %w", file, err)
	}

	return nil
}
