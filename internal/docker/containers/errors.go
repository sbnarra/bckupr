package containers

import (
	"fmt"
)

type MissingTemplate struct {
	Message string
}

func (mt *MissingTemplate) Error() string {
	return mt.Message
}

func (mt *MissingTemplate) Is(target error) bool {
	return fmt.Sprintf("%T", target) == fmt.Sprintf("%T", mt)
}
