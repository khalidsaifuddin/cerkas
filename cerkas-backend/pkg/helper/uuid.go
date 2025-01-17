package helper

import (
	"fmt"

	"github.com/gofrs/uuid"
)

func GenerateUUUID() (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", u), nil
}
